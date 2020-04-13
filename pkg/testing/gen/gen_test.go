package gen

import (
	"testing"
)

func TestGenTestCase(t *testing.T) {
	type args struct {
		pkg     string
		feature string
		id      int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"", args{pkg: "face_api", feature: "add_user", id: 1089}},
		{"", args{pkg: "face_api", feature: "add_user", id: 1090}},
		{"", args{pkg: "face_api", feature: "register", id: 1091}},
		{"", args{pkg: "face_api", feature: "add_user", id: 1092}},
		{"", args{pkg: "ocr", feature: "add_user", id: 2080}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TestCase(tt.args.pkg, tt.args.feature, tt.args.id)
		})
	}
}
