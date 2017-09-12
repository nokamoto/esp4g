package extension

import (
	"go.uber.org/zap"
	"github.com/nokamoto/esp4g/esp4g-utils"
	"fmt"
)

func sugaredLogger(config Config) (*zap.SugaredLogger, error) {
	fmt.Println(config.Logs.Zap)

	if c := config.Logs.Zap; c != nil {
		utils.Logger.Info("build zap logger")
		if logger, err := c.Build(); err != nil {
			return nil, err
		} else {
			return logger.Sugar(), nil
		}
	}

	return nil, nil
}
