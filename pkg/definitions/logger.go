package definitions

type ILogger interface {
	Debug(msg string, keysAndValues ...interface{})

	Info(msg string, keysAndValues ...interface{})

	Warning(msg string, keysAndValues ...interface{})

	Error(msg string, keysAndValues ...interface{})

	Fatal(msg string, keysAndValues ...interface{})
}
