// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"encoding/json"
	"time"
)

const TableNameWorkspace = "workspace"

// Workspace mapped from table <workspace>
type Workspace struct {
	ID               int32           `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name             string          `gorm:"column:name;not null" json:"name"`
	PreWorkspaceID   *int32          `gorm:"column:pre_workspace_id" json:"pre_workspace_id"`
	Rank             int32           `gorm:"column:rank;not null" json:"rank"`
	Ancient          string          `gorm:"column:ancient;not null" json:"ancient"`
	Enable           bool            `gorm:"column:enable;not null" json:"enable"`
	OwnerID          int32           `gorm:"column:owner_id;not null" json:"owner_id"`
	ExpiryDate       time.Time       `gorm:"column:expiry_date;not null" json:"expiry_date"`
	Auth             json.RawMessage `gorm:"column:auth;default:json_object()" json:"auth"`
	UserAuthConst    json.RawMessage `gorm:"column:user_auth_const;default:json_object()" json:"user_auth_const"`
	UserAuthPassDown json.RawMessage `gorm:"column:user_auth_pass_down;default:json_object()" json:"user_auth_pass_down"`
	UserAuthCustom   json.RawMessage `gorm:"column:user_auth_custom;default:json_object()" json:"user_auth_custom"`
	UpdatedAt        *time.Time      `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedAt        *time.Time      `gorm:"column:created_at;default:now()" json:"created_at"`
	WUsers           []WUser         `gorm:"foreignKey:workspace_id" json:"w_users"`
	WGroups          []WGroup        `gorm:"foreignKey:workspace_id" json:"w_groups"`
	NextWorkspaces   []Workspace     `gorm:"foreignKey:pre_workspace_id" json:"next_workspaces"`
}

// TableName Workspace's table name
func (*Workspace) TableName() string {
	return TableNameWorkspace
}
