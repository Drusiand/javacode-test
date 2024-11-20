package httpconfig

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	envHTTPHost         = "HTTP_HOST"
	envHTTPPort         = "HTTP_PORT"
	envHTTPReadTimeout  = "HTTP_R_TIMEOUT"
	envHTTPWriteTimeout = "HTTP_W_TIMEOUT"
	envHTTPIdleTimeout  = "HTTP_I_TIMEOUT"
)

var (
	errVarMissing = errors.New("missing env variable")
	errBadPort    = errors.New("bad http port")
)

type HTTPConfig struct {
	HTTPHost         string
	HTTPPort         int
	HTTPReadTimeout  time.Duration
	HTTPWriteTimeout time.Duration
	HTTPIdleTimeout  time.Duration
}

func New() (*HTTPConfig, error) {

	httpHost, ok := os.LookupEnv(envHTTPHost)
	if !ok {
		return nil, fmt.Errorf("%w: %s", errVarMissing, envHTTPHost)
	}

	strHTTPPort, ok := os.LookupEnv(envHTTPPort)
	if !ok {
		return nil, fmt.Errorf("%w: %s", errVarMissing, envHTTPPort)
	}
	httpPort, err := strconv.Atoi(strHTTPPort)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errBadPort, strHTTPPort)
	}

	strHTTPReadTimeout, ok := os.LookupEnv(envHTTPReadTimeout)
	if !ok {
		return nil, fmt.Errorf("%w: %s", errVarMissing, envHTTPReadTimeout)
	}
	httpReadTimeout, err := time.ParseDuration(strHTTPReadTimeout)
	if err != nil {
		return nil, err
	}

	strHTTPWriteTimeout, ok := os.LookupEnv(envHTTPWriteTimeout)
	if !ok {
		return nil, fmt.Errorf("%w: %s", errVarMissing, envHTTPWriteTimeout)
	}
	httpWriteTimeout, err := time.ParseDuration(strHTTPWriteTimeout)
	if err != nil {
		return nil, err
	}

	strHTTPIdleTimeout, ok := os.LookupEnv(envHTTPIdleTimeout)
	if !ok {
		return nil, fmt.Errorf("%w: %s", errVarMissing, envHTTPIdleTimeout)
	}
	httpIdleTimeout, err := time.ParseDuration(strHTTPIdleTimeout)
	if err != nil {
		return nil, err
	}

	cfg := &HTTPConfig{
		HTTPHost:         httpHost,
		HTTPPort:         httpPort,
		HTTPReadTimeout:  httpReadTimeout,
		HTTPWriteTimeout: httpWriteTimeout,
		HTTPIdleTimeout:  httpIdleTimeout,
	}

	return cfg, nil
}
