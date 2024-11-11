package main

import (
	"fmt"
	"github.com/arelate/xargon"
	"os"
)

func main() {

	path := ""
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	xr, err := xargon.NewReader(path)
	if err != nil {
		panic(err)
	}
	defer xr.Close()

	fmt.Println(xr.TOC())

	fmt.Println("no errors")

}
