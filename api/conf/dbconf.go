package conf

import (
	"database/sql"
	"fmt"
	"time"
)

type MySQLConf struct {
	Ip            string `mapstructure:"ip"`
	Port          int    `mapstructure:"port"`
	User          string `mapstructure:"user"`
	Passwd        string `mapstructure:"passwd"`
	DB            string `mapstructure:"db"`
	MaxConnection int    `mapstructure:"max_connection"`
	MaxLifetime   int    `mapstructure:"max_lifetime"`
	Timeout       int    `mapstructure:"timeout"`
	DebugSQL      bool   `mapstructure:"debug_sql"`
}

func (m *MySQLConf) validate() error {
	if m == nil {
		return fmt.Errorf("mysql config should not be empty")
	}
	if m.Ip == "" {
		return fmt.Errorf("mysql ip should not be empty")
	}
	if m.Port <= 0 || m.Port > 65535 {
		return fmt.Errorf("mysql port should be in range 1-65535")
	}
	if m.User == "" {
		return fmt.Errorf("mysql user should not be empty")
	}
	if m.Passwd == "" {
		return fmt.Errorf("mysql passwd should not be empty")
	}
	if m.DB == "" {
		return fmt.Errorf("mysql db should not be empty")
	}
	return nil
}

type RedisConf struct {
	Addr         string `mapstructure:"addr"`
	Password     string `mapstructure:"password"`
	Database     int    `mapstructure:"database"`
	PoolSize     int    `mapstructure:"pool_size"`
	ConnTimeout  int    `mapstructure:"conn_timeout"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
	IdleTimeout  int    `mapstructure:"idle_timeout"`
}

const (
	MySQLConfKey = "mysql"
	RedisConfKey = "redis"
)

func LoadMySQLConf(confPath string) (*MySQLConf, error) {
	v := newViperWithConfPath(confPath)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read mysql config file: %w", err)
	}

	config := &MySQLConf{
		Ip:            "127.0.0.1",
		Port:          3306,
		DB:            "coc_lb_openapi",
		Timeout:       3000,
		MaxConnection: 300,
		MaxLifetime:   1000,
	}

	if err := v.UnmarshalKey(MySQLConfKey, config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal mysql config: %w", err)
	}
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("failed to validate mysql config: %w", err)
	}
	return config, nil
}

func NewMySQLInstance(conf *MySQLConf) (*sql.DB, error) {
	strConn := "%s:%s@tcp(%s:%d)/%s?autocommit=true&parseTime=true&timeout=%dms&loc=Asia%%2FShanghai&tx_isolation='READ-COMMITTED'"
	url := fmt.Sprintf(strConn, conf.User, conf.Passwd,
		conf.Ip, conf.Port, conf.DB, conf.Timeout)
	db, err := sql.Open("mysql", url)
	if err != nil {
		fmt.Printf("[Mysql/cfg/NewMysqlInstance] [sql.Open: %s, url: %s]\n", err.Error(), url)
		return nil, err
	}
	fmt.Printf("open mysql success\n")
	db.SetMaxOpenConns(conf.MaxConnection)
	db.SetMaxIdleConns(conf.MaxConnection)
	db.SetConnMaxLifetime(time.Second * time.Duration(conf.MaxLifetime))

	err = db.Ping()
	if err != nil {
		fmt.Printf("[Mysql/cfg/NewMySQLInstance] [db.Ping: %s]\n", err.Error())
		return nil, err
	}
	return db, nil
}

func LoadRedisConf(confPath string) (*RedisConf, error) {
	v := newViperWithConfPath(confPath)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read redis config file: %w", err)
	}

	var config RedisConf

	if err := v.UnmarshalKey(RedisConfKey, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal redis config: %w", err)
	}

	return &config, nil
}
