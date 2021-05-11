package daemon

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
)

func TestNewWorker(t *testing.T) {
	wg := new(sync.WaitGroup)
	cancelDaemon := InitDaemon(wg)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGTSTP, syscall.SIGSTOP)
	s := <-c
	fmt.Println("signal.Notify", s)

	fmt.Println("cancel ...")
	cancelDaemon()
	fmt.Println("wg wait ...")
	wg.Wait()
	fmt.Println("main end ...")
}
