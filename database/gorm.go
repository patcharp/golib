package database

import (
	"database/sql"
	"fmt"
	"github.com/carlescere/scheduler"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
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

type MySQLConfig struct {
	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
	DebugMode    bool
}

type Database struct {
	config     Config
	driver     string
	dsn        string
	ctx        *gorm.DB
	sql        *sql.DB
	gormConfig gorm.Config
	job        *scheduler.Job
}

const (
	DriverMSSQL    = "mssql"
	DriverMySQL    = "mysql"
	DriverSQLLite  = "sqlite3"
	DriverPostgres = "postgres"

	DefaultMySQLCharset = "utf8mb4"
)

func NewWithConfig(cfg Config, driver string) Database {
	return Database{
		config: cfg,
		driver: driver,
	}
}

func NewMySqlWithConfig(my MySQLConfig) Database {
	return NewWithConfig(Config{
		Host:      my.Host,
		Port:      my.Port,
		Username:  my.Username,
		Password:  my.Password,
		Name:      my.DatabaseName,
		DebugMode: my.DebugMode,
		Charset:   DefaultMySQLCharset,
	},
		DriverMySQL,
	)
}

func (db *Database) Connect() error {
	return db.ConnectWithGormConfig(gorm.Config{
		Logger: NewGormLog(),
	})
}

func (db *Database) ConnectWithGormConfig(gormCfg gorm.Config) error {
	_ = db.stopKeepAlive()
	db.gormConfig = gormCfg
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
		db.ctx, err = gorm.Open(sqlserver.Open(db.dsn), &gormCfg)
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
		db.ctx, err = gorm.Open(postgres.Open(db.dsn), &gormCfg)
	case DriverSQLLite:
		db.dsn = db.config.Filename
		db.ctx, err = gorm.Open(sqlite.Open(db.dsn), &gormCfg)
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
		db.ctx, err = gorm.Open(mysql.Open(db.dsn), &gormCfg)
	}
	if err != nil {
		return err
	}
	if db.config.DebugMode {
		db.ctx = db.ctx.Debug()
	}

	db.sql, err = db.ctx.DB()
	if err != nil {
		return err
	}
	if err = db.startKeepAlive(); err != nil {
		return err
	}
	return nil
}

func (db *Database) Reconnect() error {
	return db.ConnectWithGormConfig(db.gormConfig)
}

func (db *Database) Ctx() *gorm.DB {
	return db.ctx
}

func (db *Database) SqlDB() *sql.DB {
	return db.sql
}

func (db *Database) MigrateDatabase(tables []interface{}) error {
	tx := db.ctx.Begin()
	for _, t := range tables {
		if err := tx.AutoMigrate(t); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (db *Database) startKeepAlive() error {
	var err error
	db.job, err = scheduler.Every(15).Seconds().Run(func() {
		if err := db.sql.Ping(); err != nil {
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
	if err := db.sql.Close(); err != nil {
		return err
	}
	return nil
}
