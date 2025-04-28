package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/schema"
)

type client struct {
	db   *bun.DB
	conf *Conf
	//logger log.Logger
}

// CreateClient 创建数据库连接
func CreateClient(conf *Conf) Client {
	if conf.MaxOpen == 0 {
		conf.MaxOpen = 10
	}
	if conf.MaxIdle == 0 {
		conf.MaxIdle = 3
	}
	if conf.ConnMaxLifetime == 0 {
		conf.ConnMaxLifetime = 8 * time.Hour
	}
	if conf.CharSet == "" {
		conf.CharSet = "utf8mb4"
	}
	if conf.SslMode == "" {
		conf.SslMode = "disable"
	}

	marshal, err := json.Marshal(conf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("<<<<<<<<< db config >>>>>>>>> \n%s\n", string(marshal))

	// 连接
	var dbc *client
	switch conf.Dialect {
	case "mysql":
		dbc = createMysqlClient(conf)
	case "postgres":
		dbc = createPostgresClient(conf)
	default:
		panic(fmt.Sprintf("unkown database dialect '%s'", conf.Dialect))
	}

	//dbc.SetLogger(logger)
	//dbc.db.AddQueryHook(dbc) // 设置钩子

	// 关闭表名复数形式
	if conf.Plural == "disable" {
		schema.SetTableNameInflector(func(tableName string) string {
			return tableName
		})
	}

	return dbc
}

func (c *client) SetMaxOpenConns(n int) {
	c.db.SetMaxOpenConns(n)
}

func (c *client) SetMaxIdleConns(n int) {
	c.db.SetMaxIdleConns(n)
}

func (c *client) Stats() sql.DBStats {
	return c.db.Stats()
}

//func (c *client) SetLogger(logger log.Logger) {
//	c.logger = logger
//}

// BeforeQuery SQL执行前拦截监控
func (c *client) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	_ = event
	return ctx
}

func (c *client) Close() (err error) {
	return
}

// GetDB 获得数据库连接，无事务
func (c *client) GetDB(_ context.Context) DB {
	return c.db
}

// GetTx 获得数据库连接，如果tx=nil则返回无事务连接，否则返回tx
func (c *client) GetTx(_ context.Context, tx *Tx) DB {
	if tx != nil {
		return tx
	}
	return c.db
}

// BeginTx 开启新事务
func (c *client) BeginTx(ctx context.Context) (*Tx, error) {
	return c.beginTx(ctx, nil)
}

// beginTx 开启新事务
func (c *client) beginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if tx, err := c.db.BeginTx(ctx, opts); err != nil {
		return nil, err
	} else {
		return &Tx{Tx: &tx}, nil
	}
}

// GetSelect 获得通用处理器：查询
func (c *client) GetSelect(ctx context.Context, model interface{}) *bun.SelectQuery {
	return c.GetDB(ctx).NewSelect().Model(model)
}

// GetInsert 获得通用处理器：写入
func (c *client) GetInsert(ctx context.Context, model interface{}, tx *Tx) *bun.InsertQuery {
	return c.GetTx(ctx, tx).NewInsert().Model(model)
}

// GetUpdate 获得通用处理器：更新
func (c *client) GetUpdate(ctx context.Context, model interface{}, tx *Tx) *bun.UpdateQuery {
	return c.GetTx(ctx, tx).NewUpdate().Model(model)
}

// GetDelete 获得通用处理器：删除
func (c *client) GetDelete(ctx context.Context, model interface{}, tx *Tx) *bun.DeleteQuery {
	return c.GetTx(ctx, tx).NewDelete().Model(model)
}

// GetSoftDelete 获得通用处理器：逻辑删除
func (c *client) GetSoftDelete(ctx context.Context, model interface{}, tx *Tx) *bun.UpdateQuery {
	return c.GetTx(ctx, tx).NewUpdate().Model(model).Set("update_time=current_timestamp").Set("status=-1")
}

// Insert 通用处理：插入数据(单条及批量处理，批量太大时不要不要使用)
func (c *client) Insert(ctx context.Context, data interface{}, tx *Tx) (int64, error) {
	return c.ParseResult(c.GetTx(ctx, tx).NewInsert().Model(data).Exec(ctx))
}

func (c *client) BatchInsert(ctx context.Context, data []interface{}, tx *Tx, batchSize int) (int64, error) {
	if len(data) == 0 {
		return 0, nil
	}
	pageNumber := 1
	if batchSize <= 10 {
		batchSize = 100
	} else {
		batchSize = 100
	}
	count := int64(0)
	for {
		canSplit, start, end := SplitIndex(pageNumber, batchSize, len(data))
		if !canSplit {
			break
		}
		batchList := data[start:end]
		cnt, err := c.ParseResult(c.GetTx(ctx, tx).NewInsert().Model(&batchList).Exec(ctx))
		if err != nil {
			return count, err
		}
		count += cnt
		pageNumber++
	}
	return count, nil
}

// SoftDelete 通用处理：逻辑删除(id可以是单个也可以是数组)
func (c *client) SoftDelete(ctx context.Context, model interface{}, tx *Tx, id interface{}) (int64, error) {
	handler := c.GetSoftDelete(ctx, model, tx)
	if reflect.TypeOf(id).Kind() == reflect.Slice {
		handler.Where("id in (?)", bun.In(id))
	} else {
		handler.Where("id=?", id)
	}
	return c.ParseResult(handler.Exec(ctx))
}

// Delete 通用处理：物理删除(id可以是单个也可以是数组)
func (c *client) Delete(ctx context.Context, model interface{}, tx *Tx, id interface{}) (int64, error) {
	handler := c.GetDelete(ctx, model, tx)
	if reflect.TypeOf(id).Kind() == reflect.Slice {
		handler.Where("id in (?)", bun.In(id))
	} else {
		handler.Where("id=?", id)
	}
	return c.ParseResult(handler.Exec(ctx))
}

// QueryAll 通用处理：查询全量
func (c *client) QueryAll(ctx context.Context, result interface{}) error {
	return c.GetSelect(ctx, result).Order("id desc").Scan(ctx, result)
}

// QueryPage 通用处理：分页查询
func (c *client) QueryPage(ctx context.Context, result interface{}, filter DataFilter) (int, error) {

	offset, limit, disable := filter.GetPage()

	if disable {
		count := 0
		err := filter.Filter(c.GetDB(ctx)).Order(filter.GetOrders()...).Scan(ctx, result)
		v := reflect.ValueOf(result)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.Slice {
			count = v.Len()
		}
		return count, err
	}

	// 分页
	return filter.Filter(c.GetDB(ctx)).Order(filter.GetOrders()...).Offset(offset).Limit(limit).ScanAndCount(ctx, result)
}

// ParseErr 通用处理：单条查询后消化掉空数据错误
func (c *client) ParseErr(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}

// ParseResult 通用处理：获取执行结果影响的记录数量
func (c *client) ParseResult(result sql.Result, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

///////////////////////////// private ///////////////////////////

// 创建mysql数据库连接
func createMysqlClient(cfg *Conf) *client {

	// 连接
	passwd := cfg.Password
	connStr := cfg.User + ":" + passwd + "@tcp(" + cfg.Url + ")/" + cfg.DbName + "?charset=" + cfg.CharSet
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(cfg.MaxOpen)
	db.SetMaxIdleConns(cfg.MaxIdle)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	bunDB := bun.NewDB(db, mysqldialect.New())
	if err := bunDB.Ping(); err != nil {
		panic(err)
	}

	return &client{db: bunDB, conf: cfg}
}

// 创建postgres数据库连接
func createPostgresClient(cfg *Conf) *client {
	// 连接
	passwd := cfg.Password
	dsn := "postgres://" + cfg.User + ":" + passwd + "@" + cfg.Url + "/" + cfg.DbName + "?sslmode=" + cfg.SslMode
	db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db.SetMaxOpenConns(cfg.MaxOpen)
	db.SetMaxIdleConns(cfg.MaxIdle)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// 包装
	bunDB := bun.NewDB(db, pgdialect.New())
	if err := bunDB.Ping(); err != nil {
		panic(err)
	}

	return &client{db: bunDB, conf: cfg}
}
