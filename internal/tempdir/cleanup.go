package tempdir

import (
	"os"
	"os/signal"
	"syscall"
)

func init() {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		if singleton != nil {
			singleton.Delete()
		}
		os.Exit(0)
	}()
}
