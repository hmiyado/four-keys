package cli

import (
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
)

type CliContextWrapper struct {
	context *cli.Context
}

func (c *CliContextWrapper) isDebug() bool {
	return c.context.Bool("debug")
}

func (c *CliContextWrapper) Debugf(format string, a ...any) {
	debug := c.context.Bool("debug")
	if debug {
		fmt.Print("[Debug] ")
		fmt.Printf(format, a...)
	}
}

func (c *CliContextWrapper) Debugln(a ...any) {
	debug := c.context.Bool("debug")
	if debug {
		fmt.Print("[Debug] ")
		fmt.Println(a...)
	}
}

func (c *CliContextWrapper) StartTimer(key string) {
	if c.isDebug() {
		if timerMap == nil {
			timerMap = make(map[string]time.Time)
		}
		timerMap[key] = time.Now()
		c.Debugln("StartTimer:", key)
	}
}

func (c *CliContextWrapper) StopTimer(key string) {
	if c.isDebug() {
		c.Debugln("Stop_Timer:", key, "\t", time.Since(timerMap[key]))
		delete(timerMap, key)
	}
}

func (c *CliContextWrapper) Error(err error) {
	c.context.App.ErrWriter.Write([]byte(err.Error()))
}

func (c *CliContextWrapper) Write(p []byte) {
	c.context.App.Writer.Write(p)
}
