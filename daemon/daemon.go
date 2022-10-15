package daemon

import (
	"log"
	"os"
	"reflect"
	"sync"
	"time"
)

// StructDaemon :
type StructDaemon struct {
	SleepTime int
	WaitGroup bool
	Logger    *log.Logger
	ListFunc  []StructFunction
}

// StructFunction :
type StructFunction struct {
	Function func() error
}

// New : SleepTime default 5
func New() *StructDaemon {
	return &StructDaemon{
		SleepTime: 5,
		Logger:    log.New(os.Stderr, "", log.LstdFlags),
	}
}

// StartDaemon :
func (p *StructDaemon) StartDaemon() {
	p.Logger.Println("| Daemon | Starting")

	started := 0
Start:
	if started == 0 {
		p.Logger.Printf("| Daemon | Started | Configured to sleep every %d seconds\n", p.SleepTime)
		started = 1
	}

	if p.WaitGroup == true {
		p.ExecuteHandlersWg()
	} else {
		p.ExecuteHandlers()
	}

	time.Sleep(time.Duration(p.SleepTime) * time.Second)

	goto Start
}

// ExecuteHandlers :
func (p *StructDaemon) ExecuteHandlers() {
	for _, structFunction := range p.ListFunc {
		go func(currentFunction func() error) {
			currentFunction()
		}(structFunction.Function)
	}
}

// ExecuteHandlersWg :
func (p *StructDaemon) ExecuteHandlersWg() {
	wg := &sync.WaitGroup{}
	for _, structFunction := range p.ListFunc {
		wg.Add(1)
		go func(currentFunction func() error) {
			currentFunction()
			wg.Done()
		}(structFunction.Function)
	}
	wg.Wait()
}

// AddHandler :
func (p *StructDaemon) AddHandler(f func() error) {
	p.ListFunc = append(p.ListFunc, StructFunction{
		f,
	})
}

// StructNameToString :
func (p *StructDaemon) StructNameToString(myvar interface{}) (res string) {
	t := reflect.TypeOf(myvar)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		res += "*"
	}

	return res + t.Name()
}
