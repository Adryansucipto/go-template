package config

import "fmt"

type RedisConfig struct {
	Enabled         bool   `yaml:"enabled"`
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	Password        string `yaml:"password"`
	DB              int    `yaml:"db"`
	DefaultTTLMin   int    `yaml:"default_ttl_minutes"`
}

func (r *RedisConfig) Addr() string {
	host := r.Host
	if host == "" {
		host = "localhost"
	}
	port := r.Port
	if port == "" {
		port = "6379"
	}
	return fmt.Sprintf("%s:%s", host, port)
}

