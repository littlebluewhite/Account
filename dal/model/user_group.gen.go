// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameUserGroup = "user_group"

// UserGroup mapped from table <user_group>
type UserGroup struct {
	UserID  int32 `gorm:"column:user_id;primaryKey" json:"user_id"`
	GroupID int32 `gorm:"column:group_id;primaryKey" json:"group_id"`
}

// TableName UserGroup's table name
func (*UserGroup) TableName() string {
	return TableNameUserGroup
}