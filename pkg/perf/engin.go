package perf

import (
	"context"
	"sync"
	"time"

	"bingo/pkg/utils"
	"bingo/pkg/utils/log"
)

const (
	LimitByUndefined = "LimitByUndefined"
	LimitByAmount    = "LimitByAmount"
	LimitByTime      = "LimitByTime"
	LimitByPassrate  = "LimitByPassrate"
	LimitByCost      = "LimitByCost"
	LimitByMixed     = "LimitByMixed"
)

type Engin struct {
	// 0: limit by pv amount
	Max int

	// 1: limit by time
	PlanDuration time.Duration

	// 2: limit by pass rate
	PassRate float32

	// 3: limit by expected cost request rate
	ExpectedCost int
	ExpectedRate float32

	Sema int

	sumCode   map[int]int
	sumSucc   map[bool]int
	sumCostMS []float64
	results   [][]FnInfo

	launchedCounter, completedCounter int
	wayby                             string
	stopLaunch, closed                bool

	since       time.Time
	timer       *time.Timer
	waitCh      chan bool
	launchedCh  chan bool
	completedCh chan FnInfo
	assignCh    chan int
	// exitCh  chan bool

	context.Context
	context.CancelFunc
	once *sync.Once
}

var ep = sync.Pool{
	New: func() interface{} { return &Engin{} },
}

func newE() *Engin {
	return ep.Get().(*Engin)
	// return &Engin{}
}

// NewEnginByAmount returns
func NewEnginByAmount(max, sema int) *Engin {
	e := newE()
	e.init()
	e.Max = max
	e.Sema = sema
	e.wayby = LimitByAmount
	e.Context, e.CancelFunc = context.WithCancel(context.Background())
	return e
}

// NewEnginByDuration returns
func NewEnginByDuration(dur time.Duration, sema int) *Engin {
	e := newE()
	e.init()
	e.Sema = sema
	e.wayby = LimitByTime
	e.PlanDuration = dur
	e.Context, e.CancelFunc = context.WithTimeout(context.Background(), dur)
	return e
}
func (e *Engin) init() {
	e.Max = 0
	e.PlanDuration = 0
	e.PassRate = 0
	e.ExpectedCost = 0
	e.ExpectedRate = 0
	e.Sema = 0

	e.launchedCounter = 0
	e.completedCounter = 0
	e.closed = false
	e.stopLaunch = false

	e.since = time.Now()
	e.timer = time.NewTimer(time.Hour * 24 * 100)
	e.wayby = LimitByUndefined

	e.waitCh = make(chan bool, 10)
	e.launchedCh = make(chan bool, 1000)
	e.completedCh = make(chan FnInfo, 1000)
	e.assignCh = make(chan int, e.Sema*2)
	e.once = new(sync.Once)

	e.sumCode = make(map[int]int, 1000)
	e.sumSucc = make(map[bool]int, 1000)
	e.sumCostMS = make([]float64, 1000)
	e.results = make([][]FnInfo, 1000)

	// timer * time.Timer
}
func (e *Engin) debug() {
	// go func() {
	// 	for range time.Tick(time.Second * 10) {
	// 		if e.closed {
	// 			log.Warning("                --> exit calculator")
	// 			return
	// 		}

	// 		secs := time.Since(e.since).Seconds()
	// 		for _, arr := range e.results {
	// 			PerfCalculate(arr, secs)
	// 		}
	// 	}
	// }()

	go func() {
		for range time.Tick(time.Second * 10) {
			if e.closed {
				log.Debug("                --> exit monitor")
				return
			}
			// debug.FreeOSMemory()
			log.Debug("%s  instance: %d/%d c:%d/r:%d/w:%d/a:%d", utils.Now(), e.completedCounter, e.launchedCounter, len(e.completedCh), len(e.launchedCh), len(e.waitCh), len(e.assignCh))
		}
	}()
}
func (e *Engin) assign() {
	i := 0
	for {
		e.assignCh <- i // e.launchedCounter // 写入初始数据，保证每次执行都能写入数据
		i++
		if e.stopLaunch || e.closed {
			return
		}
	}
}
func (e *Engin) stopCheck(b bool) bool {
	if b {
		log.Debug("===========init exit %d/%d  c:%d/w:%d/r:%d==========", e.completedCounter, e.launchedCounter, len(e.completedCh), len(e.waitCh), len(e.launchedCh))
		// e.exitCh <- true
		e.closed = true
		e.CancelFunc()
		e.waitCh <- true
		e.timer.Reset(0 * time.Second)
		log.Debug("===========init exit done %d/%d c:%d/w:%d/r:%d==========", e.completedCounter, e.launchedCounter, len(e.completedCh), len(e.waitCh), len(e.launchedCh))
	}
	return b
}
func (e *Engin) run() {
	for {
		select {
		case <-e.Context.Done():
			e.stopCheck(true)
			return
		case c := <-e.launchedCh:
			if c && e.closed {
				log.Debug("                --> exit launchedCh %v/%v", c, e.closed)
			}
			e.launchedCounter++
			// log.Warning("                --> launchedCh %v/%v", c, e.launchedCounter)
			if e.wayby == LimitByAmount {
				e.stopLaunch = e.launchedCounter >= e.Max
			}
		case c := <-e.completedCh:
			e.completedCounter++
			if !c.end {
				e.results[c.index] = append(e.results[c.index], c)
			}
			log.Debug("                --> read from completedCounter chan: %v - %d c:%d/w:%d/r:%d/a:%d", c, e.completedCounter, len(e.completedCh), len(e.waitCh), len(e.launchedCh), len(e.assignCh))

			// TODO
			switch e.wayby {
			case LimitByAmount:
				if e.stopCheck(e.completedCounter >= e.Max) {
					log.Debug("                --> exit amout limitation")
					return
				}
			}
		case <-e.timer.C:
			log.Debug("                --> exit time's up")
			e.stopCheck(true)
			return
		}
	}
}
func (e *Engin) launch() {
	e.once.Do(func() {
		switch e.wayby {
		case LimitByAmount:
		case LimitByTime:
			e.timer = time.NewTimer(e.PlanDuration)
			fallthrough
		default:
			// do()
		}
		// if e.timer == nil {
		// 	e.timer = time.NewTimer(time.Hour * 24 * 100)
		// }

		log.Seperator()
		log.Debug("--------start limit: %+v------", e)
		// hi.Send("zhujunhao", "--------limit: %+v------", e)

		// e.debug()

		go e.assign()
		go e.run()
	})
}
func (e *Engin) check() bool { return !e.closed } // !<-e.exitCh }

// Wait will
func (e *Engin) Wait() {
	log.Debug("++++++++++++++++++++++++++++++ waiting.....  - c:%d/w:%d/r:%d ++++++++++++++++++++++++++++++", len(e.completedCh), len(e.waitCh), len(e.launchedCh))
	<-e.waitCh

	// close(e.exitCh)
	close(e.waitCh)
	close(e.completedCh)

	//secs := time.Since(e.since).Seconds()
	//for i, arr := range e.results {
	//	PerfCalculate(i, arr, secs)
	//}
	log.Debug("exit. cost: %s. %d/%d c:%d/r:%d/w:%d/a:%d",
		utils.FormatDuration(time.Since(e.since)),
		e.completedCounter, e.launchedCounter,
		len(e.completedCh), len(e.launchedCh), len(e.waitCh), len(e.assignCh))

	// for k, v := range e.sumSucc {
	// 	log.Notice("--------succ:%v -> completedCounter:%d", k, v)
	// }
	// for k, v := range e.sumCode {
	// 	log.Notice("=======code:%v -> completedCounter:%d", k, v)
	// }

	// hi.Send("zouxun01,zhujunhao", "exit. cost: %s. c:%d/w:%d/r:%d", utils.FormatDuration(time.Since(e.since)), len(e.completedCh), len(e.waitCh), len(e.launchedCh))

	// time.Sleep(time.Second * 5)
}

// Work 有defer方法的work
func (e *Engin) Work(para bool, deferFunc func(), fns ...func() (bool, int, time.Duration)) {
	if len(fns) == 0 {
		log.Error("please set hook functions")
		return
	}

	e.launch()

	for range fns {
		// for cnt := range fns {
		// 	e.assignCh <- cnt // 写入初始数据，保证每次执行都能写入数据
		e.results = append(e.results, []FnInfo{})
	}

	exe := func() {
		if e.stopLaunch || e.closed {
			return
		}
		defer func() {
			if deferFunc != nil {
				deferFunc()
			}
			utils.Recover()
		}()

		// log.Notice("------> running...")
		e.launchedCh <- false

		i := 0

		if len(fns) > 1 {
			idx := <-e.assignCh
			i = idx % len(fns)
		}

		succ, code, cost := fns[i]()
		if e.closed {
			return
		}

		e.completedCh <- FnInfo{
			index: i,
			end:   false,
			Succ:  succ,
			Code:  code,
			Cost:  cost,
		}

		// fmt.Printf("-----> %s  working. func:%d %d/%d c:%d/r:%d/w:%d/a:%d\n", utils.Now(), i, e.completedCounter, e.launchedCounter, len(e.completedCh), len(e.launchedCh), len(e.waitCh), len(e.assignCh))
	}

	if para {
		log.Debug("=====Paraller====")
		// utils.Paraller(e.Context, e.Sema, 0, exe)
		utils.Paraller(e.Context, e.Sema, 0, exe, e.check)
	} else {
		log.Debug("=====Concurrency====")
		utils.Timer(e.Context, e.Sema, 0, time.Second, true, exe, e.check)

		// for range time.Tick(time.Second) {
		// 	for i := 0; i < e.Sema; i++ {
		// 		go exe()
		// 	}
		// }
	}
	e.launchedCh <- true
	e.completedCh <- FnInfo{end: true}

	log.Debug("done function")
	// e.waitCh <- true
}
