package user_server

import (
	"context"
	"fmt"
	dbs2 "github.com/littlebluewhite/Account/app/dbs"
	"github.com/littlebluewhite/Account/dal/model"
	"github.com/littlebluewhite/Account/dal/query"
	"github.com/littlebluewhite/Account/entry/domain"
	"github.com/littlebluewhite/Account/entry/e_user"
	"github.com/littlebluewhite/Account/util/config"
	"github.com/littlebluewhite/Account/util/my_log"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func setUpServer() (us *UserServer) {
	Config := config.NewConfig[config.Config]("../../..", "config", "config", config.Yaml)
	dbTestLog := my_log.NewLog("app/dbs/test")
	DBS := dbs2.NewDbs(dbTestLog, true, Config)
	us = NewUserServer(DBS)
	return
}

// Example test for createToken method
func TestUserServer(t *testing.T) {
	us := setUpServer()
	t.Run("CreateToken", func(t *testing.T) {
		id := 1
		us.createToken(id)

		// Verify that a token was created and stored correctly
		// This might involve checking the mockDbs' methods were called with expected arguments
		// For simplicity, let's assume we directly access the cache (which would actually be mocked)
		tokenMap := us.getID2Token()
		fmt.Println(tokenMap)
		idMap := us.getToken2ID()
		fmt.Println(idMap)
		token, exists := tokenMap[id]

		assert.True(t, exists, "Token should exist for user ID")
		assert.NotEmpty(t, token.value, "Token value should not be empty")
	})
	t.Run("DeleteToken", func(t *testing.T) {
		id := 1
		us.createToken(id)
		err := us.DeleteToken([]int{id})
		assert.NoError(t, err)
		// Verify the token was deleted
		tokenMap := us.getID2Token()
		_, exists := tokenMap[id]
		fmt.Println(tokenMap)
		fmt.Println("---------")
		fmt.Println(exists)
		assert.False(t, exists, "Token should not exist for user ID after deletion")
	})
	t.Run("Login", func(t *testing.T) {

		testUser := model.User{ID: 1, Username: "testuser", Password: "testpassword"}
		userByIDCacheMap := map[int]model.User{1: testUser}
		userByUsernameCacheMap := map[string]model.User{"testuser": testUser}
		us.setUserMaps(userByIDCacheMap, userByUsernameCacheMap)

		// Test correct login
		_, err := us.Login("testuser", "testpassword")
		assert.NoError(t, err)

		// Test wrong password
		_, err = us.Login("testuser", "wrongpassword")
		assert.ErrorIs(t, err, WrongPassword)

		// Test non-existent username
		_, err = us.Login("nonexistent", "password")
		assert.ErrorIs(t, err, NoUsername)
	})
	t.Run("Create", func(t *testing.T) {
		us := setUpServer()

		// Create new user objects
		userCreates := []*e_user.UserCreate{
			{
				Username: "newuser1",
				Password: "password1",
			},
			{
				Username: "newuser2",
				Password: "password2",
			},
		}

		// Perform create
		users, err := us.Create(userCreates)
		assert.NoError(t, err)

		// Check if the users were created
		assert.Len(t, users, len(userCreates))

		// Verify that the users exist in the cache
		userByIDCacheMap, userByUsernameCacheMap := us.getUserMaps()
		for _, createdUser := range users {
			cachedUserByID, existsByID := userByIDCacheMap[int(createdUser.ID)]
			cachedUserByUsername, existsByUsername := userByUsernameCacheMap[createdUser.Username]

			assert.True(t, existsByID, "User should exist in cache by ID")
			assert.True(t, existsByUsername, "User should exist in cache by username")
			assert.Equal(t, createdUser, cachedUserByID, "Cached user by ID should match created user")
			assert.Equal(t, createdUser, cachedUserByUsername, "Cached user by username should match created user")
		}
	})
	t.Run("Update", func(t *testing.T) {

		testName := "testName"
		testUser := model.User{ID: 1, Username: "testuser", Name: &testName}
		userByIDCacheMap := map[int]model.User{1: testUser}
		userByUsernameCacheMap := map[string]model.User{"testuser": testUser}
		us.setUserMaps(userByIDCacheMap, userByUsernameCacheMap)

		// Create update object
		testNewName := "testNewName"
		updateUser := &e_user.UserUpdate{ID: 1, Name: &testNewName}
		updates := []*e_user.UserUpdate{updateUser}

		// Perform update
		_, err := us.Update(updates)
		assert.NoError(t, err)

		// Check if the password was updated
		userByIDCacheMap, _ = us.getUserMaps()
		updatedUser := userByIDCacheMap[1]
		assert.Equal(t, "testNewName", *updatedUser.Name)
		assert.Equal(t, "testuser", updatedUser.Username)
	})
	t.Run("Delete", func(t *testing.T) {
		// Preload users for testing
		result, _ := us.Create([]*e_user.UserCreate{{Username: "user1", Password: "password1"}, {Username: "user2", Password: "password2"}})
		deleteIDs := make([]int32, 0, len(result))
		for _, user := range result {
			deleteIDs = append(deleteIDs, user.ID)
		}
		err := us.Delete(deleteIDs)
		assert.NoError(t, err)

		// Verify the users were deleted from the cache
		userByIDCacheMap, userByUsernameCacheMap := us.getUserMaps()
		for _, id := range deleteIDs {
			_, existsByID := userByIDCacheMap[int(id)]
			assert.False(t, existsByID, "User should not exist in cache by ID after deletion")

			// As usernames are unique, this works correctly
			_, existsByUsername := userByUsernameCacheMap[fmt.Sprintf("user%d", id)]
			assert.False(t, existsByUsername, "User should not exist in cache by username after deletion")
		}

		// Verify the tokens were deleted
		tokenMap := us.getID2Token()
		for _, id := range deleteIDs {
			_, exists := tokenMap[int(id)]
			assert.False(t, exists, "Token should not exist for user ID after deletion")
		}
	})
	t.Run("checkToken", func(t *testing.T) {
		us := setUpServer()
		userID := 1
		us.createToken(userID)
		us.checkToken()
		tokenMap := us.getID2Token()
		token := tokenMap[userID]
		// Check if the timeout decreased
		assert.Equal(t, 599, token.timeout)
	})
	t.Run("setAndGetTokenMap", func(t *testing.T) {
		us := setUpServer()
		// Set token maps
		id2Token := map[int]tokenTime{1: {value: "token1", timeout: 600}}
		token2ID := map[string]int{"token1": 1}
		us.setID2Token(id2Token)
		us.setToken2ID(token2ID)

		// Get token maps and check values
		retrievedID2Token := us.getID2Token()
		retrievedToken2ID := us.getToken2ID()
		assert.Equal(t, id2Token, retrievedID2Token)
		assert.Equal(t, token2ID, retrievedToken2ID)
	})
	t.Run("listDB", func(t *testing.T) {
		us := setUpServer()
		users, err := us.listDB()
		assert.NoError(t, err)
		assert.NotNil(t, users)
	})
	t.Run("findDB", func(t *testing.T) {
		us := setUpServer()
		ctx := context.Background()

		// Preload a user for testing
		testUser := &model.User{ID: 1, Username: "testuser", Password: "testpassword"}
		us.d.GetSql().Create(testUser)

		// Find user by ID
		q := query.Use(us.d.GetSql())
		users, err := us.findDB(ctx, q, []int32{1})
		assert.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, "testuser", users[0].Username)
	})
	t.Run("setAndGetUserMaps", func(t *testing.T) {
		us := setUpServer()

		// Preload a user for testing
		testUser := model.User{ID: 1, Username: "testuser", Password: "testpassword"}
		userByIDCacheMap := map[int]model.User{1: testUser}
		userByUsernameCacheMap := map[string]model.User{"testuser": testUser}

		us.setUserMaps(userByIDCacheMap, userByUsernameCacheMap)
		retrievedUserByIDCacheMap, retrievedUserByUsernameCacheMap := us.getUserMaps()

		assert.Equal(t, userByIDCacheMap, retrievedUserByIDCacheMap)
		assert.Equal(t, userByUsernameCacheMap, retrievedUserByUsernameCacheMap)
	})
}

func TestRegister(t *testing.T) {
	// 假設 user_server 有個結構體並實作了 Register 方法
	// 你需要依據實際程式碼修改此例
	us := setUpServer()

	testInput := domain.Register{
		Username: "testuser",
		Name:     "Test User",
		Password: "Passw0rd",
		Birthday: "2000-01-01",
		Email:    "test@example.com",
		Phone:    "123456789",
		Country:  "TW",
	}

	// 呼叫要測試的 Register 函式
	err := us.Register(testInput)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// 根據你的實際實作邏輯進行驗證
	// 例如，如果 Register 會將使用者資訊存進資料庫，你可以在此查詢資料庫確認資料正確性
	// 以下只是示意
	user, err := us.GetUserByUsername("testuser")
	if err != nil {
		t.Fatalf("Expected to find user 'testuser', got error %v", err)
	}

	if user == nil {
		t.Fatalf("Expected to find a user, got nil")
	}

	// 假設 user 物件內的生日是 time.Time，且已經轉為 UTC
	expectedBirthday, err := time.Parse("2006-01-02", testInput.Birthday)
	if err != nil {
		t.Fatalf("Error parsing birthday: %v", err)
	}

	if user.Birthday == nil {
		t.Errorf("Expected Birthday is not nil")
	} else {
		y, m, d := user.Birthday.Date()
		ey, em, ed := expectedBirthday.Date()
		if y != ey || m != em || d != ed {
			t.Errorf("Expected Birthday %v, got %v", expectedBirthday, user.Birthday)
		}
	}
	// 驗證其他欄位
	if user.Email == nil || *user.Email != testInput.Email {
		t.Errorf("Expected Email %s, got %v", testInput.Email, user.Email)
	}
	if user.Phone == nil || *user.Phone != testInput.Phone {
		t.Errorf("Expected Phone %s, got %v", testInput.Phone, user.Phone)
	}
	if user.Country == nil || *user.Country != testInput.Country {
		t.Errorf("Expected Country %s, got %v", testInput.Country, user.Country)
	}

}

func TestLoginWithToken(t *testing.T) {
	us := setUpServer()

	// Mock a user and token
	testUser := model.User{
		ID:       1,
		Username: "testuser",
		Password: "testpassword",
	}

	// Preload user maps
	userByIDCacheMap := map[int]model.User{1: testUser}
	userByUsernameCacheMap := map[string]model.User{"testuser": testUser}
	us.setUserMaps(userByIDCacheMap, userByUsernameCacheMap)

	// Create a token for the user
	us.createToken(int(testUser.ID))
	tokenMap := us.getID2Token()
	token := tokenMap[int(testUser.ID)].value

	t.Run("Valid Token", func(t *testing.T) {
		// Call LoginWithToken with a valid token
		user, err := us.LoginWithToken(token)
		assert.NoError(t, err, "Expected no error for valid token")
		assert.Equal(t, testUser.ID, user.ID, "User ID should match")
		assert.Equal(t, testUser.Username, user.Username, "Username should match")
	})

	t.Run("Invalid Token", func(t *testing.T) {
		// Call LoginWithToken with an invalid token
		_, err := us.LoginWithToken("invalidToken")
		assert.ErrorIs(t, err, NoToken, "Expected NoToken error for invalid token")
	})

	t.Run("User Not Found", func(t *testing.T) {
		// Remove the user from the cache but keep the token
		us.setUserMaps(map[int]model.User{}, userByUsernameCacheMap)

		_, err := us.LoginWithToken(token)
		assert.ErrorIs(t, err, NoUser, "Expected NoUser error for missing user")
	})
}
