/*
 * Copyright 2022 Han Xin, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetLogLevelComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected logrus.Level
	}{
		{
			name:     "empty string",
			input:    "",
			expected: logrus.InfoLevel,
		},
		{
			name:     "debug level",
			input:    "debug",
			expected: logrus.DebugLevel,
		},
		{
			name:     "info level",
			input:    "info",
			expected: logrus.InfoLevel,
		},
		{
			name:     "warn level",
			input:    "warn",
			expected: logrus.WarnLevel,
		},
		{
			name:     "error level",
			input:    "error",
			expected: logrus.ErrorLevel,
		},
		{
			name:     "fatal level",
			input:    "fatal",
			expected: logrus.FatalLevel,
		},
		{
			name:     "panic level",
			input:    "panic",
			expected: logrus.PanicLevel,
		},
		{
			name:     "trace level",
			input:    "trace",
			expected: logrus.TraceLevel,
		},
		{
			name:     "uppercase level",
			input:    "DEBUG",
			expected: logrus.DebugLevel,
		},
		{
			name:     "mixed case level",
			input:    "WaRn",
			expected: logrus.WarnLevel,
		},
		{
			name:     "invalid level",
			input:    "invalid",
			expected: logrus.InfoLevel, // 默认值
		},
		{
			name:     "numeric string",
			input:    "123",
			expected: logrus.InfoLevel, // 默认值
		},
		{
			name:     "special characters",
			input:    "debug!@#",
			expected: logrus.InfoLevel, // 默认值
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetLogLevel(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetLogLevelWithCapture(t *testing.T) {
	// 测试日志输出
	oldLevel := logrus.GetLevel()
	defer logrus.SetLevel(oldLevel)
	
	// 捕获日志输出
	var logOutput []string
	hook := &testLogHook{entries: &logOutput}
	logrus.AddHook(hook)
	defer logrus.StandardLogger().ReplaceHooks(make(logrus.LevelHooks))
	
	level := GetLogLevel("debug")
	
	assert.Equal(t, logrus.DebugLevel, level)
	assert.NotEmpty(t, logOutput)
	assert.Contains(t, logOutput[0], "Set log level to debug")
}

// testLogHook 用于捕获测试期间的日志输出
type testLogHook struct {
	entries *[]string
}

func (h *testLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *testLogHook) Fire(entry *logrus.Entry) error {
	*h.entries = append(*h.entries, entry.Message)
	return nil
}

func BenchmarkGetLogLevel(b *testing.B) {
	levels := []string{"debug", "info", "warn", "error", "invalid"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		level := levels[i%len(levels)]
		GetLogLevel(level)
	}
}