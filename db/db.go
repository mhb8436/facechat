package db

import (
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbMap = map[string]*gorm.DB{}
var syncLock sync.Mutex

func init() {
	// InitDB("facechat")
	InitPostgres("lilac")
}

func InitDB(dbName string) {
	var e error
	// if prod env , you should change mysql driver for yourself !!!
	realPath, _ := filepath.Abs("./")
	configFilePath := realPath + "/db/facechat.db"
	syncLock.Lock()
	dbMap[dbName], e = gorm.Open(sqlite.Open(configFilePath), &gorm.Config{})
	syncLock.Unlock()
	if e != nil {
		logrus.Error("connect db fail:%s", e.Error())
	}
}

func InitPostgres(dbName string) {
	// dsn := "host=localhost user=mhb8436 password=dbwj1234 dbname=lilac port=5432"
	dsn := "host=localhost user=**** password=**** dbname=**** port=5432"
	syncLock.Lock()
	db, e := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	dbMap[dbName] = db
	syncLock.Unlock()
	if e != nil {
		logrus.Error("postgres connect db fail:%s", e.Error())
	}
}

func GetDb(dbName string) (db *gorm.DB) {
	if db, ok := dbMap[dbName]; ok {
		return db
	} else {
		return nil
	}
}

type DbFaceChat struct {
}

func (*DbFaceChat) GetDbName() string {
	return "facechat"
}
