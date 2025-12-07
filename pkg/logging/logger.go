// Copyright 2020 the Exposure Notifications Server authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package logging sets up and configures logging.
package logging

import (
	"context"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// contextKey is a private string type for zap logger to pass through
// context and prevent collisions in the context map
type contextKey string

// loggerKey point to zap logger pointer that store in context value
const loggerKey = contextKey("logger")

// New return zap sugered logger with development or prod config
func New(level string, development bool) *zap.SugaredLogger {
	var config *zap.Config

	if development {
		config = &zap.Config{
			Level:            zap.NewAtomicLevelAt(levelToZapLevel(level)),
			Development:      true,
			Encoding:         encodingConsole,
			EncoderConfig:    debugEncoderConfig,
			OutputPaths:      outputStderr,
			ErrorOutputPaths: outputStderr,
		}
	} else {
		config = &zap.Config{
			Level:            zap.NewAtomicLevelAt(levelToZapLevel(level)),
			Sampling:         &prodSampleConfig,
			Encoding:         encodingJSON,
			EncoderConfig:    prodEncoderConfig,
			OutputPaths:      outputStderr,
			ErrorOutputPaths: outputStderr,
		}
	}

	logger, err := config.Build()

	if err != nil {
		return NewNop()
	}

	return logger.Sugar()
}

// NewNop return no op zap sugered logger
func NewNop() *zap.SugaredLogger {
	return zap.NewNop().Sugar()
}

// NewLoggerFromEnv create new logger base on env:
//   - LOG_LEVEL: determine the minimal level that logger will displace
//   - RUN_MODE: determine if it is "development" or not
func NewLoggerFromEnv() *zap.SugaredLogger {
	level := os.Getenv("LOG_LEVEL")
	development := strings.ToLower(strings.TrimSpace(os.Getenv("RUN_MODE"))) == "development"
	return New(level, development)
}

// WithLogger create new context from old ctx, and
// store zap logger into new context
func WithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext get zap sugared logger that store in value of context.
// If no such logger, a no op logger is returned
func FromContext(ctx context.Context) *zap.SugaredLogger {
	if logger, ok := ctx.Value(loggerKey).(*zap.SugaredLogger); ok {
		return logger
	}

	return NewNop()
}

const (
	levelDebug     = "DEBUG"
	levelInfo      = "INFO"
	levelWarning   = "WARNING"
	levelError     = "ERROR"
	levelCritical  = "CRITICAL"
	levelAlert     = "ALERT"
	levelEmergency = "ENERGENCY"

	timestamp  = "timestamp"
	severity   = "severity"
	logger     = "logger"
	caller     = "caller"
	message    = "message"
	function   = "function"
	stacktrace = "stacktrace"

	encodingConsole = "console"
	encodingJSON    = "json"
)

var outputStderr = []string{"stderr"}

var prodSampleConfig = zap.SamplingConfig{
	Initial:    100,  // log fist 100
	Thereafter: 1000, // after that 100 log one
}

var prodEncoderConfig = zapcore.EncoderConfig{
	MessageKey:     message,
	LevelKey:       severity,
	TimeKey:        timestamp,
	NameKey:        logger,
	CallerKey:      caller,
	FunctionKey:    function,
	StacktraceKey:  stacktrace,
	SkipLineEnding: false,
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    levelEncoder(),
	EncodeTime:     timeEncoder(),
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

var debugEncoderConfig = zapcore.EncoderConfig{
	MessageKey:     message,
	LevelKey:       severity,
	TimeKey:        timestamp,
	NameKey:        logger,
	CallerKey:      caller,
	FunctionKey:    function,
	StacktraceKey:  stacktrace,
	SkipLineEnding: false,
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    levelEncoder(),
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.StringDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

func levelToZapLevel(s string) zapcore.Level {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case levelDebug:
		return zapcore.DebugLevel
	case levelInfo:
		return zapcore.InfoLevel
	case levelWarning:
		return zapcore.WarnLevel
	case levelError:
		return zapcore.ErrorLevel
	case levelCritical:
		return zapcore.DPanicLevel
	case levelAlert:
		return zapcore.PanicLevel
	case levelEmergency:
		return zapcore.FatalLevel
	default:
		return zapcore.WarnLevel
	}
}

func levelEncoder() zapcore.LevelEncoder {
	return func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		switch l {
		case zapcore.DebugLevel:
			enc.AppendString(levelDebug)
		case zapcore.InfoLevel:
			enc.AppendString(levelInfo)
		case zapcore.WarnLevel:
			enc.AppendString(levelWarning)
		case zapcore.ErrorLevel:
			enc.AppendString(levelError)
		case zapcore.DPanicLevel:
			enc.AppendString(levelCritical)
		case zapcore.PanicLevel:
			enc.AppendString(levelAlert)
		case zapcore.FatalLevel:
			enc.AppendString(levelEmergency)
		}
	}
}

func timeEncoder() zapcore.TimeEncoder {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC3339Nano))
	}
}
