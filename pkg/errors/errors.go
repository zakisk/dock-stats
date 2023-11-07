package errors

import (
	"github.com/zakisk/dock-stats/pkg/logger"
)

func ErrDockerNotInstalled() *logger.Message {
	return &logger.Message{
		Message: "Unable to access `docker`",
		KeyValues: map[string]any{
			logger.KeyError:         "Encountered an error while accessing `docker`",
			logger.KeyProbableCause: "Docker may not be installed on your machine or not accessible from terminal.",
			logger.KeyRemedy:        "If `docker` is installed but not accessible from terminal then add its location in `PATH` environment variable.",
		},
	}
}
