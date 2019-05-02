package sigwait

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/fe0b6/tools"
)

var (
	exitChan      chan struct{}
	waitTime      int
	ignoreSignals []string
	wg            sync.WaitGroup
)

func init() {
	waitTime = 5

	// Создаем канал, к которому будут подключаться в ожидании выхода
	exitChan = make(chan struct{})

	go runWaiter()
}

func runWaiter() {
	// Перехват сигналов
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)

	waitExit(c)

	go func() {
		time.Sleep(time.Duration(waitTime) * time.Second)
		log.Println("[error]", "Неудалось завершить работу корректно")

		os.Exit(2)
	}()

	wg.Done()

	log.Println("[info]", "Работа завершена корректно")

	os.Exit(0)
}

func waitExit(c chan os.Signal) {
	for {
		select {
		case s := <-c:
			if !tools.InArray(ignoreSignals, s.String()) {
				log.Println("[info]", "Получен сигнал: ", s)
				return
			}

		case <-exitChan:
			log.Println("[info]", "Самоинициализированный выход")
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
	wg.Add(1)
	_ = <-exitChan
}

// Release - отмечаем что поток готов к выходу
func Release() {
	wg.Done()
}

// Exit - функция корректного выхода
func Exit() {
	close(exitChan)
}
