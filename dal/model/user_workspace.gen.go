// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"encoding/json"
	"time"
)

const TableNameUserWorkspace = "user_workspace"

// UserWorkspace mapped from table <user_workspace>
type UserWorkspace struct {
	UserID      int32           `gorm:"column:user_id;primaryKey" json:"user_id"`
	WorkspaceID int32           `gorm:"column:workspace_id;primaryKey" json:"workspace_id"`
	Enable      bool            `gorm:"column:enable;not null" json:"enable"`
	Auth        json.RawMessage `gorm:"column:auth;default:json_object()" json:"auth"`
	CreatedAt   *time.Time      `gorm:"column:created_at;default:now()" json:"created_at"`
}

// TableName UserWorkspace's table name
func (*UserWorkspace) TableName() string {
	return TableNameUserWorkspace
}
