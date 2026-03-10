package provider

import (
	"database/sql"
	"fmt"
	"sync"

	ezgorm "github.com/itsLeonB/ezutil/v2/gorm"
	"github.com/itsLeonB/ungerr"
	"github.com/reflect-homini/stora/internal/core/config"
	"github.com/reflect-homini/stora/internal/core/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DataSources struct {
	Gorm *gorm.DB
	SQL  *sql.DB
}

func (ds *DataSources) Shutdown() error {
	if err := ds.SQL.Close(); err != nil {
		return ungerr.Wrap(err, "error closing SQL db")
	}
	return nil
}

func ProvideDataSource() (*DataSources, error) {
	gormDB, sqlDB, err := provideAndConfigureSQL(config.Global.DB)
	if err != nil {
		return nil, err
	}

	return &DataSources{
		Gorm: gormDB,
		SQL:  sqlDB,
	}, nil
}

var (
	sqlInstance *sqlConnection
	sqlOnce     sync.Once
)

type sqlConnection struct {
	gormDB *gorm.DB
	sqlDB  *sql.DB
}

func provideAndConfigureSQL(cfg config.DB) (*gorm.DB, *sql.DB, error) {
	var err error
	sqlOnce.Do(func() {
		gormDB, e := gorm.Open(postgres.Open(dsn(cfg)), &gorm.Config{
			Logger: ezgorm.NewGormLogger(logger.Global),
		})
		if e != nil {
			err = ungerr.Wrap(e, "error opening gorm connection")
			return
		}

		sqlDB, e := gormDB.DB()
		if e != nil {
			err = ungerr.Wrap(e, "error obtaining sql.DB instance")
			return
		}

		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

		if e = sqlDB.Ping(); e != nil {
			err = ungerr.Wrap(e, "error pinging SQL DB")
			return
		}

		sqlInstance = &sqlConnection{gormDB: gormDB, sqlDB: sqlDB}
	})

	if err != nil {
		return nil, nil, err
	}

	return sqlInstance.gormDB, sqlInstance.sqlDB, nil
}

func dsn(cfg config.DB) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
	)
}
