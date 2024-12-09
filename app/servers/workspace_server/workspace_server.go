package workspace_server

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/littlebluewhite/Account/api"
	"github.com/littlebluewhite/Account/dal/model"
	"github.com/littlebluewhite/Account/dal/query"
	"github.com/littlebluewhite/Account/entry/e_w_group"
	"github.com/littlebluewhite/Account/entry/e_w_user"
	"github.com/littlebluewhite/Account/entry/e_workspace"
	"github.com/littlebluewhite/Account/util"
	"github.com/littlebluewhite/Account/util/convert"
	"github.com/littlebluewhite/Account/util/my_log"
	"github.com/patrickmn/go-cache"
	"gorm.io/gen/field"
	"sync"
	"time"
)

type userServer interface {
	ReloadCacheByIDs(ids []int32) (e error)
	Close()
}

type WorkspaceServer struct {
	d  api.Dbs
	l  api.Logger
	wg *sync.WaitGroup
	us userServer
}

func NewWorkspaceServer(dbs api.Dbs, userServer userServer) *WorkspaceServer {
	l := my_log.NewLog("app/workspace_server")
	return &WorkspaceServer{d: dbs, l: l, us: userServer, wg: new(sync.WaitGroup)}
}

func (w *WorkspaceServer) Start(ctx context.Context) {
	w.l.Infoln("Workspace server start")
	e := w.reloadCache()
	if e != nil {
		panic(e)
	}

	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		w.listen(ctx)
	}()
}

func (w *WorkspaceServer) listen(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)
	for {
		select {
		case <-ctx.Done():
			w.l.Infoln("Workspace server stop")
			return
		case <-ticker.C:
		}
	}
}

func (w *WorkspaceServer) Close() {
	w.wg.Wait()
	w.us.Close()
	w.l.Infoln("Workspace server stop")
}

func (w *WorkspaceServer) setCacheMap(cacheMap map[int]model.Workspace) {
	w.d.GetCache().Set("workspace", cacheMap, cache.NoExpiration)
}

func (w *WorkspaceServer) getCacheMap() map[int]model.Workspace {
	var cacheMap map[int]model.Workspace
	if x, found := w.d.GetCache().Get("workspace"); found {
		cacheMap = x.(map[int]model.Workspace)
	} else {
		return make(map[int]model.Workspace)
	}
	return cacheMap
}

func (w *WorkspaceServer) getCacheList() []model.Workspace {
	m := w.getCacheMap()
	result := make([]model.Workspace, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result
}

func (w *WorkspaceServer) listDB() ([]*model.Workspace, error) {
	ws := query.Use(w.d.GetSql()).Workspace
	ctx := context.Background()
	tt, err := ws.WithContext(ctx).Preload(field.Associations).Preload(
		ws.WGroups.WUserGroups).Preload(ws.WUsers.WUserGroups).Find()
	if err != nil {
		return nil, err
	}
	return tt, nil
}

func (w *WorkspaceServer) findDB(ctx context.Context, q *query.Query, ids []int32) ([]*model.Workspace, error) {
	ws := q.Workspace
	workspace, err := ws.WithContext(ctx).Preload(field.Associations).Where(ws.ID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}
	return workspace, nil
}

func (w *WorkspaceServer) reloadCache() (e error) {
	ws, err := w.listDB()
	if err != nil {
		e = err
		return
	}
	cacheMap := make(map[int]model.Workspace)
	for i := 0; i < len(ws); i++ {
		entry := ws[i]
		cacheMap[int(entry.ID)] = *entry
	}
	w.setCacheMap(cacheMap)
	return
}

func (w *WorkspaceServer) Create(ec []*e_workspace.WorkspaceCreate) ([]model.Workspace, error) {
	q := query.Use(w.d.GetSql())
	ctx := context.Background()
	cacheMap := w.getCacheMap()
	ids := make([]int32, 0, len(ec))
	workspaces := convert.CreateConvert[model.Workspace, e_workspace.WorkspaceCreate](ec)
	result := make([]model.Workspace, 0, len(workspaces))
	err := q.Transaction(func(tx *query.Query) error {
		if err := tx.Workspace.WithContext(ctx).CreateInBatches(workspaces, 100); err != nil {
			return err
		}
		for _, item := range workspaces {
			ids = append(ids, item.ID)
		}
		newWorkspaces, err := w.findDB(ctx, tx, ids)
		if err != nil {
			return err
		}
		for _, workspace := range newWorkspaces {
			cacheMap[int(workspace.ID)] = *workspace
			result = append(result, *workspace)
		}
		w.setCacheMap(cacheMap)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return result, nil
}

func (w *WorkspaceServer) Update(eu []*e_workspace.WorkspaceUpdate) error {
	cacheMap := w.getCacheMap()
	cacheList := w.getCacheList()
	result, err := convert.UpdateConvert[model.Workspace, *e_workspace.WorkspaceUpdate](cacheList, eu, "ID")
	if err != nil {
		return err
	}
	ids := make([]int32, 0, len(result))
	q := query.Use(w.d.GetSql())
	ctx := context.Background()
	err = q.Transaction(func(tx *query.Query) error {
		wUpdate := make([]map[string]interface{}, 0, len(result))

		wuUpdate := make([]map[string]interface{}, 0, len(result))
		wuCreate := make([]*model.WUser, 0, len(result))
		wuDelete := make([]int32, 0, len(result))

		wgUpdate := make([]map[string]interface{}, 0, len(result))
		wgCreate := make([]*model.WGroup, 0, len(result))
		wgDelete := make([]int32, 0, len(result))

		wugCreate := make([]*model.WUserGroup, 0, len(result))
		wugDelete := make([]int32, 0, len(result))
		for _, workspace := range result {
			ids = append(ids, workspace.ID)
			for _, wUser := range workspace.WUsers {
				wu := wUser
				switch {
				case wu.ID < 0:
					wuDelete = append(wuDelete, -wu.ID)
				case wu.ID == 0:
					wu.WorkspaceID = workspace.ID
					wuCreate = append(wuCreate, &wu)
				case wu.ID > 0:
					wuUpdate = append(wuUpdate, util.StructToMap(wu))
					for _, wug := range wu.WUserGroups {
						wgID := wug.WGroupID
						switch {
						case wgID > 0:
							wug.WUserID = wu.ID
							wugCreate = append(wugCreate, &wug)
						case wgID < 0:
							wugDelete = append(wugDelete, wu.ID)
						}
					}
				}
			}
			for _, wGroup := range workspace.WGroups {
				wg := wGroup
				switch {
				case wg.ID < 0:
					wgDelete = append(wgDelete, -wg.ID)
				case wg.ID == 0:
					wg.WorkspaceID = workspace.ID
					wgCreate = append(wgCreate, &wg)
				case wg.ID > 0:
					wgUpdate = append(wgUpdate, util.StructToMap(wg))
					for _, wug := range wg.WUserGroups {
						wuID := wug.WUserID
						switch {
						case wuID > 0:
							wug.WGroupID = wg.ID
							wugCreate = append(wugCreate, &wug)
						case wuID < 0:
							wugDelete = append(wugDelete, -wuID)
						}
					}
				}
			}
			t := util.StructToMap(workspace)
			delete(t, "updated_at")
			delete(t, "created_at")
			delete(t, "w_users")
			delete(t, "w_groups")
			delete(t, "next_workspaces")
			wUpdate = append(wUpdate, t)
		}
		// sql operations
		for _, workspace := range wUpdate {
			if _, err = tx.Workspace.WithContext(ctx).Where(tx.Workspace.ID.Eq(workspace["id"].(int32))).Updates(workspace); err != nil {
				return err
			}
		}
		for _, wuU := range wuUpdate {
			delete(wuU, "updated_at")
			delete(wuU, "created_at")
			delete(wuU, "w_user_groups")
			if _, err = tx.WUser.WithContext(ctx).Where(tx.WUser.ID.Eq(wuU["id"].(int32))).Updates(wuU); err != nil {
				return err
			}
		}
		for _, wgU := range wgUpdate {
			delete(wgU, "updated_at")
			delete(wgU, "created_at")
			if _, err = tx.WGroup.WithContext(ctx).Where(tx.WGroup.ID.Eq(wgU["id"].(int32))).Updates(wgU); err != nil {
				return err
			}
		}
		if err = tx.WUser.WithContext(ctx).CreateInBatches(wuCreate, 100); err != nil {
			return err
		}
		if err = tx.WGroup.WithContext(ctx).CreateInBatches(wgCreate, 100); err != nil {
			return err
		}
		if _, err = tx.WUser.WithContext(ctx).Where(tx.WUser.ID.In(wuDelete...)).Delete(); err != nil {
			return err
		}
		if _, err = tx.WGroup.WithContext(ctx).Where(tx.WGroup.ID.In(wgDelete...)).Delete(); err != nil {
			return err
		}
		if _, err = tx.WUserGroup.WithContext(ctx).Where(tx.WUserGroup.WUserID.In(wugDelete...)).Delete(); err != nil {
		}
		// update cache
		newWorkspace, err := w.findDB(ctx, tx, ids)
		if err != nil {
			return err
		}
		for _, workspace := range newWorkspace {
			cacheMap[int(workspace.ID)] = *workspace
		}
		w.setCacheMap(cacheMap)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (w *WorkspaceServer) getWorkspaceById(workspaceId int32) model.Workspace {
	cacheMap := w.getCacheMap()
	workspace, ok := cacheMap[int(workspaceId)]
	if !ok {
		w.l.Infoln("Workspace ID %d not found in cacheMap", workspaceId)
		return model.Workspace{}
	}
	return workspace
}

func (w *WorkspaceServer) getDefaultUserAuth(workspace model.Workspace) json.RawMessage {
	var constMap map[string]interface{}
	var passDownMap map[string]interface{}
	var customMap map[string]interface{}
	if err := json.Unmarshal(workspace.UserAuthConst, &constMap); err != nil {
		w.l.Infoln("Failed to unmarshal UserAuthConst for Workspace ID %d: %v", workspace.ID, err)
	}
	if err := json.Unmarshal(workspace.UserAuthPassDown, &passDownMap); err != nil {
		w.l.Infoln("Failed to unmarshal UserAuthPassDown for Workspace ID %d: %v", workspace.ID, err)
	}
	if err := json.Unmarshal(workspace.UserAuthCustom, &customMap); err != nil {
		w.l.Infoln("Failed to unmarshal UserAuthCustom for Workspace ID %d: %v", workspace.ID, err)
	}
	result := make(map[string]interface{})
	for k, v := range customMap {
		result[k] = v
	}
	for k, v := range passDownMap {
		result[k] = v
	}
	for k, v := range constMap {
		result[k] = v
	}

	if workspace.ID == 0 {
		w.l.Infoln("no workspace")
		return nil
	}

	defaultAuth, err := json.Marshal(result)
	if err != nil {
		w.l.Infoln("Failed to marshal defaultAuthMap for Workspace ID %d: %v", workspace.ID, err)
		return nil
	}
	return defaultAuth
}

func (w *WorkspaceServer) syncUserGroupAuth(workSpaceId int32) {
	workspace := w.getWorkspaceById(workSpaceId)
	defaultUserAuth := w.getDefaultUserAuth(workspace)
	wUserUpdate := make([]e_w_user.WUserUpdate, 0, len(workspace.WUsers))
	wGroupUpdate := make([]e_w_group.WGroupUpdate, 0, len(workspace.WGroups))
	workspaceUpdate := e_workspace.WorkspaceUpdate{ID: workSpaceId}
	for _, wUser := range workspace.WUsers {
		newAuthByte, err := convert.MergeJSON(wUser.Auth, defaultUserAuth)
		var newAuth json.RawMessage = newAuthByte
		if err != nil {
			w.l.Infoln("Failed to merge defaultUserAuth for WUser ID %d: %v", wUser.ID, err)
			continue
		}
		wUserUpdate = append(wUserUpdate, e_w_user.WUserUpdate{
			ID:   wUser.ID,
			Auth: &newAuth,
		})
	}
	for _, wGroup := range workspace.WGroups {
		newAuthByte, err := convert.MergeJSON(wGroup.DefaultAuth, defaultUserAuth)
		var newAuth json.RawMessage = newAuthByte
		if err != nil {
			w.l.Infoln("Failed to merge defaultUserAuth for WGroup ID %d: %v", wGroup.ID, err)
		}
		wGroupUpdate = append(wGroupUpdate, e_w_group.WGroupUpdate{
			ID:          wGroup.ID,
			DefaultAuth: &newAuth,
		})
	}
	workspaceUpdate.WUsers = wUserUpdate
	err := w.Update([]*e_workspace.WorkspaceUpdate{&workspaceUpdate})
	if err != nil {
		w.l.Errorln("Failed to update WUser for Workspace ID %d: %v", workSpaceId, err)
	}
}
