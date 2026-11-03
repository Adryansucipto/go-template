package config

import (
	"fmt"
	"net/url"
)

type Database struct {
	PostgreSql DatabaseCfg `yaml:"postgreSql"`
}

type DatabaseCfg struct {
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	DBName          string `yaml:"name"`
	Driver          string `yaml:"driver"`
	Schema          string `yaml:"schema"`
	SSLMode         string `yaml:"sslmode"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime_minutes"`
}

func (db *DatabaseCfg) GenerateConnectionPostgresString() string {
	driver := db.Driver
	if driver == "" {
		driver = "postgresql"
	}
	username := url.QueryEscape(db.Username)
	password := url.QueryEscape(db.Password)
	host := db.Host
	port := db.Port
	name := db.DBName
	schema := db.Schema
	if schema == "" {
		schema = "public"
	}
	sslmode := db.SSLMode
	if sslmode == "" {
		sslmode = "disable"
	}

	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s&search_path=%s",
		driver, username, password, host, port, name, sslmode, schema)
}
