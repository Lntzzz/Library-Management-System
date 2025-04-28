package db

import "time"

type Conf struct {
	Dialect         string        `yaml:"db.dialect"`         //方言: mysql、postgres
	User            string        `yaml:"db.user"`            //用户名
	Password        string        `yaml:"db.password"`        //密码
	Url             string        `yaml:"db.url"`             //地址
	DbName          string        `yaml:"db.dbName"`          //库名
	CharSet         string        `yaml:"db.charSet"`         //字符集
	MaxOpen         int           `yaml:"db.maxOpen"`         //连接池大小
	MaxIdle         int           `yaml:"db.maxIdle"`         //最大空闲连接
	ConnMaxLifetime time.Duration `yaml:"db.connMaxLifetime"` //连接有效时间
	SslMode         string        `yaml:"db.sslmode"`         //postgres有效
	Debug           bool          `yaml:"db.debug"`           //打印SQL
	Plural          string        `yaml:"db.plural"`          //表名复数化: enable、disable(默认)
}
