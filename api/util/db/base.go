package db

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"github.com/uptrace/bun"
)

// BaseDao other dao can extend this to use the common method
type BaseDao[T any] struct {
	Client
}

func NewBaseDao[T Identifiable](dbc Client) *BaseDao[T] {
	return &BaseDao[T]{Client: dbc}
}
func (d *BaseDao[T]) TxBatchInsert(ctx context.Context, tx *Tx, params []*T) (count int64, err error) {
	if len(params) == 0 {
		return 0, nil
	}
	pageNumber := 1
	batch := 100
	count = int64(0)
	for {
		canSplit, start, end := SplitIndex(pageNumber, batch, len(params))
		if !canSplit {
			break
		}
		batchList := params[start:end]
		cnt, err := d.ParseResult(d.GetTx(ctx, tx).NewInsert().Model(&batchList).Exec(ctx))
		if err != nil {
			return count, err
		}
		count += cnt
		pageNumber++
	}
	return count, nil
}

func (d *BaseDao[T]) Create(ctx context.Context, data *T) (sql.Result, error) {
	return d.GetDB(ctx).NewInsert().Model(data).Exec(ctx)
}

func (d *BaseDao[T]) TxCreate(ctx context.Context, tx *Tx, data *T) (sql.Result, error) {
	return d.GetTx(ctx, tx).NewInsert().Model(data).Exec(ctx)
}

func (d *BaseDao[T]) TxDelete(ctx context.Context, tx *Tx, id string) error {
	if _, err := d.GetTx(ctx, tx).NewDelete().Model((*T)(nil)).
		Where("`id`=?", id).Exec(ctx); err != nil {
		return d.ParseErr(err)
	}
	return nil
}

func (d *BaseDao[T]) TxBatchDelete(ctx context.Context, tx *Tx, ids []string) error {
	if _, err := d.GetTx(ctx, tx).NewDelete().Model((*T)(nil)).
		Where("`id` IN (?)", bun.In(ids)).Exec(ctx); err != nil {
		return d.ParseErr(err)
	}
	return nil
}

// List list data by filters, od: 0 DESC 1 ASC
func (d *BaseDao[T]) List(ctx context.Context, filters map[string]interface{}, limit int, offset int, order string, od int) ([]*T, error) {
	list := make([]*T, 0)
	query := d.GetDB(ctx).NewSelect().Model(&list)
	for k, val := range filters {
		rt := reflect.TypeOf(val)
		switch rt.Kind() {
		case reflect.Slice, reflect.Array:
			query.Where(k+" in (?) ", bun.In(val))
		default:
			query.Where(k+" = ? ", val)
		}
	}
	if offset > 0 {
		query.Offset(offset)
	}
	if limit > 0 {
		query.Limit(limit)
	}
	if len(order) > 0 { // 0 DESC 1 ASC, default DESC
		if od == 0 {
			query.Order(fmt.Sprintf("%s DESC", order))
		} else {
			query.Order(fmt.Sprintf("%s", order))
		}
	}
	return list, query.Scan(ctx)
}

type FilterOp string

const (
	FilterOpIn   FilterOp = "in"
	FilterOpEq   FilterOp = "eq"
	FilterOpNeq  FilterOp = "neq" // not equal
	FilterOpLike FilterOp = "like"
)

type Filter struct {
	Col string
	Op  FilterOp
	Val interface{}
}

func (d *BaseDao[T]) Filtrate(ctx context.Context, filters []Filter) ([]*T, error) {
	list := make([]*T, 0)
	query := d.GetDB(ctx).NewSelect().Model(&list)
	if err := d.appendFilter(query, filters); err != nil {
		return nil, err
	}
	return list, query.Scan(ctx)
}

func (d *BaseDao[T]) CountByFiltrate(ctx context.Context, filters []Filter) (int32, error) {
	query := d.GetDB(ctx).NewSelect().Model((*T)(nil))
	if err := d.appendFilter(query, filters); err != nil {
		return 0, err
	}
	count, err := query.Count(ctx)
	return int32(count), err
}

func (*BaseDao[T]) appendFilter(query *bun.SelectQuery, filters []Filter) error {
	for _, filter := range filters {
		var (
			col = filter.Col
			op  = filter.Op
			val = filter.Val
		)
		if col == "" || val == nil {
			return fmt.Errorf("filter col, val should not be empty, filter: %+v", filter)
		}
		if op == "" {
			op = FilterOpEq // default eq
		}
		switch op {
		case FilterOpEq:
			valType := reflect.TypeOf(val)
			switch valType.Kind() {
			case reflect.Slice, reflect.Array:
				query.Where(col+" IN (?) ", bun.In(val))
			default:
				query.Where(col+" = ? ", val)
			}
		case FilterOpNeq:
			rt := reflect.TypeOf(val)
			switch rt.Kind() {
			case reflect.Slice, reflect.Array:
				query.Where(col+" NOT IN (?) ", bun.In(val))
			default:
				query.Where(col+" != ? ", val)
			}
		case FilterOpLike:
			valStr, ok := val.(string)
			if !ok {
				return fmt.Errorf("filter op %s value should be string, filter.col: %s", op, col)
			}
			query.Where(col+" like ? ", "%"+valStr+"%")
		case FilterOpIn:
			valType := reflect.TypeOf(val)
			switch valType.Kind() {
			case reflect.Slice, reflect.Array:
				query.Where(col+" IN (?) ", bun.In(val))
			default:
				query.Where(col+" IN (?) ", bun.In([]interface{}{val}))
			}
		default:
			return fmt.Errorf("unsupported filter op: %s", op)
		}
	}
	return nil
}

func (d *BaseDao[T]) ListByFilter(ctx context.Context, filters map[string]interface{}) ([]*T, error) {
	list := make([]*T, 0)
	query := d.GetDB(ctx).NewSelect().Model(&list)
	for k, val := range filters {
		rt := reflect.TypeOf(val)
		switch rt.Kind() {
		case reflect.Slice, reflect.Array:
			query.Where(k+" in (?) ", bun.In(val))
		default:
			query.Where(k+" = ? ", val)
		}
	}
	return list, query.Scan(ctx)
}

func (d *BaseDao[T]) ListByFilterWithFuzzy(ctx context.Context, eqFilters map[string]interface{}, fuzzyFilters map[string]string) ([]*T, error) {
	list := make([]*T, 0)
	query := d.GetDB(ctx).NewSelect().Model(&list)
	for k, val := range eqFilters {
		rt := reflect.TypeOf(val)
		switch rt.Kind() {
		case reflect.Slice, reflect.Array:
			query.Where(k+" in (?) ", bun.In(val))
		default:
			query.Where(k+" = ? ", val)
		}
	}
	for field, fuzzy := range fuzzyFilters {
		query.Where(field+" like ? ", "%"+fuzzy+"%")
	}
	return list, query.Scan(ctx)
}

func (d *BaseDao[T]) CountByFilter(ctx context.Context, filters map[string]interface{}) (int32, error) {
	query := d.GetDB(ctx).NewSelect().Model((*T)(nil))
	for k, val := range filters {
		rt := reflect.TypeOf(val)
		switch rt.Kind() {
		case reflect.Slice, reflect.Array:
			query.Where(k+" in (?) ", bun.In(val))
		default:
			query.Where(k+" = ? ", val)
		}
	}
	count, err := query.Count(ctx)
	return int32(count), err
}

func (d *BaseDao[T]) TxListByFilter(ctx context.Context, tx *Tx, filters map[string]interface{}) ([]*T, error) {
	list := make([]*T, 0)
	query := d.GetTx(ctx, tx).NewSelect().Model(&list)
	for k, val := range filters {
		rt := reflect.TypeOf(val)
		switch rt.Kind() {
		case reflect.Slice, reflect.Array:
			query.Where(k+" in (?) ", bun.In(val))
		default:
			query.Where(k+" = ? ", val)
		}
	}
	return list, query.Scan(ctx)
}

func (d *BaseDao[T]) TxUpdateFieldsById(ctx context.Context, tx *Tx, id string, updateFields map[string]interface{}) error {
	return d.TxUpdateFields(ctx, tx, updateFields, map[string]interface{}{"id": id})
}

func (d *BaseDao[T]) TxUpdateFields(ctx context.Context, tx *Tx, updateFields, whereFields map[string]interface{}) error {
	update := d.GetUpdate(ctx, (*T)(nil), tx)
	for updateKey, updateField := range updateFields {
		update.Set(updateKey+" = ?", updateField)
	}
	for whereKey, where := range whereFields {
		rt := reflect.TypeOf(where)
		switch rt.Kind() {
		case reflect.Slice, reflect.Array:
			update.Where(whereKey+" in (?) ", bun.In(where))
		default:
			update.Where(whereKey+" = ?", where)
		}
	}
	rows, err := update.Exec(ctx)
	_, err = d.ParseResult(rows, err)
	return err
}

func (d *BaseDao[T]) UpdateFieldsById(ctx context.Context, id string, updateFields map[string]interface{}) error {
	return d.UpdateFields(ctx, updateFields, map[string]interface{}{"id": id})
}

func (d *BaseDao[T]) UpdateFields(ctx context.Context, updateFields, whereFields map[string]interface{}) error {
	update := d.GetDB(ctx).NewUpdate().Model((*T)(nil))
	for updateKey, updateField := range updateFields {
		update.Set(updateKey+" = ?", updateField)
	}
	for whereKey, where := range whereFields {
		rt := reflect.TypeOf(where)
		switch rt.Kind() {
		case reflect.Slice, reflect.Array:
			update.Where(whereKey+" in (?) ", bun.In(where))
		default:
			update.Where(whereKey+" = ?", where)
		}
	}
	rows, err := update.Exec(ctx)
	_, err = d.ParseResult(rows, err)
	return err
}

func (d *BaseDao[T]) FindById(ctx context.Context, id string) (*T, error) {
	var ret = new(T)
	return ret, d.GetDB(ctx).NewSelect().
		Model(ret).
		Where("id = ?", id).
		Scan(ctx)
}

func (d *BaseDao[T]) TxFindById(ctx context.Context, tx *Tx, id string) (*T, error) {
	var ret = new(T)
	return ret, d.GetTx(ctx, tx).NewSelect().
		Model(ret).
		Where("id = ?", id).
		Scan(ctx)
}

func (d *BaseDao[T]) TxFindByIdForUpdate(ctx context.Context, tx *Tx, id string) (*T, error) {
	var ret = new(T)
	return ret, d.GetTx(ctx, tx).NewSelect().
		Model(ret).
		For("update").
		Where("id = ?", id).
		Scan(ctx)
}

func (d *BaseDao[T]) FindByIdAndUser(ctx context.Context, id, userId, region string) (*T, error) {
	var ret = new(T)
	return ret, d.GetDB(ctx).NewSelect().
		Model(ret).
		Where("id = ?", id).
		Where("user_id = ?", userId).
		Where("region = ?", region).
		Scan(ctx)
}

func (d *BaseDao[T]) FindByUserAndId(ctx context.Context, userId string, region string, Id string) (*T, error) {
	var ret = new(T)
	return ret, d.GetDB(ctx).NewSelect().
		Model(ret).
		Where("id = ?", Id).
		Where("user_id = ?", userId).
		Where("region = ?", region).
		Scan(ctx)
}

// TxGetForUpdate get data for update, must has id field, and use in transaction
func (d *BaseDao[T]) TxGetForUpdate(ctx context.Context, tx *Tx, id interface{}) (*T, error) {
	var ret = new(T)
	err := d.GetTx(ctx, tx).NewSelect().Model(ret).For("update").Where("id=?", id).Scan(ctx)
	return ret, err
}

func (d *BaseDao[T]) ListByCocIds(ctx context.Context, cocIDs []string) ([]*T, error) {
	list := make([]*T, 0)
	return list, d.GetDB(ctx).NewSelect().Model(&list).Where("coc_id in (?)", bun.In(cocIDs)).Scan(ctx, &list)
}

func (d *BaseDao[T]) TxListByCocIds(ctx context.Context, tx *Tx, cocIDs []string) ([]*T, error) {
	list := make([]*T, 0)
	return list, d.GetTx(ctx, tx).NewSelect().Model(&list).Where("coc_id in (?)", bun.In(cocIDs)).Scan(ctx, &list)
}

var ErrNotFound = sql.ErrNoRows
