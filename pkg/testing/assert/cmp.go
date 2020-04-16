package assert

// https://darjun.github.io/2020/03/20/godailylib/go-cmp/
import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Compare(x, y interface{}) bool          { return cmp.Equal(x, y) }
func IgnoreCompare(x, y, z interface{}) bool { return cmp.Equal(x, y, cmpopts.IgnoreUnexported(z)) }

func Diff(x, y interface{}) string { return cmp.Diff(x, y) }
