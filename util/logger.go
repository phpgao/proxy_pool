package util

import (
	"fmt"
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/fatih/structs"
	"io"
	"os"
	"sync"
	"time"
)

var logger *log.Logger

type Handler struct {
	mu     sync.Mutex
	Writer io.Writer
}

func (h *Handler) HandleLog(e *log.Entry) error {
	color := text.Colors[e.Level]
	level := text.Strings[e.Level]
	names := e.Fields.Names()

	h.mu.Lock()
	defer h.mu.Unlock()

	_, _ = fmt.Fprintf(h.Writer, "\033[%dm%6s\033[0m[%s] %-25s", color, level, time.Now().Format("2006-01-02 15:04:05"), e.Message)

	for _, name := range names {
		_, _ = fmt.Fprintf(h.Writer, " \033[%dm%s\033[0m=%v", color, name, e.Fields.Get(name))
	}

	_, _ = fmt.Fprintln(h.Writer)

	return nil
}

func GetLogger() *log.Logger {
	if logger == nil {
		var level log.Level
		if ServerConf.Debug {
			level = log.DebugLevel
		} else {
			level = log.InfoLevel
		}
		logger = &log.Logger{
			Handler: &Handler{
				Writer: os.Stderr,
			},
			Level: level,
		}
		logger.WithField("config", structs.Map(ServerConf)).Debug("loaded config")
	}

	return logger
}
