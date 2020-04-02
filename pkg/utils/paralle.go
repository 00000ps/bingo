package utils

import (
	"bingo/pkg/log"
	"context"
	"sync"
	"time"
)

var (
	loopmap = NewMap()
)

type Semaphore struct {
	c           chan struct{}
	batch, leng int
	wg          *sync.WaitGroup
}

func NewSemaphore(batch, length int) *Semaphore {
	s := &Semaphore{
		c:     make(chan struct{}, batch),
		batch: batch,
		leng:  length,
		wg:    new(sync.WaitGroup),
	}
	s.wg.Add(length)
	return s
}

func (s *Semaphore) Add(delta int) {
	// s.wg.Add(delta)
	for i := 0; i < delta; i++ {
		s.c <- struct{}{}
	}
}

func (s *Semaphore) Done() {
	<-s.c
	s.wg.Done()
}

func (s *Semaphore) Wait() {
	s.wg.Wait()
	close(s.c)
}

type Parall struct {
	c           chan struct{}
	para, count int
}

func NewParall(parall int) *Parall {
	s := &Parall{
		c:    make(chan struct{}, parall),
		para: parall,
	}
	return s
}

func (s *Parall) Add(delta int) {
	for i := 0; i < delta; i++ {
		s.c <- struct{}{}
	}
}

func (s *Parall) Done() { <-s.c }

// Paraller returns
func Paraller(ctx context.Context, batch, max int, fn func(), condition ...func() bool) {
	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		return
	// 	default:
	// 	}
	// }
	if max > 0 {
		s := NewSemaphore(batch, max)
		for i := 0; i < max; i++ {
			if len(condition) > 0 && !condition[0]() {
				log.Debug("exit Paraller by condition")
				return

				// // 结束信号为true，消耗掉剩余的次数
				// defer s.Done()
				// s.Add(1)
				// continue
			}
			go func() {
				s.Add(1)
				fn()
				s.Done()
			}()
		}
		s.Wait()
	} else {
		ch := make(chan bool, batch)
		// log.Notice("Paraller batch:%d", batch)
		for len(condition) == 0 || condition[0]() {
			// log.Debug("Paraller 11")
			ch <- true
			// log.Notice("Paraller 12")
			go func() {
				fn()
				// log.Notice("----Paraller 13")
				<-ch
				// log.Debug("----Paraller 14")
			}()
		}
		close(ch)
	}
}

// Timer will do something periodly
func Timer(ctx context.Context, batch, max int, d time.Duration, para bool, fn func(), condition ...func() bool) {
	count := 0
	exitCh := make(chan bool, 2)
	ch := make(chan bool, batch*2)
	exitCh <- false

	do := func() {
		fn()
		ch <- false
	}

	go func() {
		for range time.Tick(d) {
			if exit := <-exitCh; exit || (len(condition) > 0 && !condition[0]()) {
				log.Warning("exit Timer by condition")
				//ch <- true
				return
			}
			exitCh <- false
			for i := 0; i < batch; i++ {
				//para、qps方式
				if para {
					go do()
				} else {
					//普通case
					do()
				}
				if len(condition) > 0 && !condition[0]() {
					log.Warning("exit Timer by condition")
					ch <- true
					close(ch)
					return
				}
			}
		}
	}()
	for {
		select {
		case <-ch:
			count++
			if (max > 0 && count > max) || (len(condition) > 0 && !condition[0]()) {
				// exitCh <- true
				//在这里关channel会send on closed channel
				//close(ch)
				log.Warning("exit Timer by condition, self count")
				return
			}
		}
	}
}

// Looper used to do specified function periodly,
// if condition is set, it will exit when condition returns false
func Looper(period time.Duration, parallel, doFirst bool, condition func() bool, do func()) int {
	id := int(GetID())
	f := func(donow bool) {
		loopmap.Store(id, false)
		for {
			v, _ := loopmap.Load(id)
			exit := v.(bool)

			if !donow {
				time.Sleep(period)
			}

			if (condition != nil && !condition()) || exit {
				loopmap.Delete(id)
				log.Debug("utils.Looper %d-%v: exit", id, exit)
				return
			}

			log.Debug("utils.Looper %d-%v: doing...", id, exit)
			do()

			if donow {
				time.Sleep(period)
			}
		}
	}
	if parallel {
		go f(doFirst)
	} else {
		f(doFirst)
	}
	return id
}

// KillLooper ss
func KillLooper(id int) {
	_, ok := loopmap.Load(id)
	if ok {
		loopmap.Store(id, true)
		log.Debug("utils.Looper %d-%v: kill", id, true)
	}
}

// Waiter used to do specified function once
// only when condition returns true
func Waiter(period time.Duration, parallel, doFirst bool, condition func() bool, do func()) {
	f := func(donow bool) {
		for {
			if !donow {
				time.Sleep(period)
			}

			if condition() {
				do()
				return
			}

			if donow {
				time.Sleep(period)
			}
		}
	}
	if parallel {
		go f(doFirst)
	} else {
		f(doFirst)
	}
}

// ParallelDo used in scenarios which need parallel working
func ParallelDo(arr []interface{}, do func(int, interface{}), done ...func()) {
	var wg sync.WaitGroup
	wg.Add(len(arr))
	for i, v := range arr {
		go func(i int, v interface{}) {
			defer wg.Done() // 操作完成，减少一个计数
			do(i, v)
		}(i, v)
	}
	wg.Wait() // 等待，直到计数为0

	for _, f := range done {
		f()
	}
}

// OR returns
func OR(input interface{}, selection ...interface{}) bool {

	return true
}

// AND returns
func AND(input interface{}, selection ...interface{}) bool {
	return true
}
