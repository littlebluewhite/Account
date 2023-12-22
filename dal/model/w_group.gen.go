// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameWGroup = "w_group"

// WGroup mapped from table <w_group>
type WGroup struct {
	ID          int32      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name        string     `gorm:"column:name;not null" json:"name"`
	CreatorID   int32      `gorm:"column:creator_id;not null" json:"creator_id"`
	WorkspaceID int32      `gorm:"column:workspace_id;not null" json:"workspace_id"`
	Enable      bool       `gorm:"column:enable;not null" json:"enable"`
	CreatedAt   *time.Time `gorm:"column:created_at;default:now()" json:"created_at"`
	users       []WUser    `gorm:"many2many:user_group" json:"users"`
}

// TableName WGroup's table name
func (*WGroup) TableName() string {
	return TableNameWGroup
}
