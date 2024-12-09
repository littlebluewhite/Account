package e_w_user

import (
	"github.com/goccy/go-json"
	"time"
)

type WUser struct {
	ID          int32           `json:"id"`
	UserID      int32           `json:"user_id"`
	WorkspaceID int32           `json:"workspace_id"`
	Enable      bool            `json:"enable"`
	Auth        json.RawMessage `json:"auth"`
	UpdatedAt   *time.Time      `json:"updated_at"`
	CreatedAt   *time.Time      `json:"created_at"`
	WUserGroups []WUserGroup    `json:"w_user_groups"`
}

type WUserGroup struct {
	WUserID  int32 `json:"w_user_id"`
	WGroupID int32 `json:"w_group_id"`
}

type WUserUpdate struct {
	ID          int32              `json:"id"`
	UserID      *int32             `json:"user_id"`
	Enable      *bool              `json:"enable"`
	Auth        *json.RawMessage   `json:"auth"`
	WUserGroups []WUserGroupUpdate `json:"w_user_groups"`
}

type WUserGroupUpdate struct {
	WGroupID int32 `json:"w_group_id"`
}
