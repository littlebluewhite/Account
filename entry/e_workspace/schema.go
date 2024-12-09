package e_workspace

import (
	"github.com/goccy/go-json"
	"github.com/littlebluewhite/Account/entry/e_w_group"
	"github.com/littlebluewhite/Account/entry/e_w_user"
	"time"
)

type Workspace struct {
	ID               int32              `json:"id"`
	Name             string             `json:"name"`
	PreWorkspaceID   *int32             `json:"pre_workspace_id"`
	Rank             int32              `json:"rank"`
	Ancient          string             `json:"ancient"`
	Enable           bool               `json:"enable"`
	OwnerID          int32              `json:"owner_id"`
	ExpiryDate       time.Time          `json:"expiry_date"`
	Auth             json.RawMessage    `json:"auth"`
	UserAuthConst    json.RawMessage    `json:"user_auth_const"`
	UserAuthPassDown json.RawMessage    `json:"user_auth_pass_down"`
	UserAuthCustom   json.RawMessage    `json:"user_auth_custom"`
	UpdatedAt        *time.Time         `json:"updated_at"`
	CreatedAt        *time.Time         `json:"created_at"`
	WUsers           []e_w_user.WUser   `json:"w_users"`
	WGroups          []e_w_group.WGroup `json:"w_groups"`
	NextWorkspaces   []Workspace        `json:"next_workspaces"`
}

type WorkspaceCreate struct {
	Name             string          `json:"name"`
	PreWorkspaceID   *int32          `json:"pre_workspace_id"`
	Rank             int32           `json:"rank"`
	Ancient          string          `json:"ancient"`
	Enable           bool            `json:"enable"`
	OwnerID          int32           `json:"owner_id"`
	ExpiryDate       time.Time       `json:"expiry_date"`
	Auth             json.RawMessage `json:"auth"`
	UserAuthConst    json.RawMessage `json:"user_auth_const"`
	UserAuthPassDown json.RawMessage `json:"user_auth_pass_down"`
	UserAuthCustom   json.RawMessage `json:"user_auth_custom"`
}

type WorkspaceUpdate struct {
	ID               int32                    `json:"id"`
	Name             *string                  `json:"name"`
	PreWorkspaceID   *int32                   `json:"pre_workspace_id"`
	Rank             *int32                   `json:"rank"`
	Ancient          *string                  `json:"ancient"`
	Enable           *bool                    `json:"enable"`
	OwnerID          *int32                   `json:"owner_id"`
	ExpiryDate       *time.Time               `json:"expiry_date"`
	Auth             *json.RawMessage         `json:"auth"`
	UserAuthConst    *json.RawMessage         `json:"user_auth_const"`
	UserAuthPassDown *json.RawMessage         `json:"user_auth_pass_down"`
	UserAuthCustom   *json.RawMessage         `json:"user_auth_custom"`
	WUsers           []e_w_user.WUserUpdate   `json:"w_users"`
	WGroups          []e_w_group.WGroupUpdate `json:"w_groups"`
}

func (wsu *WorkspaceUpdate) GetKey(key string) int {
	switch key {
	case "id":
		return int(wsu.ID)
	default:
		return 0
	}
}
