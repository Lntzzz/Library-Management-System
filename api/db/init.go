package db

import (
	"fmt"
	"log"
	"time"

	"Library-Management-System/api/conf"
	//"Library-Management-System/api/log"
	"Library-Management-System/api/util/db"
)

func Init(mySQLConf *conf.MySQLConf, logger log.Logger) error {
	newDBConf := db.Conf{
		Dialect:         "mysql",
		User:            mySQLConf.User,
		Password:        mySQLConf.Passwd,
		Url:             fmt.Sprintf("%s:%d", mySQLConf.Ip, mySQLConf.Port),
		DbName:          mySQLConf.DB,
		CharSet:         "utf8mb4",
		MaxOpen:         mySQLConf.MaxConnection,
		MaxIdle:         10,
		ConnMaxLifetime: time.Duration(mySQLConf.MaxLifetime) * time.Second,
		Debug:           mySQLConf.DebugSQL,
		Plural:          "disable",
	}
	dbClient := db.CreateClient(&newDBConf)
	InitBookDao(dbClient)
	InitLoanRecordDao(dbClient)
	InitUserDao(dbClient)
	return nil
}
