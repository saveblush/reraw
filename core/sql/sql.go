package sql

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/saveblush/reraw/core/utils"
)

var (
	// Database global variable
	Database = &gorm.DB{}
)

var (
	defaultMaxIdleConns = 10
	defaultMaxOpenConns = 30
	defaultMaxLifetime  = time.Minute
	defaultCharset      = "utf8mb4"
)

// gorm config
var defaultConfig = &gorm.Config{
	PrepareStmt:            false,
	SkipDefaultTransaction: true,
	DisableAutomaticPing:   false,
	Logger:                 logger.Default.LogMode(logger.Error),
}

// Session session
type Session struct {
	Database *gorm.DB
}

// Configuration config mysql
type Configuration struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
	Charset      string
	Timezone     string
	MaxIdleConns int
	MaxOpenConns int
	MaxLifetime  time.Duration
}

// InitConnectionMysql open initialize a new db connection.
func InitConnection(cf *Configuration) (*Session, error) {
	var db *gorm.DB
	var err error

	// check params config
	if cf.Charset == "" {
		cf.Charset = defaultCharset
	}
	if cf.Timezone == "" {
		cf.Timezone = utils.TimeZone()
	}
	if cf.MaxIdleConns == 0 {
		cf.MaxIdleConns = defaultMaxIdleConns
	}
	if cf.MaxOpenConns == 0 {
		cf.MaxOpenConns = defaultMaxOpenConns
	}
	if cf.MaxLifetime == 0 {
		cf.MaxLifetime = defaultMaxLifetime
	}

	// create database
	err = createDatabase(cf)
	if err != nil {
		return nil, err
	}

	// connect db postgres
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s TimeZone=%s sslmode=disable",
		cf.Username,
		cf.Password,
		cf.Host,
		cf.Port,
		cf.DatabaseName,
		cf.Timezone,
	)
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                 dsn,
		WithoutQuotingCheck: true,
	}), defaultConfig)
	if err != nil {
		return nil, err
	}

	// connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cf.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cf.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cf.MaxLifetime)

	fmt.Printf("%s \n", "--------------------------------------------------")
	fmt.Printf("DB Stats [host: %s dbname: %s]\n", cf.Host, cf.DatabaseName)
	fmt.Printf("MaxOpenConnections: %v\n", sqlDB.Stats().MaxOpenConnections)
	fmt.Printf("%s \n\n", "--------------------------------------------------")

	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	return &Session{Database: db}, nil
}

// CloseConnection close connection db
func CloseConnection(db *gorm.DB) error {
	c, err := db.DB()
	if err != nil {
		return err
	}

	err = c.Close()
	if err != nil {
		return err
	}

	return nil
}

// DebugDatabase set debug sql
func DebugDatabase() {
	Database = Database.Debug()
}
