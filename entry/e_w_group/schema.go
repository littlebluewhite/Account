package e_w_group

import (
	"github.com/goccy/go-json"
	"time"
)

type WGroup struct {
	ID          int32           `json:"id"`
	Name        string          `json:"name"`
	CreatorID   int32           `json:"creator_id"`
	WorkspaceID int32           `json:"workspace_id"`
	Enable      bool            `json:"enable"`
	DefaultAuth json.RawMessage `json:"default_auth"`
	UpdatedAt   *time.Time      `json:"updated_at"`
	CreatedAt   *time.Time      `json:"created_at"`
	WUserGroups []WUserGroup    `json:"w_user_groups"`
}

type WUserGroup struct {
	WUserID  int32 `json:"w_user_id"`
	WGroupID int32 `json:"w_group_id"`
}

type WGroupUpdate struct {
	ID          int32              `json:"id"`
	Name        *string            `json:"name"`
	CreatorID   int32              `json:"creator_id"`
	WorkspaceID int32              `json:"workspace_id"`
	DefaultAuth *json.RawMessage   `json:"default_auth"`
	Enable      *bool              `json:"enable"`
	WUserGroups []WUserGroupUpdate `json:"w_user_groups"`
}

type WUserGroupUpdate struct {
	WUserID int32 `json:"w_user_id"`
}
