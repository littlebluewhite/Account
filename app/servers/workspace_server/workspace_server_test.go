// workspace_server_test.go
package workspace_server

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/littlebluewhite/Account/api"
	dbs2 "github.com/littlebluewhite/Account/app/dbs"
	"github.com/littlebluewhite/Account/app/servers/user_server"
	"github.com/littlebluewhite/Account/dal/model"
	"github.com/littlebluewhite/Account/dal/query" // Import the query package
	"github.com/littlebluewhite/Account/entry/e_w_user"
	"github.com/littlebluewhite/Account/entry/e_workspace"
	"github.com/littlebluewhite/Account/util/config"
	"github.com/littlebluewhite/Account/util/convert"
	"github.com/littlebluewhite/Account/util/my_log"
	"gorm.io/gorm"
	"reflect"
	"testing"
	"time"
)

// setUpServer initializes the WorkspaceServer with a test database and user server.
func setUpServer(ctx context.Context) *WorkspaceServer {
	Config := config.NewConfig[config.Config]("../../..", "config", "config", config.Yaml)

	// Initialize logger
	dbTestLog := my_log.NewLog("app/dbs/test")

	// Initialize database connections (ensure this points to a test DB)
	DBS := dbs2.NewDbs(dbTestLog, true, Config)

	// Initialize UserServer
	us := user_server.NewUserServer(DBS)

	// Initialize WorkspaceServer
	ws := NewWorkspaceServer(DBS, us)

	// Start the WorkspaceServer
	go func() {
		ws.Start(ctx)
	}()

	return ws
}

// cleanUpDatabase cleans up the workspace table after tests.
func cleanUpDatabase(t *testing.T, dbs api.Dbs) {
	q := query.Use(dbs.GetSql()) // Initialize the query builder
	ctx := context.Background()

	// Use the underlying GORM DB to delete all records
	db := q.Workspace.WithContext(ctx).UnderlyingDB()
	err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Workspace{}).Error
	if err != nil {
		t.Fatalf("Failed to clean up database: %v", err)
	}
}

func TestWorkspaceServer(t *testing.T) {
	// Create a cancellable context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Set up the server
	ws := setUpServer(ctx)
	defer ws.Close()

	// Ensure the database is clean before starting tests
	cleanUpDatabase(t, ws.d)

	t.Run("Create", func(t *testing.T) {
		// Define sample workspaces to create
		workspaceCreates := []*e_workspace.WorkspaceCreate{
			{
				Name:             "Workspace A",
				PreWorkspaceID:   nil,
				Rank:             1,
				Ancient:          "Ancient A",
				Enable:           true,
				OwnerID:          1,
				ExpiryDate:       time.Now().AddDate(1, 0, 0),
				Auth:             json.RawMessage(`{"role": "admin"}`),
				UserAuthConst:    json.RawMessage(`{"access": "full"}`),
				UserAuthPassDown: json.RawMessage(`{"pass": "down"}`),
				UserAuthCustom:   json.RawMessage(`{"custom": "rules"}`),
			},
			{
				Name:             "Workspace B",
				PreWorkspaceID:   nil,
				Rank:             2,
				Ancient:          "Ancient B",
				Enable:           false,
				OwnerID:          1,
				ExpiryDate:       time.Now().AddDate(2, 0, 0),
				Auth:             json.RawMessage(`{"role": "user"}`),
				UserAuthConst:    json.RawMessage(`{"access": "limited"}`),
				UserAuthPassDown: json.RawMessage(`{"pass": "up"}`),
				UserAuthCustom:   json.RawMessage(`{"custom": "permissions"}`),
			},
		}

		// Call the Create method
		createdWorkspaces, err := ws.Create(workspaceCreates)
		if err != nil {
			t.Fatalf("Create method failed: %v", err)
		}

		// Verify the number of created workspaces
		if len(createdWorkspaces) != len(workspaceCreates) {
			t.Fatalf("Expected %d workspaces, got %d", len(workspaceCreates), len(createdWorkspaces))
		}

		// Verify each created workspace
		for i, cw := range createdWorkspaces {
			ec := workspaceCreates[i]
			if cw.Name != ec.Name {
				t.Errorf("Expected Name %s, got %s", ec.Name, cw.Name)
			}
			if cw.Rank != ec.Rank {
				t.Errorf("Expected Rank %d, got %d", ec.Rank, cw.Rank)
			}
			if cw.Ancient != ec.Ancient {
				t.Errorf("Expected Ancient %s, got %s", ec.Ancient, cw.Ancient)
			}
			if cw.Enable != ec.Enable {
				t.Errorf("Expected Enable %v, got %v", ec.Enable, cw.Enable)
			}
			if cw.OwnerID != ec.OwnerID {
				t.Errorf("Expected OwnerID %d, got %d", ec.OwnerID, cw.OwnerID)
			}
			// Additional field checks can be added here
		}

		// Verify that the cache is updated
		cacheMap := ws.getCacheMap()
		for _, cw := range createdWorkspaces {
			if cached, exists := cacheMap[int(cw.ID)]; !exists {
				t.Errorf("Workspace ID %d not found in cache", cw.ID)
			} else if cached.Name != cw.Name {
				t.Errorf("Cache mismatch for Workspace ID %d: expected Name %s, got %s", cw.ID, cw.Name, cached.Name)
			}
		}
		return
	})

	t.Run("Update", func(t *testing.T) {
		// First, create a workspace to update
		now := time.Now()
		workspaceCreate := &e_workspace.WorkspaceCreate{
			Name:             "Workspace To Update",
			PreWorkspaceID:   nil,
			Rank:             3,
			Ancient:          "Ancient C",
			Enable:           true,
			OwnerID:          1,
			ExpiryDate:       now.AddDate(3, 0, 0),
			Auth:             json.RawMessage(`{"role": "editor"}`),
			UserAuthConst:    json.RawMessage(`{"access": "partial"}`),
			UserAuthPassDown: json.RawMessage(`{"pass": "neutral"}`),
			UserAuthCustom:   json.RawMessage(`{"custom": "settings"}`),
		}

		createdWorkspaces, err := ws.Create([]*e_workspace.WorkspaceCreate{workspaceCreate})
		if err != nil {
			t.Fatalf("Setup Create method failed: %v", err)
		}
		if len(createdWorkspaces) != 1 {
			t.Fatalf("Expected 1 workspace, got %d", len(createdWorkspaces))
		}
		workspace := createdWorkspaces[0]

		// Define updates
		newName := "Workspace Updated"
		newEnable := false
		newRank := int32(4)
		newAncient := "Ancient Updated"
		newExpiredTime := now.AddDate(3, 0, 5)
		workspaceUpdate := &e_workspace.WorkspaceUpdate{
			ID:               workspace.ID,
			Name:             &newName,
			Enable:           &newEnable,
			Rank:             &newRank,
			Ancient:          &newAncient,
			OwnerID:          nil, // No change
			PreWorkspaceID:   nil, // No change
			ExpiryDate:       &newExpiredTime,
			Auth:             &workspaceCreate.Auth, // Assuming no change
			UserAuthConst:    &workspaceCreate.UserAuthConst,
			UserAuthPassDown: &workspaceCreate.UserAuthPassDown,
			UserAuthCustom:   &workspaceCreate.UserAuthCustom,
			// Assuming no changes to WUsers and WGroups
			WUsers:  nil,
			WGroups: nil,
		}

		// Call the Update method
		err = ws.Update([]*e_workspace.WorkspaceUpdate{workspaceUpdate})
		if err != nil {
			t.Fatalf("Update method failed: %v", err)
		}

		// Initialize the query builder
		q := query.Use(ws.d.GetSql())

		// Fetch the updated workspace from the database
		updatedWorkspace, err := q.Workspace.WithContext(ctx).Where(q.Workspace.ID.Eq(workspace.ID)).First()
		if err != nil {
			t.Fatalf("Failed to fetch updated workspace: %v", err)
		}

		// Verify the updates
		if updatedWorkspace.Name != newName {
			t.Errorf("Expected Name %s, got %s", newName, updatedWorkspace.Name)
		}
		if updatedWorkspace.Enable != newEnable {
			t.Errorf("Expected Enable %v, got %v", newEnable, updatedWorkspace.Enable)
		}
		if updatedWorkspace.Rank != newRank {
			t.Errorf("Expected Rank %d, got %d", newRank, updatedWorkspace.Rank)
		}
		if updatedWorkspace.Ancient != newAncient {
			t.Errorf("Expected Ancient %s, got %s", newAncient, updatedWorkspace.Ancient)
		}
		updatedYear, updatedMonth, updatedDay := updatedWorkspace.ExpiryDate.Date()
		newYear, newMonth, newDay := newExpiredTime.Date()
		if updatedYear != newYear || updatedMonth != newMonth || updatedDay != newDay {
			t.Errorf("Expected ExpiryDate %s, got %s", newExpiredTime, updatedWorkspace.ExpiryDate)
		}

		// Verify that the cache is updated
		cacheMap := ws.getCacheMap()
		if cached, exists := cacheMap[int(workspace.ID)]; !exists {
			t.Errorf("Workspace ID %d not found in cache after update", workspace.ID)
		} else {
			if cached.Name != newName {
				t.Errorf("Cache mismatch for Workspace ID %d: expected Name %s, got %s", workspace.ID, newName, cached.Name)
			}
			if cached.Enable != newEnable {
				t.Errorf("Cache mismatch for Workspace ID %d: expected Enable %v, got %v", workspace.ID, newEnable, cached.Enable)
			}
			if cached.Rank != newRank {
				t.Errorf("Cache mismatch for Workspace ID %d: expected Rank %d, got %d", workspace.ID, newRank, cached.Rank)
			}
			if cached.Ancient != newAncient {
				t.Errorf("Cache mismatch for Workspace ID %d: expected Ancient %s, got %s", workspace.ID, newAncient, cached.Ancient)
			}
		}
	})

}

func TestGetDefaultUserAuth(t *testing.T) {
	// Set up the WorkspaceServer with a mocked cacheMap
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	ws := setUpServer(ctx)

	// Populate the cacheMap with a test workspace
	cacheMap := map[int]model.Workspace{1: model.Workspace{
		ID:               1,
		UserAuthConst:    json.RawMessage(`{"constKey":"constValue"}`),
		UserAuthPassDown: json.RawMessage(`{"passDownKey":"passDownValue"}`),
		UserAuthCustom:   json.RawMessage(`{"customKey":"customValue"}`),
	}}
	ws.setCacheMap(cacheMap)

	// Call getDefaultUserAuth with a valid workspace ID
	cache := ws.getCacheMap()
	workspace := cache[1]
	result := ws.getDefaultUserAuth(workspace)

	// Expected JSON structure
	expectedMap := map[string]interface{}{
		"constKey":    "constValue",
		"passDownKey": "passDownValue",
		"customKey":   "customValue",
	}

	expectedResult, err := json.Marshal(expectedMap)
	fmt.Println(string(expectedResult))
	if err != nil {
		t.Fatalf("Error marshaling expected result: %v", err)
	}

	// Compare the result with the expected JSON
	if !jsonEqual(result, expectedResult) {
		t.Errorf("Expected %s, got %s", string(expectedResult), string(result))
	}

	// Call getDefaultUserAuth with an invalid workspace ID
	workspace2 := cache[2]
	result = ws.getDefaultUserAuth(workspace2)
	if result != nil {
		t.Errorf("Expected nil for invalid workspace ID, got %s", string(result))
	}
}

// Helper function to compare two JSON RawMessages
func jsonEqual(a, b json.RawMessage) bool {
	var objA interface{}
	var objB interface{}

	if err := json.Unmarshal(a, &objA); err != nil {
		return false
	}
	if err := json.Unmarshal(b, &objB); err != nil {
		return false
	}

	return reflect.DeepEqual(objA, objB)
}

func TestSyncUserAuth(t *testing.T) {
	// Set up context and WorkspaceServer
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	ws := setUpServer(ctx)
	defer ws.Close()

	// Clean up database before starting
	cleanUpDatabase(t, ws.d)

	// Create a workspace with specific UserAuthConst, UserAuthPassDown, UserAuthCustom
	workspaceCreate := &e_workspace.WorkspaceCreate{
		Name:             "Test Workspace",
		PreWorkspaceID:   nil,
		Rank:             1,
		Ancient:          "Test Ancient",
		Enable:           true,
		OwnerID:          1,
		ExpiryDate:       time.Now().AddDate(1, 0, 0),
		Auth:             json.RawMessage(`{"role": "admin"}`),
		UserAuthConst:    json.RawMessage(`{"constKey": "constValue"}`),
		UserAuthPassDown: json.RawMessage(`{"passDownKey": "passDownValue"}`),
		UserAuthCustom:   json.RawMessage(`{"customKey": "customValue"}`),
	}

	// Create the workspace
	createdWorkspaces, err := ws.Create([]*e_workspace.WorkspaceCreate{workspaceCreate})
	if err != nil {
		t.Fatalf("Failed to create workspace: %v", err)
	}
	if len(createdWorkspaces) != 1 {
		t.Fatalf("Expected 1 workspace created, got %d", len(createdWorkspaces))
	}
	workspace := createdWorkspaces[0]

	// Create WUsers with specific Auth and assign them to the workspace via ws.Update
	wUser1Auth := json.RawMessage(`{"constKey":"constValue","customKey":"customValue","passDownKey":"passDownValue"}`)
	wUser2Auth := json.RawMessage(`{"constKey":"constValue","customKey":"customValue","passDownKey":"passDownValue"}`)

	userID1 := int32(1)
	userID2 := int32(2)
	var auth1 json.RawMessage = []byte("{}")
	var auth2 json.RawMessage = []byte("{}")

	wUserUpdate1 := e_w_user.WUserUpdate{
		ID:          0, // New WUser
		UserID:      &userID1,
		Enable:      nil,
		Auth:        &auth1,
		WUserGroups: nil,
	}
	wUserUpdate2 := e_w_user.WUserUpdate{
		ID:          0,
		UserID:      &userID2,
		Enable:      nil,
		Auth:        &auth2,
		WUserGroups: nil,
	}

	workspaceUpdate := &e_workspace.WorkspaceUpdate{
		ID:     workspace.ID,
		WUsers: []e_w_user.WUserUpdate{wUserUpdate1, wUserUpdate2},
	}

	// Call the Update method to add WUsers
	err = ws.Update([]*e_workspace.WorkspaceUpdate{workspaceUpdate})
	if err != nil {
		t.Fatalf("Failed to update workspace with WUsers: %v", err)
	}

	// Now call syncUserGroupAuth with workspace.ID
	ws.syncUserGroupAuth(workspace.ID)

	// Fetch WUsers from the database and check their Auth fields
	q := query.Use(ws.d.GetSql())
	ctx = context.Background()

	wUsers, err := q.WUser.WithContext(ctx).Where(q.WUser.WorkspaceID.Eq(workspace.ID)).Find()
	if err != nil {
		t.Fatalf("Failed to fetch WUsers: %v", err)
	}

	// For each WUser, check that their Auth is the merged result of their initial Auth and defaultUserAuth
	defaultUserAuth := ws.getDefaultUserAuth(workspace)

	// Map of initial Auths by UserID
	initialAuthMap := map[int32]json.RawMessage{
		userID1: wUser1Auth,
		userID2: wUser2Auth,
	}

	for _, wuser := range wUsers {
		initialAuth, exists := initialAuthMap[wuser.UserID]
		if !exists {
			t.Errorf("Unexpected UserID %d", wuser.UserID)
			continue
		}

		expectedAuth, err := convert.MergeJSON(initialAuth, defaultUserAuth)
		if err != nil {
			t.Fatalf("Failed to merge JSON: %v", err)
		}

		if !jsonEqual(wuser.Auth, expectedAuth) {
			t.Errorf("WUser ID %d Auth mismatch. Expected %s, got %s", wuser.ID, string(expectedAuth), string(wuser.Auth))
		}
	}
}
