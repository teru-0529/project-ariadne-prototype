package wails

type Logger interface {
	Error(msg string, args ...any)
}
