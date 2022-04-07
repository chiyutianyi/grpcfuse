package utils_test

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/sirupsen/logrus"

	"github.com/chiyutianyi/grpcfuse/pkg/utils"
)

func TestGetLogLevel(t *testing.T) {
	// default
	assert.Equal(t, logrus.InfoLevel, utils.GetLogLevel(""))
	assert.Equal(t, logrus.DebugLevel, utils.GetLogLevel("debug"))
}
