package user_server

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/littlebluewhite/Account/dal/model"
	"github.com/littlebluewhite/Account/dal/query"
	"github.com/littlebluewhite/Account/entry/domain"
	"github.com/littlebluewhite/Account/entry/e_user"
	"github.com/littlebluewhite/Account/util"
	"github.com/littlebluewhite/Account/util/convert"
	"github.com/littlebluewhite/Account/util/my_cache"
	"github.com/littlebluewhite/Account/util/my_log"
	"github.com/patrickmn/go-cache"
	"gorm.io/gen/field"
	"sync"
	"time"
)

type UserServer struct {
	d    domain.Dbs
	l    domain.Logger
	salt string
	wg   *sync.WaitGroup
}

type tokenTime struct {
	value   string
	timeout int
}

func NewUserServer(d domain.Dbs) *UserServer {
	l := my_log.NewLog("app/user_server")
	return &UserServer{d: d, l: l, salt: "Wilson", wg: new(sync.WaitGroup)}
}

func (u *UserServer) Start(ctx context.Context) {
	u.l.Infoln("User server start")
	e := u.reloadCache()
	if e != nil {
		panic(e)
	}

	u.wg.Add(1)
	go func() {
		defer u.wg.Done()
		u.listen(ctx)
	}()
}

func (u *UserServer) Close() {
	u.wg.Wait()
	u.l.Infoln("User server stop")
}

func (u *UserServer) listen(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			u.l.Infoln("User server stop")
			return
		case <-ticker.C:
			go u.checkToken()
		}
	}
}

func (u *UserServer) setUserMaps(userByIDCacheMap map[int]model.User, userByUsernameCacheMap map[string]model.User) {
	u.d.GetCache().Set("user_by_id", userByIDCacheMap, cache.NoExpiration)
	u.d.GetCache().Set("user_by_username", userByUsernameCacheMap, cache.NoExpiration)
}

func (u *UserServer) getUserMaps() (map[int]model.User, map[string]model.User) {
	userByIdCacheMap := my_cache.GetCacheMap[model.User, int](u.d.GetCache(), "user_by_id")
	userByUsernameCacheMap := my_cache.GetCacheMap[model.User, string](u.d.GetCache(), "user_by_username")
	return userByIdCacheMap, userByUsernameCacheMap
}

func (u *UserServer) getUserList() []model.User {
	userByIDCacheMap, _ := u.getUserMaps()
	userByIDCacheList := make([]model.User, 0, len(userByIDCacheMap))
	for _, v := range userByIDCacheMap {
		userByIDCacheList = append(userByIDCacheList, v)
	}
	return userByIDCacheList
}

func (u *UserServer) listDB() ([]*model.User, error) {
	wu := query.Use(u.d.GetSql()).User
	ctx := context.Background()
	tt, err := wu.WithContext(ctx).Preload(field.Associations).Preload(
		wu.WUsers.WUserGroups).Find()
	if err != nil {
		return nil, err
	}
	return tt, nil
}

func (u *UserServer) listDBByIDs(ids []int32) ([]*model.User, error) {
	wu := query.Use(u.d.GetSql()).User
	ctx := context.Background()
	tt, err := wu.WithContext(ctx).Preload(field.Associations).Preload(
		wu.WUsers.WUserGroups).Where(wu.ID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}
	return tt, nil
}

func (u *UserServer) findDB(ctx context.Context, q *query.Query, ids []int32) ([]*model.User, error) {
	wu := q.User
	users, err := wu.WithContext(ctx).Preload(field.Associations).Preload(
		wu.WUsers.WUserGroups).Where(wu.ID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserServer) reloadCache() (e error) {
	wu, err := u.listDB()
	if err != nil {
		e = err
		return
	}
	userByIDCacheMap := make(map[int]model.User)
	userByUsernameCacheMap := make(map[string]model.User)
	for i := 0; i < len(wu); i++ {
		entry := wu[i]
		userByIDCacheMap[int(entry.ID)] = *entry
		userByUsernameCacheMap[entry.Username] = *entry
	}
	u.setUserMaps(userByIDCacheMap, userByUsernameCacheMap)
	return
}

func (u *UserServer) ReloadCacheByIDs(ids []int32) (e error) {
	wu, err := u.listDBByIDs(ids)
	if err != nil {
		e = err
		return
	}
	userByIDCacheMap, userByUsernameCacheMap := u.getUserMaps()
	for i := 0; i < len(wu); i++ {
		entry := wu[i]
		userByIDCacheMap[int(entry.ID)] = *entry
		userByUsernameCacheMap[entry.Username] = *entry
	}
	u.setUserMaps(userByIDCacheMap, userByUsernameCacheMap)
	return
}

func (u *UserServer) getID2Token() map[int]tokenTime {
	var cacheMap map[int]tokenTime
	if x, found := u.d.GetCache().Get("id2token"); found {
		cacheMap = x.(map[int]tokenTime)
	} else {
		return make(map[int]tokenTime)
	}
	return cacheMap
}

func (u *UserServer) setID2Token(tokenMap map[int]tokenTime) {
	u.d.GetCache().Set("id2token", tokenMap, cache.NoExpiration)
}

func (u *UserServer) getToken2ID() map[string]int {
	var cacheMap map[string]int
	if x, found := u.d.GetCache().Get("token2id"); found {
		cacheMap = x.(map[string]int)
	} else {
		return make(map[string]int)
	}
	return cacheMap
}

func (u *UserServer) setToken2ID(tokenMap map[string]int) {
	u.d.GetCache().Set("token2id", tokenMap, cache.NoExpiration)
}

func (u *UserServer) getTokenCache() (token2ID map[string]int, id2Token map[int]tokenTime) {
	return u.getToken2ID(), u.getID2Token()
}

func (u *UserServer) setTokenCache(token2ID map[string]int, id2Token map[int]tokenTime) {
	u.d.GetCache().Set("token2id", token2ID, cache.NoExpiration)
	u.d.GetCache().Set("id2token", id2Token, cache.NoExpiration)
}

func (u *UserServer) checkToken() {
	tokenId, idToken := u.getTokenCache()
	for id, t := range idToken {
		if t.timeout <= 0 {
			delete(idToken, id)
			delete(tokenId, t.value)
		} else {
			t.timeout--
			idToken[id] = t
		}
	}
	u.setTokenCache(tokenId, idToken)
}

func (u *UserServer) createToken(id int) {
	timestamp := time.Now().Unix()
	data := fmt.Sprintf("%d:%d", id, timestamp)
	hasher := sha256.New()
	hash := []byte(data)
	for i := 0; i < 3; i++ {
		hasher.Write(append(hash, []byte(u.salt)...))
		hash = hasher.Sum(nil)
	}
	t := hex.EncodeToString(hash)
	u.SetToken(id, t)
}

func (u *UserServer) DeleteToken(ids []int) (err error) {
	tokenId, idToken := u.getTokenCache()
	for _, id := range ids {
		t, ok := idToken[id]
		if !ok {
			continue
		}
		delete(idToken, id)
		delete(tokenId, t.value)
	}
	u.setTokenCache(tokenId, idToken)
	return
}

func (u *UserServer) SetToken(id int, token string) {
	tokenId, idToken := u.getTokenCache()
	t, ok := idToken[id]
	// check this user has token
	if ok {
		delete(tokenId, t.value)
	}
	idToken[id] = tokenTime{value: token, timeout: 600}
	tokenId[token] = id
	u.setTokenCache(tokenId, idToken)
}

func (u *UserServer) Create(ec []*e_user.UserCreate) ([]model.User, error) {
	q := query.Use(u.d.GetSql())
	ctx := context.Background()
	userByIDCacheMap, userByUsernameCacheMap := u.getUserMaps()
	ids := make([]int32, 0, len(ec))
	users := convert.CreateConvert[model.User, e_user.UserCreate](ec)
	result := make([]model.User, 0, len(users))
	err := q.Transaction(func(tx *query.Query) error {
		if err := tx.User.WithContext(ctx).CreateInBatches(users, 100); err != nil {
			return err
		}
		for _, item := range users {
			ids = append(ids, item.ID)
		}
		newUsers, err := u.findDB(ctx, tx, ids)
		if err != nil {
			return err
		}
		for _, user := range newUsers {
			userByIDCacheMap[int(user.ID)] = *user
			userByUsernameCacheMap[user.Username] = *user
			result = append(result, *user)
		}
		u.setUserMaps(userByIDCacheMap, userByUsernameCacheMap)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *UserServer) Update(eu []*e_user.UserUpdate) ([]model.User, error) {
	userByIDCacheMap, userByUsernameCacheMap := u.getUserMaps()
	userByIDCacheList := u.getUserList()
	result, err := convert.UpdateConvert[model.User, *e_user.UserUpdate](userByIDCacheList, eu, "ID")
	if err != nil {
		return nil, err
	}
	ids := make([]int32, 0, len(result))
	q := query.Use(u.d.GetSql())
	ctx := context.Background()
	err = q.Transaction(func(tx *query.Query) error {
		for _, item := range result {
			ids = append(ids, item.ID)
			t := util.StructToMap(item)
			delete(t, "updated_at")
			delete(t, "created_at")
			delete(t, "w_users")
			_, err = tx.User.WithContext(ctx).Where(tx.User.ID.Eq(item.ID)).Updates(t)
			if err != nil {
				return err
			}
		}
		newUsers, e := u.findDB(ctx, tx, ids)
		if e != nil {
			return e
		}
		for _, user := range newUsers {
			userByIDCacheMap[int(user.ID)] = *user
			userByUsernameCacheMap[user.Username] = *user
		}
		u.setUserMaps(userByIDCacheMap, userByUsernameCacheMap)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *UserServer) Delete(ids []int32) error {
	userByIDCacheMap, userByUsernameCacheMap := u.getUserMaps()
	q := query.Use(u.d.GetSql())
	ctx := context.Background()
	err := q.Transaction(func(tx *query.Query) error {
		if _, err := tx.User.WithContext(ctx).Where(tx.User.ID.In(ids...)).Delete(); err != nil {
			return err
		}
		for _, id := range ids {
			deleteIDs := make([]int, 0, len(ids))
			user, ok := userByIDCacheMap[int(id)]
			if ok {
				delete(userByIDCacheMap, int(id))
				delete(userByUsernameCacheMap, user.Username)
				deleteIDs = append(deleteIDs, int(id))
			}
			_ = u.DeleteToken(deleteIDs)
		}
		return nil
	})
	if err != nil {
		return err
	}
	u.setUserMaps(userByIDCacheMap, userByUsernameCacheMap)
	return nil
}

func (u *UserServer) Login(username string, password string) (user model.User, err error) {
	_, userByUsernameCacheMap := u.getUserMaps()
	user, ok := userByUsernameCacheMap[username]
	if !ok {
		err = NoUsername
		return
	}
	if user.Password != password {
		err = WrongPassword
		return
	}
	u.createToken(int(user.ID))
	return
}

func (u *UserServer) Register(register domain.Register) error {
	_, userByUsernameCacheMap := u.getUserMaps()
	_, ok := userByUsernameCacheMap[register.Username]
	if ok {
		return UsernameExist
	}
	CreateUsers := convert.CreateConvert[e_user.UserCreate, domain.Register]([]*domain.Register{&register})
	_, err := u.Create(CreateUsers)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserServer) GetUserByUsername(username string) (*model.User, error) {
	_, userByUsernameCacheMap := u.getUserMaps()
	user, ok := userByUsernameCacheMap[username]
	if !ok {
		return nil, NoUsername
	}
	return &user, nil
}

func (u *UserServer) LoginWithToken(token string) (user model.User, err error) {
	userByIDCacheMap, _ := u.getUserMaps()
	id, ok := u.getToken2ID()[token]
	if !ok {
		err = NoToken
		return
	}
	user, ok = userByIDCacheMap[id]
	if !ok {
		err = NoUser
		return
	}
	u.SetToken(id, token)
	return
}
