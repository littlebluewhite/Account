// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"gorm.io/gen"

	"gorm.io/plugin/dbresolver"
)

func Use(db *gorm.DB, opts ...gen.DOOption) *Query {
	return &Query{
		db:            db,
		DefaultAuth:   newDefaultAuth(db, opts...),
		UserGroup:     newUserGroup(db, opts...),
		UserWorkspace: newUserWorkspace(db, opts...),
		WGroup:        newWGroup(db, opts...),
		WUser:         newWUser(db, opts...),
		Workspace:     newWorkspace(db, opts...),
	}
}

type Query struct {
	db *gorm.DB

	DefaultAuth   defaultAuth
	UserGroup     userGroup
	UserWorkspace userWorkspace
	WGroup        wGroup
	WUser         wUser
	Workspace     workspace
}

func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db:            db,
		DefaultAuth:   q.DefaultAuth.clone(db),
		UserGroup:     q.UserGroup.clone(db),
		UserWorkspace: q.UserWorkspace.clone(db),
		WGroup:        q.WGroup.clone(db),
		WUser:         q.WUser.clone(db),
		Workspace:     q.Workspace.clone(db),
	}
}

func (q *Query) ReadDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Read))
}

func (q *Query) WriteDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Write))
}

func (q *Query) ReplaceDB(db *gorm.DB) *Query {
	return &Query{
		db:            db,
		DefaultAuth:   q.DefaultAuth.replaceDB(db),
		UserGroup:     q.UserGroup.replaceDB(db),
		UserWorkspace: q.UserWorkspace.replaceDB(db),
		WGroup:        q.WGroup.replaceDB(db),
		WUser:         q.WUser.replaceDB(db),
		Workspace:     q.Workspace.replaceDB(db),
	}
}

type queryCtx struct {
	DefaultAuth   *defaultAuthDo
	UserGroup     *userGroupDo
	UserWorkspace *userWorkspaceDo
	WGroup        *wGroupDo
	WUser         *wUserDo
	Workspace     *workspaceDo
}

func (q *Query) WithContext(ctx context.Context) *queryCtx {
	return &queryCtx{
		DefaultAuth:   q.DefaultAuth.WithContext(ctx),
		UserGroup:     q.UserGroup.WithContext(ctx),
		UserWorkspace: q.UserWorkspace.WithContext(ctx),
		WGroup:        q.WGroup.WithContext(ctx),
		WUser:         q.WUser.WithContext(ctx),
		Workspace:     q.Workspace.WithContext(ctx),
	}
}

func (q *Query) Transaction(fc func(tx *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *QueryTx {
	tx := q.db.Begin(opts...)
	return &QueryTx{Query: q.clone(tx), Error: tx.Error}
}

type QueryTx struct {
	*Query
	Error error
}

func (q *QueryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *QueryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *QueryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *QueryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}
