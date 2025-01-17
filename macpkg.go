package macpkg

import (
	"encoding/xml"
	"fmt"
	"github.com/tk103331/macpkg/cpio"
	"github.com/tk103331/macpkg/xar"
	"io"
	"os"
	"path/filepath"
	"slices"
)

type MacPkg struct {
}

type PkgReader struct {
	xarReader *xar.XarReader
	xarFiles  []string

	pkgIds []string
}

type Distribution struct {
	PkgRef []struct {
		ID string
	}
}

func Open(path string) (*PkgReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return NewReader(file)
}

func NewReader(file *os.File) (*PkgReader, error) {
	xarReader, err := xar.NewReader(file)
	if err != nil {
		return nil, err
	}
	files, err := xarReader.ReadFiles()
	if err != nil {
		return nil, err
	}
	reader := &PkgReader{
		xarReader: xarReader,
		xarFiles:  files,
	}
	return reader, nil
}

func (r *PkgReader) ReadPkgIDs() ([]string, error) {
	if len(r.pkgIds) == 0 {
		ids, err := r.readPkgIDs()
		if err != nil {
			return nil, err
		}
		r.pkgIds = ids
	}
	return r.pkgIds, nil
}

func (r *PkgReader) readPkgIDs() ([]string, error) {
	reader, err := r.Open("Distribution")
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	dist := &Distribution{}
	err = xml.Unmarshal(data, dist)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(dist.PkgRef))
	for i := 0; i < len(dist.PkgRef); i++ {
		ids[i] = dist.PkgRef[i].ID
	}
	return ids, nil
}

func (r *PkgReader) ReadFiles() ([]string, error) {
	if len(r.xarFiles) == 0 {
		files, err := r.xarReader.ReadFiles()
		if err != nil {
			return nil, err
		}
		r.xarFiles = files
	}
	return r.xarFiles, nil
}

func (r *PkgReader) readFiles() ([]string, error) {
	files, err := r.xarReader.ReadFiles()
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (r *PkgReader) Open(file string) (io.Reader, error) {
	return r.xarReader.Open(file)
}

func (r *PkgReader) OpenPayload(pkgId string) (*cpio.Reader, error) {
	if len(r.pkgIds) > 0 && !slices.Contains(r.pkgIds, pkgId) {
		return nil, fmt.Errorf("invalid package id: %s", pkgId)
	}
	payload := filepath.Join(pkgId, "Payload")
	reader, err := r.xarReader.Open(payload)
	if err != nil {
		return nil, err
	}
	cpioReader := cpio.NewReader(reader)
	return cpioReader, nil
}
