package sigwait

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/fe0b6/tools"
	"gogs.3l8.ru/dnk/golog"
)

var (
	selfExitChan  = make(chan struct{})
	exitChan      = make(chan struct{})
	waitTime      = 5
	ignoreSignals []string
	wg            sync.WaitGroup
)

func init() {
	go runWaiter()
}

func runWaiter() {
	wg.Add(1)
	// Перехват сигналов
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)

	waitExit(c)
	close(exitChan)
	wg.Done()

	go func() {
		time.Sleep(time.Duration(waitTime) * time.Second)
		golog.Error("Неудалось завершить работу корректно")

		os.Exit(2)
	}()
}

// RunWaiter - Включаем ожидание выхода
func RunWaiter() {
	wg.Wait()
	golog.Info("Работа завершена корректно")
}

func waitExit(c chan os.Signal) {
	for {
		select {
		case s := <-c:
			if !tools.InArray(ignoreSignals, s.String()) {
				golog.Info("Получен сигнал:", s)
				return
			}

		case <-selfExitChan:
			golog.Info("Самоинициализированный выход")
			return
		}
	}
}

// SetIgnoreSignal - указываем какие сигналы игнорируем
func SetIgnoreSignal(arr []string) {
	ignoreSignals = arr
}

// SetWaitTime - указываем какое время ждать перед выходом
func SetWaitTime(t int) {
	waitTime = t
}

// Wait - ожидаем команды на выход
func Wait() {
	AddWait()
	_ = <-exitChan
}

// AddWait - Добавляем счетчик работающих потоков
func AddWait() {
	wg.Add(1)
}

// Release - отмечаем что поток готов к выходу
func Release() {
	wg.Done()
}

// CheckExited - Если программа завершается, то вернет true
func CheckExited() (ok bool) {
	if tools.IsClosedChan(exitChan) {
		ok = true
	}
	return
}

// Exit - функция корректного выхода
func Exit() {
	close(selfExitChan)
}
