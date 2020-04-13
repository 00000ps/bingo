package main

import (
	"bingo/internal/app/server"
	"bingo/pkg/ps"
	"bingo/pkg/testing/gen"
)

func main() {
	gen.GenTestCase("face_api", "add_user", 10002)

	server.Serve()
	ps.Perform()
}
