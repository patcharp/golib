package database

import (
	"fmt"
	"github.com/carlescere/scheduler"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Host       string
	Port       string
	Username   string
	Password   string
	Name       string
	DebugMode  bool
	Filename   string
	SSLEnabled bool
	Charset    string
}

type Database struct {
	config Config
	driver string
	dsn    string
	ctx    *gorm.DB
	job    *scheduler.Job
}

const (
	DriverMSSQL    = "mssql"
	DriverMySQL    = "mysql"
	DriverSQLLite  = "sqlite3"
	DriverPostgres = "postgres"

	DefaultMySQLCharset = "utf8mb4"
)

var dbContext []*gorm.DB

func GetConnectionContext() []*gorm.DB {
	return dbContext
}

func NewWithConfig(cfg Config, driver string) Database {
	return Database{
		config: cfg,
		driver: driver,
	}
}

func (db *Database) Connect() error {
	var err error
	switch db.driver {
	case DriverMSSQL:
		db.dsn = fmt.Sprintf(
			"sqlserver://%s:%s@%s:%s?database=%s",
			db.config.Username,
			db.config.Password,
			db.config.Host,
			db.config.Port,
			db.config.Name,
		)
	case DriverPostgres:
		sslMode := "disable"
		if db.config.SSLEnabled {
			sslMode = "enable"
		}
		db.dsn = fmt.Sprintf(
			"postgresql://%s@%s:%s/%s?sslmode=%s",
			db.config.Username,
			db.config.Host,
			db.config.Port,
			db.config.Name,
			sslMode,
		)
	case DriverSQLLite:
		db.dsn = db.config.Filename
	default:
		db.driver = DriverMySQL
		if db.config.Charset == "" {
			db.config.Charset = DefaultMySQLCharset
		}
		db.dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=Local",
			db.config.Username,
			db.config.Password,
			db.config.Host,
			db.config.Port,
			db.config.Name,
			db.config.Charset,
		)
	}
	db.ctx, err = gorm.Open(db.driver, db.dsn)
	if err != nil {
		return err
	}
	db.ctx.LogMode(db.config.DebugMode)
	if err := db.startKeepAlive(); err != nil {
		return err
	}
	dbContext = append(dbContext, db.ctx)
	return nil
}

func (db *Database) Reconnect() error {
	var err error
	db.ctx, err = gorm.Open(db.driver, db.dsn)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) Ctx() *gorm.DB {
	return db.ctx
}

func (db *Database) MigrateDatabase(tables []interface{}) error {
	tx := db.ctx.Begin()
	for _, t := range tables {
		if err := tx.AutoMigrate(t).Error; err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (db *Database) SetDebugMode(mode bool) {
	db.ctx.LogMode(mode)
}

func (db *Database) startKeepAlive() error {
	var err error
	db.job, err = scheduler.Every(15).Seconds().Run(func() {
		if err := db.ctx.DB().Ping(); err != nil {
			logrus.Errorln("Database keep alive error ->", err)
			if err := db.Reconnect(); err != nil {
				logrus.Errorln("Trying to reconnect to database error ->", err)
			} else {
				logrus.Infoln("Database reconnect success.")
			}
		}
	})
	return err
}

func (db *Database) stopKeepAlive() error {
	if db.job != nil {
		db.job.Quit <- true
	}
	return nil
}

func (db *Database) Close() error {
	_ = db.stopKeepAlive()
	if err := db.ctx.Close(); err != nil {
		return err
	}
	return nil
}
