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

func newWUserGroup(db *gorm.DB, opts ...gen.DOOption) wUserGroup {
	_wUserGroup := wUserGroup{}

	_wUserGroup.wUserGroupDo.UseDB(db, opts...)
	_wUserGroup.wUserGroupDo.UseModel(&model.WUserGroup{})

	tableName := _wUserGroup.wUserGroupDo.TableName()
	_wUserGroup.ALL = field.NewAsterisk(tableName)
	_wUserGroup.WUserID = field.NewInt32(tableName, "w_user_id")
	_wUserGroup.WGroupID = field.NewInt32(tableName, "w_group_id")

	_wUserGroup.fillFieldMap()

	return _wUserGroup
}

type wUserGroup struct {
	wUserGroupDo wUserGroupDo

	ALL      field.Asterisk
	WUserID  field.Int32
	WGroupID field.Int32

	fieldMap map[string]field.Expr
}

func (w wUserGroup) Table(newTableName string) *wUserGroup {
	w.wUserGroupDo.UseTable(newTableName)
	return w.updateTableName(newTableName)
}

func (w wUserGroup) As(alias string) *wUserGroup {
	w.wUserGroupDo.DO = *(w.wUserGroupDo.As(alias).(*gen.DO))
	return w.updateTableName(alias)
}

func (w *wUserGroup) updateTableName(table string) *wUserGroup {
	w.ALL = field.NewAsterisk(table)
	w.WUserID = field.NewInt32(table, "w_user_id")
	w.WGroupID = field.NewInt32(table, "w_group_id")

	w.fillFieldMap()

	return w
}

func (w *wUserGroup) WithContext(ctx context.Context) *wUserGroupDo {
	return w.wUserGroupDo.WithContext(ctx)
}

func (w wUserGroup) TableName() string { return w.wUserGroupDo.TableName() }

func (w wUserGroup) Alias() string { return w.wUserGroupDo.Alias() }

func (w wUserGroup) Columns(cols ...field.Expr) gen.Columns { return w.wUserGroupDo.Columns(cols...) }

func (w *wUserGroup) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := w.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (w *wUserGroup) fillFieldMap() {
	w.fieldMap = make(map[string]field.Expr, 2)
	w.fieldMap["w_user_id"] = w.WUserID
	w.fieldMap["w_group_id"] = w.WGroupID
}

func (w wUserGroup) clone(db *gorm.DB) wUserGroup {
	w.wUserGroupDo.ReplaceConnPool(db.Statement.ConnPool)
	return w
}

func (w wUserGroup) replaceDB(db *gorm.DB) wUserGroup {
	w.wUserGroupDo.ReplaceDB(db)
	return w
}

type wUserGroupDo struct{ gen.DO }

func (w wUserGroupDo) Debug() *wUserGroupDo {
	return w.withDO(w.DO.Debug())
}

func (w wUserGroupDo) WithContext(ctx context.Context) *wUserGroupDo {
	return w.withDO(w.DO.WithContext(ctx))
}

func (w wUserGroupDo) ReadDB() *wUserGroupDo {
	return w.Clauses(dbresolver.Read)
}

func (w wUserGroupDo) WriteDB() *wUserGroupDo {
	return w.Clauses(dbresolver.Write)
}

func (w wUserGroupDo) Session(config *gorm.Session) *wUserGroupDo {
	return w.withDO(w.DO.Session(config))
}

func (w wUserGroupDo) Clauses(conds ...clause.Expression) *wUserGroupDo {
	return w.withDO(w.DO.Clauses(conds...))
}

func (w wUserGroupDo) Returning(value interface{}, columns ...string) *wUserGroupDo {
	return w.withDO(w.DO.Returning(value, columns...))
}

func (w wUserGroupDo) Not(conds ...gen.Condition) *wUserGroupDo {
	return w.withDO(w.DO.Not(conds...))
}

func (w wUserGroupDo) Or(conds ...gen.Condition) *wUserGroupDo {
	return w.withDO(w.DO.Or(conds...))
}

func (w wUserGroupDo) Select(conds ...field.Expr) *wUserGroupDo {
	return w.withDO(w.DO.Select(conds...))
}

func (w wUserGroupDo) Where(conds ...gen.Condition) *wUserGroupDo {
	return w.withDO(w.DO.Where(conds...))
}

func (w wUserGroupDo) Order(conds ...field.Expr) *wUserGroupDo {
	return w.withDO(w.DO.Order(conds...))
}

func (w wUserGroupDo) Distinct(cols ...field.Expr) *wUserGroupDo {
	return w.withDO(w.DO.Distinct(cols...))
}

func (w wUserGroupDo) Omit(cols ...field.Expr) *wUserGroupDo {
	return w.withDO(w.DO.Omit(cols...))
}

func (w wUserGroupDo) Join(table schema.Tabler, on ...field.Expr) *wUserGroupDo {
	return w.withDO(w.DO.Join(table, on...))
}

func (w wUserGroupDo) LeftJoin(table schema.Tabler, on ...field.Expr) *wUserGroupDo {
	return w.withDO(w.DO.LeftJoin(table, on...))
}

func (w wUserGroupDo) RightJoin(table schema.Tabler, on ...field.Expr) *wUserGroupDo {
	return w.withDO(w.DO.RightJoin(table, on...))
}

func (w wUserGroupDo) Group(cols ...field.Expr) *wUserGroupDo {
	return w.withDO(w.DO.Group(cols...))
}

func (w wUserGroupDo) Having(conds ...gen.Condition) *wUserGroupDo {
	return w.withDO(w.DO.Having(conds...))
}

func (w wUserGroupDo) Limit(limit int) *wUserGroupDo {
	return w.withDO(w.DO.Limit(limit))
}

func (w wUserGroupDo) Offset(offset int) *wUserGroupDo {
	return w.withDO(w.DO.Offset(offset))
}

func (w wUserGroupDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *wUserGroupDo {
	return w.withDO(w.DO.Scopes(funcs...))
}

func (w wUserGroupDo) Unscoped() *wUserGroupDo {
	return w.withDO(w.DO.Unscoped())
}

func (w wUserGroupDo) Create(values ...*model.WUserGroup) error {
	if len(values) == 0 {
		return nil
	}
	return w.DO.Create(values)
}

func (w wUserGroupDo) CreateInBatches(values []*model.WUserGroup, batchSize int) error {
	return w.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (w wUserGroupDo) Save(values ...*model.WUserGroup) error {
	if len(values) == 0 {
		return nil
	}
	return w.DO.Save(values)
}

func (w wUserGroupDo) First() (*model.WUserGroup, error) {
	if result, err := w.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.WUserGroup), nil
	}
}

func (w wUserGroupDo) Take() (*model.WUserGroup, error) {
	if result, err := w.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.WUserGroup), nil
	}
}

func (w wUserGroupDo) Last() (*model.WUserGroup, error) {
	if result, err := w.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.WUserGroup), nil
	}
}

func (w wUserGroupDo) Find() ([]*model.WUserGroup, error) {
	result, err := w.DO.Find()
	return result.([]*model.WUserGroup), err
}

func (w wUserGroupDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.WUserGroup, err error) {
	buf := make([]*model.WUserGroup, 0, batchSize)
	err = w.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (w wUserGroupDo) FindInBatches(result *[]*model.WUserGroup, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return w.DO.FindInBatches(result, batchSize, fc)
}

func (w wUserGroupDo) Attrs(attrs ...field.AssignExpr) *wUserGroupDo {
	return w.withDO(w.DO.Attrs(attrs...))
}

func (w wUserGroupDo) Assign(attrs ...field.AssignExpr) *wUserGroupDo {
	return w.withDO(w.DO.Assign(attrs...))
}

func (w wUserGroupDo) Joins(fields ...field.RelationField) *wUserGroupDo {
	for _, _f := range fields {
		w = *w.withDO(w.DO.Joins(_f))
	}
	return &w
}

func (w wUserGroupDo) Preload(fields ...field.RelationField) *wUserGroupDo {
	for _, _f := range fields {
		w = *w.withDO(w.DO.Preload(_f))
	}
	return &w
}

func (w wUserGroupDo) FirstOrInit() (*model.WUserGroup, error) {
	if result, err := w.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.WUserGroup), nil
	}
}

func (w wUserGroupDo) FirstOrCreate() (*model.WUserGroup, error) {
	if result, err := w.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.WUserGroup), nil
	}
}

func (w wUserGroupDo) FindByPage(offset int, limit int) (result []*model.WUserGroup, count int64, err error) {
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

func (w wUserGroupDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = w.Count()
	if err != nil {
		return
	}

	err = w.Offset(offset).Limit(limit).Scan(result)
	return
}

func (w wUserGroupDo) Scan(result interface{}) (err error) {
	return w.DO.Scan(result)
}

func (w wUserGroupDo) Delete(models ...*model.WUserGroup) (result gen.ResultInfo, err error) {
	return w.DO.Delete(models)
}

func (w *wUserGroupDo) withDO(do gen.Dao) *wUserGroupDo {
	w.DO = *do.(*gen.DO)
	return w
}