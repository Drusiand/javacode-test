package dbconfig

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

const (
	envDBUser = "DB_USER"
	envDBPass = "DB_PASS"
	envDBHost = "DB_HOST"
	envDBPort = "DB_PORT"
	envDBName = "DB_NAME"
)

var (
	errVarMissing = errors.New("missing env variable")
	errBadPort    = errors.New("bad db port")
)

type DBConfig struct {
	DBUser string
	DBPass string
	DBHost string
	DBPort int
	DBName string
}

func New() (*DBConfig, error) {

	dbUser, ok := os.LookupEnv(envDBUser)
	if !ok {
		return nil, fmt.Errorf("%w: %s", errVarMissing, envDBUser)
	}

	dbPass, ok := os.LookupEnv(envDBPass)
	if !ok {
		return nil, fmt.Errorf("%w: %s", errVarMissing, envDBPass)
	}

	dbHost, ok := os.LookupEnv(envDBHost)
	if !ok {
		return nil, fmt.Errorf("%w: %s", errVarMissing, envDBHost)
	}

	strDBPort, ok := os.LookupEnv(envDBPort)
	if !ok {
		return nil, fmt.Errorf("%w: %s", errVarMissing, envDBPort)
	}
	dbPort, err := strconv.Atoi(strDBPort)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errBadPort, strDBPort)
	}

	dbName, ok := os.LookupEnv(envDBName)
	if !ok {
		return nil, fmt.Errorf("%w: %s", errVarMissing, envDBName)
	}

	cfg := &DBConfig{
		DBUser: dbUser,
		DBPass: dbPass,
		DBHost: dbHost,
		DBPort: dbPort,
		DBName: dbName,
	}

	return cfg, nil
}
