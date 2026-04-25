package utils

type LoggerStrategy interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type LoggerContext struct {
	LoggerStrategy LoggerStrategy
}

func (l *LoggerContext) Debug(msg string, args ...any) {
	l.LoggerStrategy.Debug(msg, args...)
}

func (l *LoggerContext) Info(msg string, args ...any) {
	l.LoggerStrategy.Info(msg, args...)
}

func (l *LoggerContext) Warn(msg string, args ...any) {
	l.LoggerStrategy.Warn(msg, args...)
}

func (l *LoggerContext) Error(msg string, args ...any) {
	l.LoggerStrategy.Error(msg, args...)
}

func (l *LoggerContext) SetLoggerStrategy(strategy LoggerStrategy) {
	l.LoggerStrategy = strategy
}

func NewLoggerContext(strategy LoggerStrategy) *LoggerContext {
	return &LoggerContext{
		LoggerStrategy: strategy,
	}
}
