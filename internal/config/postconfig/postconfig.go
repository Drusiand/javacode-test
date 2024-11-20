package postconfig

import (
	"errors"
	"fmt"
	"os"
	"time"
)

const (
	envHandlerPostTimeout = "HANDLER_POST_TIMEOUT"
)

var (
	errVarMissing = errors.New("missing env variable")
)

type PostConfig struct {
	Timeout time.Duration
}

func New() (*PostConfig, error) {
	strTimeout, ok := os.LookupEnv(envHandlerPostTimeout)
	if !ok {
		return nil, fmt.Errorf("%w: %s", errVarMissing, envHandlerPostTimeout)
	}
	timeout, err := time.ParseDuration(strTimeout)
	if err != nil {
		return nil, err
	}

	cfg := &PostConfig{
		Timeout: timeout,
	}

	return cfg, nil
}
