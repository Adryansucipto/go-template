package injection

import (
	pkgConfig "go-template/pkg/config"
	"go-template/pkg/resource/model"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(cfg pkgConfig.Config) (db model.Database, err error) {
	dsn := cfg.Database.PostgreSql.GenerateConnectionPostgresString()
	dbPostgreSql, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}

	sqlDB, err := dbPostgreSql.DB()
	if err != nil {
		return
	}

	maxIdle := cfg.Database.PostgreSql.MaxIdleConns
	if maxIdle <= 0 {
		maxIdle = 10
	}
	maxOpen := cfg.Database.PostgreSql.MaxOpenConns
	if maxOpen <= 0 {
		maxOpen = 100
	}
	maxLifetime := cfg.Database.PostgreSql.ConnMaxLifetime
	if maxLifetime <= 0 {
		maxLifetime = 60
	}

	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Minute)

	db = model.Database{
		PostgreSql: dbPostgreSql,
	}

	return
}
