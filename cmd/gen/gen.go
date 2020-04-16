package gen

// https://darjun.github.io/2020/03/14/godailylib/jennifer/
import (
	"fmt"
	"strconv"
	"strings"

	. "github.com/dave/jennifer/jen"
)

// TestCase used to generate test case automated
func TestCase(pkg, feature string, id int) {
	f := NewFilePathName("testshop/"+pkg, pkg)
	f.Func().Id("Test_"+strings.ToTitle(feature)+"_"+strconv.Itoa(id)).Params().Block(
		Var().Id("tempInt").Int(),
		Qual("fmt", "Println").Call(Lit("Hello, world")),
		Qual("bingo/pkg/utils", "Max").Call(Lit(1), Lit(-10), Lit(8738914)),
	)
	fmt.Printf("%#v\n", f)
}
