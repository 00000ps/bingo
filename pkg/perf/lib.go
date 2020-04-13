package perf

import (
	"fmt"
	"sort"
	"time"

	"bingo/pkg/log"
	"bingo/pkg/utils"
)

type FnInfo struct {
	Succ, end   bool
	Code, index int
	Cost        time.Duration
}
type PVR struct {
	P int     // position
	V float64 // value
	R float64 // rate
}
type PerfData struct {
	PV                           int
	QPS                          float64
	Pass                         int
	Passrate                     float64
	Code                         map[int]int
	CostMSQuantile               map[float32]PVR
	AvgMS, MinMS, MaxMS, TotalMS float64
	costMS                       []float64
}

func PerfCalculate(i int, arr []FnInfo, costSec float64, quantiles ...float64) {
	// if len(l.results) <= idx {
	// 	log.Error("invalid index input")
	// 	return
	// }
	// arr := l.results[idx]
	if len(arr) == 0 {
		log.Error("no func info founded")
		return
	}

	var pd PerfData

	start := time.Now()
	pd.PV = len(arr)
	pd.QPS = float64(pd.PV) / costSec
	pd.Code = make(map[int]int)
	pd.Pass = 0

	for _, info := range arr {
		pd.Code[info.Code]++
		if info.Succ {
			pd.Pass++
		}
		c := info.Cost.Seconds() * 1000
		pd.TotalMS += c
		pd.costMS = append(pd.costMS, c)
	}
	sort.Float64s(pd.costMS)
	pd.MinMS = pd.costMS[0]
	pd.MaxMS = pd.costMS[pd.PV-1]
	pd.AvgMS = pd.TotalMS / float64(pd.PV)
	pd.Passrate = float64(pd.Pass) / float64(pd.PV)

	quantile := func(slice []float64, accuracy float32) int {
		if accuracy > 1 {
			log.Error("invalid input")
			return -1
		}

		leng := len(slice)
		if !sort.Float64sAreSorted(slice) {
			sort.Float64s(slice)
		}
		pos := int(float32(leng)*accuracy) - 1
		if pos < leng-1 {
			return pos
		}
		return leng - 1
	}

	pd.CostMSQuantile = make(map[float32]PVR)

	accs := []float64{0.5, 0.6, 0.7, 0.8, 0.9, 0.99, 0.999, 0.9999, 0.99999, 0.99999, 0.999999}
	tm := make(map[float64]bool)
	accs = append(accs, quantiles...)
	for _, a := range accs {
		tm[a] = true
	}
	accs = []float64{}
	for a := range tm {
		accs = append(accs, a)
	}
	sort.Float64s(accs)
	lp := 0
	for _, a := range accs {
		v := float32(a)
		pos := quantile(pd.costMS, v)
		if pos == lp {
			break
		}
		lp = pos
		pd.CostMSQuantile[v] = PVR{
			P: pos,
			V: pd.costMS[pos],
			R: float64(pos+1) / float64(pd.PV),
		}
	}
	pd.costMS = nil
	// d, err := json.MarshalIndent(&pd, "", "  ")
	// // return string(d)
	// fmt.Printf("calculate cost: %s %s err:%s\n", utils.Since(start), d, err)
	str := fmt.Sprintf("func %d: calculate cost: %s %+v", i, utils.Since(start), pd)
	fmt.Println(str)
}
