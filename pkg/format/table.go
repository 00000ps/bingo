package format

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func Do() {
	t := &table{}
	t.setColumn(10)
	t.setAlign(AlignMiddle)
	//t.setBlank("-")
	t.setPrintIndex()
	t.AppendRow("jahfj", "wjhjeiu", "jhweuiehwu", "jewhfie")
	t.AppendRow("jahfj222", "wjhjeiu", "jhweuiehwu", "jewhfie")
	t.AppendRow("jahfj3", "wjhjeiuwqdddd", "jhweuiehwu", "jewhfie")
	t.AppendRow("jahfj12", "wjhjeiu", "jhweuiehwuqw", "jewhfiewqqqqqqqqqqqqqqqqqqqqqqqq")
	t.AppendRow("jahfj3", "wjhjeiuwqdddd", "jhweuiehwu", "jewhfie")
	t.AppendRow("jahfj3")
	t.AppendRow("jahfj3", "wjhjeiuwqdddd", "", "jewhfie")
	t.AppendRow("jahfj3", "  ", "jewhfie", "", "", "s")
	t.AppendRow("jahfj3", "  ", "jewhfie", "", "", "", "", "")
	t.AppendRow("jahfj3", "wjhjeiuwqdddd", "jhwe2121uiehwu", "jewhfie")
	t.AppendRow("jahfj3", "wjhjeiuwqdddd", "4443", "jewhfie")
	t.AppendRow("jahfj31", "wjhjeiuwqdddd", "4443", "jewhfie")
	t.AppendRow("jah", "wjhjeiuwqdddd", "4443", "jewhfie")
	t.AppendRow("ja5554y45y45y54hfj3", "wjhjeiuwqdddd", "4443", "jewhfie21e21")
	t.Output()
}

const (
	AlignLeft = iota
	AlignRight
	AlignMiddle
)

type column struct {
	length int
	rows   []string
}
type table struct {
	hasHead     bool
	printIdx    bool
	blank       string
	align       int
	columnCount int
	//colLen  int
	columns []column
	rows    []string
}

func getIntLen(num int) int           { return int(math.Log10(float64(num))) + 1 }
func (t *table) setBlank(sign string) { t.blank = sign }
func (t *table) setAlign(sign int)    { t.align = sign }
func (t *table) setPrintIndex()       { t.printIdx = true }
func (t *table) setColumn(num int) {
	t.columnCount = num
	for i := 0; i < t.columnCount; i++ {
		t.columns = append(t.columns, column{})
	}
}
func NewTable(align int, column int, printIndex bool) *table {
	t := &table{}
	t.setColumn(column)
	t.setAlign(align)
	//t.setBlank("-")
	if printIndex {
		t.setPrintIndex()
	}
	return t
}
func New(align int, printIndex bool, txt ...string) *table {
	t := &table{}
	t.setAlign(align)
	//t.setBlank("-")
	if printIndex {
		t.setPrintIndex()
	}
	t.SetHead(txt...)
	return t
}
func (t *table) AppendRow(txt ...string) {
	add := func(i int, v string) {
		if i < t.columnCount {
			t.columns[i].rows = append(t.columns[i].rows, v)
			//fmt.Printf("new=%d old=%d v=%s\n", len(v), t.columns[i].length, v)
			if len(v) > t.columns[i].length {
				t.columns[i].length = len(v)
			}
		}
	}
	if len(txt) > t.columnCount {
		fmt.Println("error, oversize values")
	}
	for i, v := range txt {
		add(i, v)
	}
	if len(txt) < t.columnCount {
		for i := len(txt); i < t.columnCount; i++ {
			add(i, t.blank)
		}
	}
}
func (t *table) SetHead(txt ...string) {
	if t == nil {
		t = &table{}
		t.setAlign(AlignLeft)
	}
	t.setColumn(len(txt))
	t.AppendRow(txt...)
}
func (t *table) Output() {
	t.GetContent()
	for _, r := range t.rows {
		fmt.Println(r)
	}
}
func (t *table) GetContent() []string {
	t.rows = []string{}
	if len(t.columns) > 0 {
		signCross := "+"
		signHorLine := "-"
		signVerLine := "|"

		ht := "+"

		rowCount := len(t.columns[0].rows)
		indexLen := 0
		if rowCount > 0 {
			indexLen = getIntLen(rowCount)
			if t.printIdx {
				ht += strings.Repeat(signHorLine, indexLen) + signCross
			}
		}

		for i := 0; i < t.columnCount; i++ {
			ht += strings.Repeat(signHorLine, t.columns[i].length+2) + signCross
		}
		//fmt.Printf("c=%d %#v\n", len(t.columns), t)
		t.rows = append(t.rows, ht)
		//fmt.Println(ht)
		for ri := 0; ri < rowCount; ri++ {
			row := signVerLine
			if t.printIdx {
				num := ri + 1
				if t.hasHead {
					num = ri
				}
				row += strings.Repeat(" ", indexLen-getIntLen(num)) + strconv.Itoa(num) + signVerLine
			}
			for ci := 0; ci < t.columnCount; ci++ {
				vlen := len(t.columns[ci].rows[ri])
				//fmt.Printf("r=%d c=%d v=%s\n", ri, ci, t.columns[ci].rows[ri])
				row += " "
				switch t.align {
				case AlignRight:
					row += strings.Repeat(" ", t.columns[ci].length-vlen) + t.columns[ci].rows[ri]
				case AlignMiddle:
					spaceLen := t.columns[ci].length - vlen
					spaceLeft := spaceLen / 2
					row += strings.Repeat(" ", spaceLeft) + t.columns[ci].rows[ri] + strings.Repeat(" ", spaceLen-spaceLeft)
				case AlignLeft:
					fallthrough
				default:
					row += t.columns[ci].rows[ri] + strings.Repeat(" ", t.columns[ci].length-vlen)
				}
				row += " " + signVerLine
			}
			//fmt.Println(row)
			//fmt.Println(ht)
			t.rows = append(t.rows, row)
			t.rows = append(t.rows, ht)
		}
	}
	return t.rows
}
