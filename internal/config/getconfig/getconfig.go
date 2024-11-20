package getconfig

import (
	"errors"
	"fmt"
	"os"
	"time"
)

const (
	envHandlerGetTimeout = "HANDLER_GET_TIMEOUT"
)

var (
	errVarMissing = errors.New("missing env variable")
)

type GetConfig struct {
	Timeout time.Duration
}

func New() (*GetConfig, error) {
	strTimeout, ok := os.LookupEnv(envHandlerGetTimeout)
	if !ok {
		return nil, fmt.Errorf("%w: %s", errVarMissing, envHandlerGetTimeout)
	}
	timeout, err := time.ParseDuration(strTimeout)
	if err != nil {
		return nil, err
	}

	cfg := &GetConfig{
		Timeout: timeout,
	}

	return cfg, nil
}
