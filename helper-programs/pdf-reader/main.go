package main

import (
	"fmt"
	"log"

	"code.sajari.com/docconv"
)

func main() {
	res, err := docconv.ConvertPath("test2.pdf")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
