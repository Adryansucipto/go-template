package config

import (
	"fmt"
	"net/url"
)

type Database struct {
	PostgreSql DatabaseCfg `yaml:"postgreSql"`
}

type DatabaseCfg struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"name"`
	Driver   string `yaml:"driver"`
	Schema   string `yaml:"schema"`
}

func (db *DatabaseCfg) GenerateConnectionPostgresString() string {
	driver := db.Driver
	username := url.QueryEscape(db.Username)
	password := url.QueryEscape(db.Password)
	host := db.Host
	port := db.Port
	name := db.DBName
	schema := db.Schema

	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		driver, username, password, host, port, name, schema)
}
