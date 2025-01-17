package cpio

import (
	"bytes"
	"fmt"
	"io"
)
import newc "github.com/cavaliergopher/cpio"
import odc "github.com/korylprince/go-cpio-odc"

const magicFieldSize = 6

var odcMagic = []byte("070707")
var newcMagic = []byte("070701")
var newcCRCMagic = []byte("070702")

type File struct {
	name   string
	isDir  bool
	reader io.Reader
}

func (f *File) Name() string {
	return f.name
}

func (f *File) IsDir() bool {
	return f.isDir
}

func (f *File) Reader() io.Reader {
	return f.reader
}

type Reader struct {
	reader      *fakeReader
	newcReader  *newc.Reader
	odcReader   *odc.Reader
	currentFile *File
}

func NewReader(r io.Reader) *Reader {
	reader := &fakeReader{reader: r}
	return &Reader{
		reader:     reader,
		newcReader: newc.NewReader(reader),
		odcReader:  odc.NewReader(reader),
	}
}

func (r *Reader) Next() (*File, error) {
	magic, err := r.reader.readMagic()
	if err != nil {
		return nil, err
	}
	if bytes.Equal(magic, newcMagic) || bytes.Equal(magic, newcCRCMagic) {
		header, err := r.newcReader.Next()
		if err != nil {
			return nil, fmt.Errorf("could not read header: %w", err)
		}
		return &File{
			name:   header.Name,
			isDir:  header.Mode.IsDir(),
			reader: r.newcReader,
		}, nil
	} else if bytes.Equal(magic, odcMagic) {
		file, err := r.odcReader.Next()
		if err != nil {
			return nil, fmt.Errorf("could not read file: %w", err)
		}
		return &File{
			name:   file.Path,
			isDir:  file.IsDir(),
			reader: file,
		}, nil
	}
	return nil, io.EOF
}

func (r *Reader) Read(p []byte) (n int, err error) {
	if r.currentFile == nil {
		return 0, io.EOF
	}
	return r.currentFile.reader.Read(p)
}

type fakeReader struct {
	magicBuf []byte
	reader   io.Reader
}

func (r *fakeReader) readMagic() ([]byte, error) {
	r.magicBuf = make([]byte, magicFieldSize)
	if _, err := io.ReadFull(r.reader, r.magicBuf); err != nil {
		return nil, err
	}
	return r.magicBuf, nil
}

func (r *fakeReader) Read(p []byte) (n int, err error) {
	if r.reader == nil {
		return 0, io.EOF
	}
	if len(r.magicBuf) > 0 {
		count := 0
		for i := 0; i < len(r.magicBuf) && i < len(p); i++ {
			p[i] = r.magicBuf[i]
			count++
		}
		r.magicBuf = nil
		return count, nil
	}
	return r.reader.Read(p)
}
