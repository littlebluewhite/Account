package workspace_server

import (
	"context"
	"github.com/littlebluewhite/Account/dal/model"
	"github.com/littlebluewhite/Account/dal/query"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

type IWorkspaceRepository interface {
	List(ctx context.Context) ([]*model.Workspace, error)
	Find(ctx context.Context, ids []int32) ([]*model.Workspace, error)
	CreateBatch(ctx context.Context, workspaces []*model.Workspace) error
	UpdateBatch(ctx context.Context, workspaces []map[string]interface{}) error
	DeleteWUsers(ctx context.Context, ids []int32) error
	DeleteWGroups(ctx context.Context, ids []int32) error
	DeleteWUserGroups(ctx context.Context, ids []int32) error
}

type WorkspaceRepository struct {
	db *gorm.DB
}

func NewWorkspaceRepository(db *gorm.DB) *WorkspaceRepository {
	return &WorkspaceRepository{db: db}
}

func (r *WorkspaceRepository) List(ctx context.Context) ([]*model.Workspace, error) {
	ws := query.Use(r.db).Workspace
	return ws.WithContext(ctx).Preload(field.Associations).
		Preload(ws.WGroups.WUserGroups).
		Preload(ws.WUsers.WUserGroups).
		Find()
}

func (r *WorkspaceRepository) Find(ctx context.Context, ids []int32) ([]*model.Workspace, error) {
	ws := query.Use(r.db).Workspace
	return ws.WithContext(ctx).Preload(field.Associations).
		Where(ws.ID.In(ids...)).Find()
}

func (r *WorkspaceRepository) CreateBatch(ctx context.Context, workspaces []*model.Workspace) error {
	return query.Use(r.db).Workspace.WithContext(ctx).CreateInBatches(workspaces, 100)
}

func (r *WorkspaceRepository) UpdateBatch(ctx context.Context, workspaces []map[string]interface{}) error {
	for _, workspace := range workspaces {
		if _, err := query.Use(r.db).Workspace.WithContext(ctx).
			Where(query.Use(r.db).Workspace.ID.Eq(workspace["id"].(int32))).
			Updates(workspace); err != nil {
			return err
		}
	}
	return nil
}

func (r *WorkspaceRepository) DeleteWUsers(ctx context.Context, ids []int32) error {
	_, err := query.Use(r.db).WUser.WithContext(ctx).
		Where(query.Use(r.db).WUser.ID.In(ids...)).Delete()
	return err
}

func (r *WorkspaceRepository) DeleteWGroups(ctx context.Context, ids []int32) error {
	_, err := query.Use(r.db).WGroup.WithContext(ctx).
		Where(query.Use(r.db).WGroup.ID.In(ids...)).Delete()
	return err
}

func (r *WorkspaceRepository) DeleteWUserGroups(ctx context.Context, ids []int32) error {
	_, err := query.Use(r.db).WUserGroup.WithContext(ctx).
		Where(query.Use(r.db).WUserGroup.WUserID.In(ids...)).Delete()
	return err
}
