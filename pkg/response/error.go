package response

import (
	"fmt"

	"go.uber.org/zap"
)

func ErrorGenerate(tag string, err error) error {
	return fmt.Errorf("%s -> %w", tag, err)
}

func ErrorHandler(logger *zap.Logger, tag string, err error) {
	if err == nil {
		return
	}
	if logger != nil {
		logger.Error("request error", zap.String("tag", tag), zap.Error(err))
		return
	}
	fmt.Printf("%s -> %v\n", tag, err)
}
