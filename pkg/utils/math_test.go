package utils

import (
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Min(t *testing.T) {
	type args struct {
		v []interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{name: "", args: args{v: []interface{}{0, 0.1, 0, 1000.87896}}, want: 0},
		{name: "", args: args{v: []interface{}{0, 0.1, 0, -1.0, -1, 1000.87896}}, want: -1.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Min(tt.args.v...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Min() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_Max(t *testing.T) {
	type args struct {
		v []interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{name: "", args: args{v: []interface{}{0, 0.1, 0, 1000.87896}}, want: 1000.87896},
		{name: "", args: args{v: []interface{}{0, 0.1, 0, -1.0, -1, 1000.87896}}, want: 1000.87896},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Max(tt.args.v...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkMax(b *testing.B) {
	var n int
	for i := 0; i < b.N; i++ {
		n++
		Max(0, 0.1, 0, -1.0, -1, 1000.87896)
	}
}

func TestAbs(t *testing.T) {
	Convey("as", t, func() {
		t.Log("dei")
		t.Log(12211413)
		t.Log(true)
		t.Log("dqdwqdwqdwqd")
		t.Log("dwddwdqwkdjqwidjei")
		
		So(Abs(-1), ShouldEqual, 1)
	})

	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{"", args{0}, 0},
		{"", args{0.0}, 0.0},
		{"", args{1}, 1},
		{"", args{-1}, 1},
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() { So(Abs(tt.args.v), ShouldEqual, tt.want) })
		// t.Run(tt.name, func(t *testing.T) {
		// 	So(Abs(tt.args.v), ShouldEqual, tt.want)
		// 	// if got := Abs(tt.args.v); !reflect.DeepEqual(got, tt.want) {
		// 	// 	t.Errorf("Abs() = %v, want %v", got, tt.want)
		// 	// }
		// })
	}
}
