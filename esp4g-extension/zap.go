package extension

import (
	"go.uber.org/zap"
	"github.com/nokamoto/esp4g/esp4g-utils"
	"github.com/nokamoto/esp4g/esp4g-extension/config"
)

func sugaredLogger(cfg config.ExtensionConfig) (*zap.SugaredLogger, error) {
	if c := cfg.Logs.Zap; c != nil {
		utils.Logger.Info("build zap logger")
		if logger, err := c.Build(); err != nil {
			return nil, err
		} else {
			return logger.Sugar(), nil
		}
	}
	return nil, nil
}
