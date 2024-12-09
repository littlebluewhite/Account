package user_server

import (
	"context"
	"fmt"
	dbs2 "github.com/littlebluewhite/Account/app/dbs"
	"github.com/littlebluewhite/Account/dal/model"
	"github.com/littlebluewhite/Account/dal/query"
	"github.com/littlebluewhite/Account/entry/e_user"
	"github.com/littlebluewhite/Account/util/config"
	"github.com/littlebluewhite/Account/util/my_log"
	"github.com/stretchr/testify/assert"
	"testing"
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
		err := us.Login("testuser", "testpassword")
		assert.NoError(t, err)

		// Test wrong password
		err = us.Login("testuser", "wrongpassword")
		assert.ErrorIs(t, err, WrongPassword)

		// Test non-existent username
		err = us.Login("nonexistent", "password")
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
