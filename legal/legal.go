package legal

import (
	"code.sajari.com/docconv/client"
	"log"
	"rsc.io/pdf"
	// "os"
)

func testpdf() {
	f, err := pdf.Open("test.pdf")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(f.NumPage())

	text := f.Page(1).Content().Text

	for _, lines := range text {
		log.Println(lines.S)
	}
}

func testdocconv() {
	c := client.New()
	res, err := client.ConvertPath(c, "test.pdf")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res)
}
