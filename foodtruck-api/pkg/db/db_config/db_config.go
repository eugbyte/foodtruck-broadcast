package dbconfig

import (
	"fmt"
	"strings"
)

type Driver int

const (
	Postgres Driver = iota
)

type Config struct {
	// 	use host.docker.internal instead of localhost if accessing from inside a docker container (https://stackoverflow.com/a/24326540)
	Host     string
	Port     string
	DbName   string
	User     string
	Password string
	SSLMode  string
}

// Returns a DB connection string corresponding to the specified DB, e.g. Postgres, MySQL
func ConnString(driver Driver, config Config) string {
	var mapped = map[string]string{
		"host":     config.Host,
		"port":     config.Port,
		"dbname":   config.DbName,
		"user":     config.User,
		"password": config.Password,
		"sslmode":  config.SSLMode,
	}

	var s strings.Builder

	for key, value := range mapped {
		if strings.TrimSpace(value) == "" {
			continue
		}
		s.WriteString(fmt.Sprintf("%s=%s ", key, value))
	}

	return strings.TrimSpace(s.String())
}
