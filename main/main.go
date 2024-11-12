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

	if err := xr.Verify(); err != nil {
		panic(err)
	}

	fmt.Println("no errors")

}
