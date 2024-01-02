// v0.3.0
// Author: DIEHL E.
// © Dec 2023

package msj

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a wrapper of zap.Logger
type Logger struct {
	zap.Logger
}

const (
	dockerLogs = "/logs"
)

// NewLogger creates a new logger.  If `fileName` starts with s3://, then it uses S3.  If it starts with cw://, then it uses
// CloudWatch.  Otherwise, it uses the local file system.
func NewLogger(fileName string, options ...Option) (*Logger, error) {
	o := collectOptions(options...)
	level := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.DebugLevel
	})
	if !o.debug {
		level = func(level zapcore.Level) bool {
			return level >= zapcore.InfoLevel
		}
	}
	if o.docker {
		fileName = filepath.Join(dockerLogs, fileName)
	}
	// write syncers
	stdoutSyncer := zapcore.Lock(os.Stdout)
	var cores []zapcore.Core
	if o.verbose {
		cores = append(cores, zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			stdoutSyncer, level))
	}
	if fileName != "" {
		const cPermission = 0644
		file, errF := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, cPermission)
		if errF != nil {
			return nil, errors.Wrap(errF, "createFile")
		}
		cores = append(cores, zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(file), level))
	}
	core := zapcore.NewTee(cores...)
	l := zap.New(core)
	switch o.debug {
	case true:
		l.Info("logger in debug mode")
	default:
		l.Info("logger in info mode")
	}
	return &Logger{*l}, nil
}

// Zap returns the pointer to the zap.Logger
func (l *Logger) Zap() *zap.Logger {
	return &l.Logger
}

type loggerOptions struct {
	debug   bool
	docker  bool
	verbose bool
}

// Option allows parameterizing function New
type Option func(opts *loggerOptions)

func collectOptions(options ...Option) *loggerOptions {
	opts := &loggerOptions{debug: false, docker: false, verbose: false}
	for _, option := range options {
		option(opts)
	}
	return opts
}

// WithDebug sets the logger in debug mode
func WithDebug() Option {
	return func(op *loggerOptions) {
		op.debug = true
	}
}

// WithDocker informs the logger that it runs in a container.  It stores the information
// in the repository "/logs".
func WithDocker() Option {
	return func(op *loggerOptions) {
		op.docker = true
	}
}

// WithVerbose informs the logger that it should output both in the
// file and stdout.
func WithVerbose() Option {
	return func(op *loggerOptions) {
		op.verbose = true
	}
}
