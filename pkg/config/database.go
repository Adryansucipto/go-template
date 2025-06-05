package config

import "fmt"

type Database struct {
}

func (db *Database) GenerateConnectionPostgresString() string {
	driver := ""
	username := ""
	password := ""
	host := ""
	port := ""
	name := ""
	schema := ""

	return fmt.Sprintf("%+v://%+v:%+v@%+v:%+v/%+v?sslmode=disable&search_path=%+v", driver, username, password, host, port, name, schema)
}
