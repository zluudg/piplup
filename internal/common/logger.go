package common

type Logger interface {
	Debug(string, ...any)
	Info(string, ...any)
	Warning(string, ...any)
	Error(string, ...any)
}

type FakeLogger struct{}

func (fl FakeLogger) Debug(string, ...any)   {}
func (fl FakeLogger) Info(string, ...any)    {}
func (fl FakeLogger) Warning(string, ...any) {}
func (fl FakeLogger) Error(string, ...any)   {}
