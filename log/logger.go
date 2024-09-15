package log

import (
	"fmt"
	"log/slog"
	"os"
)

type Logger interface {
	Print(values ...interface{})
	Printf(format string, args ...interface{})
	Fatal(values ...interface{})
	Fatalf(format string, args ...interface{})
	Panic(values ...interface{})
	Panicf(format string, args ...interface{})
	Debug(values ...interface{})
}

type logger struct {
	log *slog.Logger
}

func New() *logger {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return &logger{
		log: log,
	}
}

func (l logger) SetAsDefault() {
	slog.SetDefault(l.log)
}

func (l logger) Print(values ...interface{}) {
	l.log.Info(fmt.Sprint(values...))
}

func (l logger) Printf(format string, args ...interface{}) {
	l.log.Info(fmt.Sprintf(format, args...))
}

func (l logger) Fatal(values ...interface{}) {
	l.log.Error(fmt.Sprint(values...))
	os.Exit(1)
}

func (l logger) Fatalf(format string, args ...interface{}) {
	l.log.Error(fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (l logger) Panic(values ...interface{}) {
	l.log.Error(fmt.Sprint(values...))
	panic(values)
}

func (l logger) Panicf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.log.Error(msg)
	panic(msg)
}

func (l logger) Debug(values ...interface{}) {
	l.log.Info(fmt.Sprint(values...))
}

func Print(values ...interface{}) {
	slog.Info(fmt.Sprint(values...))
}

func Printf(format string, args ...interface{}) {
	slog.Info(fmt.Sprintf(format, args...))
}

func Fatal(values ...interface{}) {
	slog.Error(fmt.Sprint(values...))
	os.Exit(1)
}

func Fatalf(format string, args ...interface{}) {
	slog.Error(fmt.Sprintf(format, args...))
	os.Exit(1)
}

func Panic(values ...interface{}) {
	slog.Error(fmt.Sprint(values...))
	panic(values)
}

func Panicf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	slog.Error(msg)
	panic(msg)
}

func Debug(values ...interface{}) {
	slog.Info(fmt.Sprint(values...))
}
