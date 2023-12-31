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

	"account/dal/model"
)

func newWorkspace(db *gorm.DB, opts ...gen.DOOption) workspace {
	_workspace := workspace{}

	_workspace.workspaceDo.UseDB(db, opts...)
	_workspace.workspaceDo.UseModel(&model.Workspace{})

	tableName := _workspace.workspaceDo.TableName()
	_workspace.ALL = field.NewAsterisk(tableName)
	_workspace.ID = field.NewInt32(tableName, "id")
	_workspace.Name = field.NewString(tableName, "name")
	_workspace.PreWorkspace = field.NewInt32(tableName, "pre_workspace")
	_workspace.Rank = field.NewInt32(tableName, "rank")
	_workspace.Ancient = field.NewString(tableName, "ancient")
	_workspace.Enable = field.NewBool(tableName, "enable")
	_workspace.OwnerID = field.NewInt32(tableName, "owner_id")
	_workspace.ExpiryDate = field.NewTime(tableName, "expiry_date")
	_workspace.Auth = field.NewString(tableName, "auth")
	_workspace.UserAuth = field.NewString(tableName, "user_auth")
	_workspace.CreatedAt = field.NewTime(tableName, "created_at")
	_workspace.users = workspaceHasManyusers{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("users", "model.UserWorkspace"),
	}

	_workspace.groups = workspaceHasManygroups{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("groups", "model.WGroup"),
		users: struct {
			field.RelationField
			groups struct {
				field.RelationField
			}
			workspaces struct {
				field.RelationField
			}
		}{
			RelationField: field.NewRelation("groups.users", "model.WUser"),
			groups: struct {
				field.RelationField
			}{
				RelationField: field.NewRelation("groups.users.groups", "model.UserGroup"),
			},
			workspaces: struct {
				field.RelationField
			}{
				RelationField: field.NewRelation("groups.users.workspaces", "model.UserWorkspace"),
			},
		},
	}

	_workspace.fillFieldMap()

	return _workspace
}

type workspace struct {
	workspaceDo workspaceDo

	ALL          field.Asterisk
	ID           field.Int32
	Name         field.String
	PreWorkspace field.Int32
	Rank         field.Int32
	Ancient      field.String
	Enable       field.Bool
	OwnerID      field.Int32
	ExpiryDate   field.Time
	Auth         field.String
	UserAuth     field.String
	CreatedAt    field.Time
	users        workspaceHasManyusers

	groups workspaceHasManygroups

	fieldMap map[string]field.Expr
}

func (w workspace) Table(newTableName string) *workspace {
	w.workspaceDo.UseTable(newTableName)
	return w.updateTableName(newTableName)
}

func (w workspace) As(alias string) *workspace {
	w.workspaceDo.DO = *(w.workspaceDo.As(alias).(*gen.DO))
	return w.updateTableName(alias)
}

func (w *workspace) updateTableName(table string) *workspace {
	w.ALL = field.NewAsterisk(table)
	w.ID = field.NewInt32(table, "id")
	w.Name = field.NewString(table, "name")
	w.PreWorkspace = field.NewInt32(table, "pre_workspace")
	w.Rank = field.NewInt32(table, "rank")
	w.Ancient = field.NewString(table, "ancient")
	w.Enable = field.NewBool(table, "enable")
	w.OwnerID = field.NewInt32(table, "owner_id")
	w.ExpiryDate = field.NewTime(table, "expiry_date")
	w.Auth = field.NewString(table, "auth")
	w.UserAuth = field.NewString(table, "user_auth")
	w.CreatedAt = field.NewTime(table, "created_at")

	w.fillFieldMap()

	return w
}

func (w *workspace) WithContext(ctx context.Context) *workspaceDo {
	return w.workspaceDo.WithContext(ctx)
}

func (w workspace) TableName() string { return w.workspaceDo.TableName() }

func (w workspace) Alias() string { return w.workspaceDo.Alias() }

func (w workspace) Columns(cols ...field.Expr) gen.Columns { return w.workspaceDo.Columns(cols...) }

func (w *workspace) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := w.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (w *workspace) fillFieldMap() {
	w.fieldMap = make(map[string]field.Expr, 13)
	w.fieldMap["id"] = w.ID
	w.fieldMap["name"] = w.Name
	w.fieldMap["pre_workspace"] = w.PreWorkspace
	w.fieldMap["rank"] = w.Rank
	w.fieldMap["ancient"] = w.Ancient
	w.fieldMap["enable"] = w.Enable
	w.fieldMap["owner_id"] = w.OwnerID
	w.fieldMap["expiry_date"] = w.ExpiryDate
	w.fieldMap["auth"] = w.Auth
	w.fieldMap["user_auth"] = w.UserAuth
	w.fieldMap["created_at"] = w.CreatedAt

}

func (w workspace) clone(db *gorm.DB) workspace {
	w.workspaceDo.ReplaceConnPool(db.Statement.ConnPool)
	return w
}

func (w workspace) replaceDB(db *gorm.DB) workspace {
	w.workspaceDo.ReplaceDB(db)
	return w
}

type workspaceHasManyusers struct {
	db *gorm.DB

	field.RelationField
}

func (a workspaceHasManyusers) Where(conds ...field.Expr) *workspaceHasManyusers {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a workspaceHasManyusers) WithContext(ctx context.Context) *workspaceHasManyusers {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a workspaceHasManyusers) Session(session *gorm.Session) *workspaceHasManyusers {
	a.db = a.db.Session(session)
	return &a
}

func (a workspaceHasManyusers) Model(m *model.Workspace) *workspaceHasManyusersTx {
	return &workspaceHasManyusersTx{a.db.Model(m).Association(a.Name())}
}

type workspaceHasManyusersTx struct{ tx *gorm.Association }

func (a workspaceHasManyusersTx) Find() (result []*model.UserWorkspace, err error) {
	return result, a.tx.Find(&result)
}

func (a workspaceHasManyusersTx) Append(values ...*model.UserWorkspace) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a workspaceHasManyusersTx) Replace(values ...*model.UserWorkspace) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a workspaceHasManyusersTx) Delete(values ...*model.UserWorkspace) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a workspaceHasManyusersTx) Clear() error {
	return a.tx.Clear()
}

func (a workspaceHasManyusersTx) Count() int64 {
	return a.tx.Count()
}

type workspaceHasManygroups struct {
	db *gorm.DB

	field.RelationField

	users struct {
		field.RelationField
		groups struct {
			field.RelationField
		}
		workspaces struct {
			field.RelationField
		}
	}
}

func (a workspaceHasManygroups) Where(conds ...field.Expr) *workspaceHasManygroups {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a workspaceHasManygroups) WithContext(ctx context.Context) *workspaceHasManygroups {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a workspaceHasManygroups) Session(session *gorm.Session) *workspaceHasManygroups {
	a.db = a.db.Session(session)
	return &a
}

func (a workspaceHasManygroups) Model(m *model.Workspace) *workspaceHasManygroupsTx {
	return &workspaceHasManygroupsTx{a.db.Model(m).Association(a.Name())}
}

type workspaceHasManygroupsTx struct{ tx *gorm.Association }

func (a workspaceHasManygroupsTx) Find() (result []*model.WGroup, err error) {
	return result, a.tx.Find(&result)
}

func (a workspaceHasManygroupsTx) Append(values ...*model.WGroup) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a workspaceHasManygroupsTx) Replace(values ...*model.WGroup) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a workspaceHasManygroupsTx) Delete(values ...*model.WGroup) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a workspaceHasManygroupsTx) Clear() error {
	return a.tx.Clear()
}

func (a workspaceHasManygroupsTx) Count() int64 {
	return a.tx.Count()
}

type workspaceDo struct{ gen.DO }

func (w workspaceDo) Debug() *workspaceDo {
	return w.withDO(w.DO.Debug())
}

func (w workspaceDo) WithContext(ctx context.Context) *workspaceDo {
	return w.withDO(w.DO.WithContext(ctx))
}

func (w workspaceDo) ReadDB() *workspaceDo {
	return w.Clauses(dbresolver.Read)
}

func (w workspaceDo) WriteDB() *workspaceDo {
	return w.Clauses(dbresolver.Write)
}

func (w workspaceDo) Session(config *gorm.Session) *workspaceDo {
	return w.withDO(w.DO.Session(config))
}

func (w workspaceDo) Clauses(conds ...clause.Expression) *workspaceDo {
	return w.withDO(w.DO.Clauses(conds...))
}

func (w workspaceDo) Returning(value interface{}, columns ...string) *workspaceDo {
	return w.withDO(w.DO.Returning(value, columns...))
}

func (w workspaceDo) Not(conds ...gen.Condition) *workspaceDo {
	return w.withDO(w.DO.Not(conds...))
}

func (w workspaceDo) Or(conds ...gen.Condition) *workspaceDo {
	return w.withDO(w.DO.Or(conds...))
}

func (w workspaceDo) Select(conds ...field.Expr) *workspaceDo {
	return w.withDO(w.DO.Select(conds...))
}

func (w workspaceDo) Where(conds ...gen.Condition) *workspaceDo {
	return w.withDO(w.DO.Where(conds...))
}

func (w workspaceDo) Order(conds ...field.Expr) *workspaceDo {
	return w.withDO(w.DO.Order(conds...))
}

func (w workspaceDo) Distinct(cols ...field.Expr) *workspaceDo {
	return w.withDO(w.DO.Distinct(cols...))
}

func (w workspaceDo) Omit(cols ...field.Expr) *workspaceDo {
	return w.withDO(w.DO.Omit(cols...))
}

func (w workspaceDo) Join(table schema.Tabler, on ...field.Expr) *workspaceDo {
	return w.withDO(w.DO.Join(table, on...))
}

func (w workspaceDo) LeftJoin(table schema.Tabler, on ...field.Expr) *workspaceDo {
	return w.withDO(w.DO.LeftJoin(table, on...))
}

func (w workspaceDo) RightJoin(table schema.Tabler, on ...field.Expr) *workspaceDo {
	return w.withDO(w.DO.RightJoin(table, on...))
}

func (w workspaceDo) Group(cols ...field.Expr) *workspaceDo {
	return w.withDO(w.DO.Group(cols...))
}

func (w workspaceDo) Having(conds ...gen.Condition) *workspaceDo {
	return w.withDO(w.DO.Having(conds...))
}

func (w workspaceDo) Limit(limit int) *workspaceDo {
	return w.withDO(w.DO.Limit(limit))
}

func (w workspaceDo) Offset(offset int) *workspaceDo {
	return w.withDO(w.DO.Offset(offset))
}

func (w workspaceDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *workspaceDo {
	return w.withDO(w.DO.Scopes(funcs...))
}

func (w workspaceDo) Unscoped() *workspaceDo {
	return w.withDO(w.DO.Unscoped())
}

func (w workspaceDo) Create(values ...*model.Workspace) error {
	if len(values) == 0 {
		return nil
	}
	return w.DO.Create(values)
}

func (w workspaceDo) CreateInBatches(values []*model.Workspace, batchSize int) error {
	return w.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (w workspaceDo) Save(values ...*model.Workspace) error {
	if len(values) == 0 {
		return nil
	}
	return w.DO.Save(values)
}

func (w workspaceDo) First() (*model.Workspace, error) {
	if result, err := w.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Workspace), nil
	}
}

func (w workspaceDo) Take() (*model.Workspace, error) {
	if result, err := w.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Workspace), nil
	}
}

func (w workspaceDo) Last() (*model.Workspace, error) {
	if result, err := w.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Workspace), nil
	}
}

func (w workspaceDo) Find() ([]*model.Workspace, error) {
	result, err := w.DO.Find()
	return result.([]*model.Workspace), err
}

func (w workspaceDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Workspace, err error) {
	buf := make([]*model.Workspace, 0, batchSize)
	err = w.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (w workspaceDo) FindInBatches(result *[]*model.Workspace, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return w.DO.FindInBatches(result, batchSize, fc)
}

func (w workspaceDo) Attrs(attrs ...field.AssignExpr) *workspaceDo {
	return w.withDO(w.DO.Attrs(attrs...))
}

func (w workspaceDo) Assign(attrs ...field.AssignExpr) *workspaceDo {
	return w.withDO(w.DO.Assign(attrs...))
}

func (w workspaceDo) Joins(fields ...field.RelationField) *workspaceDo {
	for _, _f := range fields {
		w = *w.withDO(w.DO.Joins(_f))
	}
	return &w
}

func (w workspaceDo) Preload(fields ...field.RelationField) *workspaceDo {
	for _, _f := range fields {
		w = *w.withDO(w.DO.Preload(_f))
	}
	return &w
}

func (w workspaceDo) FirstOrInit() (*model.Workspace, error) {
	if result, err := w.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Workspace), nil
	}
}

func (w workspaceDo) FirstOrCreate() (*model.Workspace, error) {
	if result, err := w.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Workspace), nil
	}
}

func (w workspaceDo) FindByPage(offset int, limit int) (result []*model.Workspace, count int64, err error) {
	result, err = w.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = w.Offset(-1).Limit(-1).Count()
	return
}

func (w workspaceDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = w.Count()
	if err != nil {
		return
	}

	err = w.Offset(offset).Limit(limit).Scan(result)
	return
}

func (w workspaceDo) Scan(result interface{}) (err error) {
	return w.DO.Scan(result)
}

func (w workspaceDo) Delete(models ...*model.Workspace) (result gen.ResultInfo, err error) {
	return w.DO.Delete(models)
}

func (w *workspaceDo) withDO(do gen.Dao) *workspaceDo {
	w.DO = *do.(*gen.DO)
	return w
}
