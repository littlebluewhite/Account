// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/littlebluewhite/Account/dal/model"
)

func newUserWorkspace(db *gorm.DB, opts ...gen.DOOption) userWorkspace {
	_userWorkspace := userWorkspace{}

	_userWorkspace.userWorkspaceDo.UseDB(db, opts...)
	_userWorkspace.userWorkspaceDo.UseModel(&model.UserWorkspace{})

	tableName := _userWorkspace.userWorkspaceDo.TableName()
	_userWorkspace.ALL = field.NewAsterisk(tableName)
	_userWorkspace.UserID = field.NewInt32(tableName, "user_id")
	_userWorkspace.WorkspaceID = field.NewInt32(tableName, "workspace_id")
	_userWorkspace.Enable = field.NewBool(tableName, "enable")
	_userWorkspace.Auth = field.NewBytes(tableName, "auth")
	_userWorkspace.CreatedAt = field.NewTime(tableName, "created_at")

	_userWorkspace.fillFieldMap()

	return _userWorkspace
}

type userWorkspace struct {
	userWorkspaceDo userWorkspaceDo

	ALL         field.Asterisk
	UserID      field.Int32
	WorkspaceID field.Int32
	Enable      field.Bool
	Auth        field.Bytes
	CreatedAt   field.Time

	fieldMap map[string]field.Expr
}

func (u userWorkspace) Table(newTableName string) *userWorkspace {
	u.userWorkspaceDo.UseTable(newTableName)
	return u.updateTableName(newTableName)
}

func (u userWorkspace) As(alias string) *userWorkspace {
	u.userWorkspaceDo.DO = *(u.userWorkspaceDo.As(alias).(*gen.DO))
	return u.updateTableName(alias)
}

func (u *userWorkspace) updateTableName(table string) *userWorkspace {
	u.ALL = field.NewAsterisk(table)
	u.UserID = field.NewInt32(table, "user_id")
	u.WorkspaceID = field.NewInt32(table, "workspace_id")
	u.Enable = field.NewBool(table, "enable")
	u.Auth = field.NewBytes(table, "auth")
	u.CreatedAt = field.NewTime(table, "created_at")

	u.fillFieldMap()

	return u
}

func (u *userWorkspace) WithContext(ctx context.Context) *userWorkspaceDo {
	return u.userWorkspaceDo.WithContext(ctx)
}

func (u userWorkspace) TableName() string { return u.userWorkspaceDo.TableName() }

func (u userWorkspace) Alias() string { return u.userWorkspaceDo.Alias() }

func (u userWorkspace) Columns(cols ...field.Expr) gen.Columns {
	return u.userWorkspaceDo.Columns(cols...)
}

func (u *userWorkspace) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := u.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (u *userWorkspace) fillFieldMap() {
	u.fieldMap = make(map[string]field.Expr, 5)
	u.fieldMap["user_id"] = u.UserID
	u.fieldMap["workspace_id"] = u.WorkspaceID
	u.fieldMap["enable"] = u.Enable
	u.fieldMap["auth"] = u.Auth
	u.fieldMap["created_at"] = u.CreatedAt
}

func (u userWorkspace) clone(db *gorm.DB) userWorkspace {
	u.userWorkspaceDo.ReplaceConnPool(db.Statement.ConnPool)
	return u
}

func (u userWorkspace) replaceDB(db *gorm.DB) userWorkspace {
	u.userWorkspaceDo.ReplaceDB(db)
	return u
}

type userWorkspaceDo struct{ gen.DO }

func (u userWorkspaceDo) Debug() *userWorkspaceDo {
	return u.withDO(u.DO.Debug())
}

func (u userWorkspaceDo) WithContext(ctx context.Context) *userWorkspaceDo {
	return u.withDO(u.DO.WithContext(ctx))
}

func (u userWorkspaceDo) ReadDB() *userWorkspaceDo {
	return u.Clauses(dbresolver.Read)
}

func (u userWorkspaceDo) WriteDB() *userWorkspaceDo {
	return u.Clauses(dbresolver.Write)
}

func (u userWorkspaceDo) Session(config *gorm.Session) *userWorkspaceDo {
	return u.withDO(u.DO.Session(config))
}

func (u userWorkspaceDo) Clauses(conds ...clause.Expression) *userWorkspaceDo {
	return u.withDO(u.DO.Clauses(conds...))
}

func (u userWorkspaceDo) Returning(value interface{}, columns ...string) *userWorkspaceDo {
	return u.withDO(u.DO.Returning(value, columns...))
}

func (u userWorkspaceDo) Not(conds ...gen.Condition) *userWorkspaceDo {
	return u.withDO(u.DO.Not(conds...))
}

func (u userWorkspaceDo) Or(conds ...gen.Condition) *userWorkspaceDo {
	return u.withDO(u.DO.Or(conds...))
}

func (u userWorkspaceDo) Select(conds ...field.Expr) *userWorkspaceDo {
	return u.withDO(u.DO.Select(conds...))
}

func (u userWorkspaceDo) Where(conds ...gen.Condition) *userWorkspaceDo {
	return u.withDO(u.DO.Where(conds...))
}

func (u userWorkspaceDo) Order(conds ...field.Expr) *userWorkspaceDo {
	return u.withDO(u.DO.Order(conds...))
}

func (u userWorkspaceDo) Distinct(cols ...field.Expr) *userWorkspaceDo {
	return u.withDO(u.DO.Distinct(cols...))
}

func (u userWorkspaceDo) Omit(cols ...field.Expr) *userWorkspaceDo {
	return u.withDO(u.DO.Omit(cols...))
}

func (u userWorkspaceDo) Join(table schema.Tabler, on ...field.Expr) *userWorkspaceDo {
	return u.withDO(u.DO.Join(table, on...))
}

func (u userWorkspaceDo) LeftJoin(table schema.Tabler, on ...field.Expr) *userWorkspaceDo {
	return u.withDO(u.DO.LeftJoin(table, on...))
}

func (u userWorkspaceDo) RightJoin(table schema.Tabler, on ...field.Expr) *userWorkspaceDo {
	return u.withDO(u.DO.RightJoin(table, on...))
}

func (u userWorkspaceDo) Group(cols ...field.Expr) *userWorkspaceDo {
	return u.withDO(u.DO.Group(cols...))
}

func (u userWorkspaceDo) Having(conds ...gen.Condition) *userWorkspaceDo {
	return u.withDO(u.DO.Having(conds...))
}

func (u userWorkspaceDo) Limit(limit int) *userWorkspaceDo {
	return u.withDO(u.DO.Limit(limit))
}

func (u userWorkspaceDo) Offset(offset int) *userWorkspaceDo {
	return u.withDO(u.DO.Offset(offset))
}

func (u userWorkspaceDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *userWorkspaceDo {
	return u.withDO(u.DO.Scopes(funcs...))
}

func (u userWorkspaceDo) Unscoped() *userWorkspaceDo {
	return u.withDO(u.DO.Unscoped())
}

func (u userWorkspaceDo) Create(values ...*model.UserWorkspace) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Create(values)
}

func (u userWorkspaceDo) CreateInBatches(values []*model.UserWorkspace, batchSize int) error {
	return u.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (u userWorkspaceDo) Save(values ...*model.UserWorkspace) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Save(values)
}

func (u userWorkspaceDo) First() (*model.UserWorkspace, error) {
	if result, err := u.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserWorkspace), nil
	}
}

func (u userWorkspaceDo) Take() (*model.UserWorkspace, error) {
	if result, err := u.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserWorkspace), nil
	}
}

func (u userWorkspaceDo) Last() (*model.UserWorkspace, error) {
	if result, err := u.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserWorkspace), nil
	}
}

func (u userWorkspaceDo) Find() ([]*model.UserWorkspace, error) {
	result, err := u.DO.Find()
	return result.([]*model.UserWorkspace), err
}

func (u userWorkspaceDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.UserWorkspace, err error) {
	buf := make([]*model.UserWorkspace, 0, batchSize)
	err = u.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (u userWorkspaceDo) FindInBatches(result *[]*model.UserWorkspace, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return u.DO.FindInBatches(result, batchSize, fc)
}

func (u userWorkspaceDo) Attrs(attrs ...field.AssignExpr) *userWorkspaceDo {
	return u.withDO(u.DO.Attrs(attrs...))
}

func (u userWorkspaceDo) Assign(attrs ...field.AssignExpr) *userWorkspaceDo {
	return u.withDO(u.DO.Assign(attrs...))
}

func (u userWorkspaceDo) Joins(fields ...field.RelationField) *userWorkspaceDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Joins(_f))
	}
	return &u
}

func (u userWorkspaceDo) Preload(fields ...field.RelationField) *userWorkspaceDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Preload(_f))
	}
	return &u
}

func (u userWorkspaceDo) FirstOrInit() (*model.UserWorkspace, error) {
	if result, err := u.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserWorkspace), nil
	}
}

func (u userWorkspaceDo) FirstOrCreate() (*model.UserWorkspace, error) {
	if result, err := u.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserWorkspace), nil
	}
}

func (u userWorkspaceDo) FindByPage(offset int, limit int) (result []*model.UserWorkspace, count int64, err error) {
	result, err = u.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = u.Offset(-1).Limit(-1).Count()
	return
}

func (u userWorkspaceDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = u.Count()
	if err != nil {
		return
	}

	err = u.Offset(offset).Limit(limit).Scan(result)
	return
}

func (u userWorkspaceDo) Scan(result interface{}) (err error) {
	return u.DO.Scan(result)
}

func (u userWorkspaceDo) Delete(models ...*model.UserWorkspace) (result gen.ResultInfo, err error) {
	return u.DO.Delete(models)
}

func (u *userWorkspaceDo) withDO(do gen.Dao) *userWorkspaceDo {
	u.DO = *do.(*gen.DO)
	return u
}
