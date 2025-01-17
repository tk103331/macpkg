package main

import (
	"fmt"
	"github.com/tk103331/macpkg/xar"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	path := ""
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	xr, err := xar.NewReader(path)
	if err != nil {
		panic(err)
	}
	defer xr.Close()

	if err := xr.CheckCertificatesSignatures(); err != nil {
		panic(err)
	}

	// testing extracting files to temp folder
	rootDir := filepath.Join(os.TempDir(), "xargon", strings.TrimSuffix(path, filepath.Ext(path)))
	fmt.Println("file://" + rootDir)

	if tocFiles, err := xr.ReadFiles(); err != nil {
		panic(err)
	} else {
		for _, tf := range tocFiles {

			checksum, method, err := xr.ChecksumMethod(tf)
			if err != nil {
				panic(err)
			}

			fmt.Println(tf, checksum, method)

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
