// package main

// import (
// 	"bingo/internal/app/server"
// 	"bingo/pkg/ps"
// 	"bingo/pkg/testing/gen"
// )

// func main() {
// 	gen.GenTestCase("face_api", "add_user", 10002)

// 	server.Serve()
// 	ps.Perform()
// }

package main

import (
	"fmt"
	"log"

	pdf "github.com/adrg/go-wkhtmltopdf"
)

func main() {
	pdf.Init()
	defer pdf.Destroy()

	// Create object from url
	object1, err := pdf.NewObject("https://www.baidu.com/")
	if err != nil {
		log.Fatal(err)
	}
	object1.SetOption("footer.right", "[page]")

	// Create converter
	converter := pdf.NewConverter()
	defer converter.Destroy()

	// Add created objects to the converter
	converter.AddObject(object1)

	// Add converter options
	converter.SetOption("documentTitle", "Sample document")
	converter.SetOption("margin.left", "10mm")
	converter.SetOption("margin.right", "10mm")
	converter.SetOption("margin.top", "10mm")
	converter.SetOption("margin.bottom", "10mm")

	// Convert the objects and get the output PDF document
	output, err := converter.Convert()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(output))
}
