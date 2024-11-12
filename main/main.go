package main

import (
	"fmt"
	"github.com/arelate/xargon"
	"io"
	"os"
	"path/filepath"
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

	// testing extracting files to temp folder
	rootDir := filepath.Join(os.TempDir(), "xargon")
	fmt.Println("file://" + rootDir)

	if tocFiles, err := xr.ReadFiles(); err != nil {
		panic(err)
	} else {
		for _, tf := range tocFiles {
			fullPath := filepath.Join(rootDir, tf)
			dir, _ := filepath.Split(fullPath)
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				if err := os.MkdirAll(dir, 0755); err != nil {
					panic(err)
				}
			}

			if src, err := xr.Open(tf); err != nil {
				panic(err)
			} else {
				dst, err := os.Create(fullPath)
				if err != nil {
					panic(err)
				}

				if _, err := io.Copy(dst, src); err != nil {
					panic(err)
				}
			}
		}
	}

	fmt.Println("all done")

}
