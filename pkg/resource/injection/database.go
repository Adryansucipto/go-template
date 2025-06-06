package injection

import (
	pkgConfig "go-template/pkg/config"
	"go-template/pkg/resource/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(cfg pkgConfig.Config) (db model.Database, err error) {
	dbPostgreSql, err := gorm.Open(postgres.Open(cfg.Database.PostgreSql.GenerateConnectionPostgresString()), &gorm.Config{})
	if err != nil {
		return
	}

	db = model.Database{
		PostgreSql: dbPostgreSql,
	}

	return
}
