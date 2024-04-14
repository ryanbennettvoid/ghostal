package memory_logger

import (
	"fmt"
	"strings"
)

type MemoryLogger struct {
	log []string
}

func NewMemoryLogger() *MemoryLogger {
	return &MemoryLogger{
		log: make([]string, 0),
	}
}

func (n *MemoryLogger) GetFullLog() string {
	if len(n.log) == 0 {
		return ""
	}
	return strings.Join(n.log, "\n")
}

func (n *MemoryLogger) Passthrough(msg string, keysAndValues ...interface{}) {
	n.log = append(n.log, fmt.Sprintf(msg, keysAndValues...))
}

func (n *MemoryLogger) Debug(msg string, keysAndValues ...interface{}) {
	n.log = append(n.log, fmt.Sprintf(msg, keysAndValues...))
}

func (n *MemoryLogger) Info(msg string, keysAndValues ...interface{}) {
	n.log = append(n.log, fmt.Sprintf(msg, keysAndValues...))
}

func (n *MemoryLogger) Warning(msg string, keysAndValues ...interface{}) {
	n.log = append(n.log, fmt.Sprintf(msg, keysAndValues...))
}

func (n *MemoryLogger) Error(msg string, keysAndValues ...interface{}) {
	n.log = append(n.log, fmt.Sprintf(msg, keysAndValues...))
}

func (n *MemoryLogger) Fatal(msg string, keysAndValues ...interface{}) {
	n.log = append(n.log, fmt.Sprintf(msg, keysAndValues...))
}
