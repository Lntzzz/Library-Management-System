package db

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
)

// Client DB操作接口
type Client interface {
	GetDB(ctx context.Context) DB
	GetTx(ctx context.Context, tx *Tx) DB
	BeginTx(ctx context.Context) (*Tx, error)
	Close() error

	ParseErr(err error) error
	ParseResult(result sql.Result, err error) (int64, error)

	// SetMaxOpenConns SetMaxIdleConns Stats 连接池管理
	SetMaxOpenConns(n int)
	SetMaxIdleConns(n int)
	Stats() sql.DBStats

	// SetLogger SetQueryHook 日志和监控
	//SetLogger(logger log.Logger)

	DML
}

type DB interface {
	bun.IDB
}

type Tx struct {
	*bun.Tx
}

// DML 数据操作接口
type DML interface {
	GetSelect(ctx context.Context, model interface{}) *bun.SelectQuery
	GetInsert(ctx context.Context, model interface{}, tx *Tx) *bun.InsertQuery
	GetUpdate(ctx context.Context, model interface{}, tx *Tx) *bun.UpdateQuery
	GetDelete(ctx context.Context, model interface{}, tx *Tx) *bun.DeleteQuery
	GetSoftDelete(ctx context.Context, model interface{}, tx *Tx) *bun.UpdateQuery

	Insert(ctx context.Context, data interface{}, tx *Tx) (int64, error)
	BatchInsert(ctx context.Context, data []interface{}, tx *Tx, batchSize int) (int64, error)
	Delete(ctx context.Context, model interface{}, tx *Tx, id interface{}) (int64, error)
	SoftDelete(ctx context.Context, model interface{}, tx *Tx, id interface{}) (int64, error)
	QueryAll(ctx context.Context, result interface{}) error
	QueryPage(ctx context.Context, result interface{}, filter DataFilter) (int, error)
}

// DataFilter 数据过滤器接口
type DataFilter interface {
	GetPage() (offset int, limit int, disable bool)
	Filter(db DB) *bun.SelectQuery
	GetOrders() []string
}

// TxFunc 事务函数类型
type TxFunc func(ctx context.Context, tx *Tx) error

// TransactionManager 事务管理器
type TransactionManager interface {
	WithinTransaction(ctx context.Context, txFunc TxFunc) error
	WithSavepoint(ctx context.Context, tx *Tx, spFunc TxFunc) error
}
