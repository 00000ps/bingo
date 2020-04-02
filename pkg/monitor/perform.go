package monitor

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"bingo/pkg/log"
)

func Perform() {
	cs, _ := cpu.Info()
	for _, c := range cs {
		log.Notice("%v", c)
	}
	cts, _ := cpu.Times(false)
	for _, c := range cts {
		log.Notice("%v", c)
	}
	v, _ := mem.VirtualMemory()
	// almost every return value is a struct
	log.Info("Total: %v, Free:%v, UsedPercent:%f%%", v.Total, v.Free, v.UsedPercent)
	// convert to JSON. String() is also implemented
	log.Notice("%v", v)
	list, _ := process.Pids()
	//fmt.Println(list)
	for i, id := range list {
		p, _ := process.NewProcess(id)
		name, _ := p.Name()
		pc, _ := p.CPUPercent()
		pm, _ := p.MemoryInfo()
		st, _ := p.Status()
		log.Notice("%d id:%d name:%s cpu:%f mem:%s status:%s\n", i, id, name, pc, pm.String(), st)
	}
}
