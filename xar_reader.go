package xargon

import (
	"bytes"
	"compress/zlib"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"io"
	"os"
)

const (
	applicationOctetStreamMimeType = "application/octet-stream"
	applicationGzipStreamMimeType  = "application/x-gzip"
)

type XarReader struct {
	file       *os.File
	header     *xarHeader
	toc        *xarToc
	heap       *io.SectionReader
	certs      []*x509.Certificate
	filesIndex map[string]*xarFile
	filesOrder []string
}

func (xr *XarReader) readHeader() error {

	if xr.file == nil {
		return errors.New("cannot read header from nil file")
	}

	h := make([]byte, xarHeaderSize)
	if n, err := xr.file.ReadAt(h, 0); err != nil {
		return err
	} else if n != xarHeaderSize {
		return errors.New("xar header size mismatch")
	}

	xh := &xarHeader{
		magic:                 binary.BigEndian.Uint32(h[0:4]),
		size:                  binary.BigEndian.Uint16(h[4:6]),
		version:               binary.BigEndian.Uint16(h[6:8]),
		tocLengthCompressed:   binary.BigEndian.Uint64(h[8:16]),
		tocLengthUncompressed: binary.BigEndian.Uint64(h[16:24]),
		checksumAlgorithm:     binary.BigEndian.Uint32(h[24:28]),
	}

	// validate expected values

	if xh.magic != xarHeaderMagic {
		return errors.New("unexpected xar format")
	}

	if xh.size != xarHeaderSize {
		return errors.New("unexpected xar h size")
	}

	if xh.version != xarHeaderVersion {
		return errors.New("unsupported xar version")
	}

	if xh.tocLengthCompressed == 0 {
		return errors.New("unexpected xar toc compressed length")
	}

	if xh.tocLengthUncompressed == 0 {
		return errors.New("unexpected xar toc uncompressed length")
	}

	if xh.checksumAlgorithm == 3 {
		return errors.New("unsupported xar checksum algorithm")
	}

	xr.header = xh

	return nil
}

func (xr *XarReader) readToc() error {

	if xr.file == nil {
		return errors.New("cannot read toc from nil file")
	}

	if xr.header == nil {
		return errors.New("cannot read toc from nil header")
	}

	toc := make([]byte, xr.header.tocLengthCompressed)
	if n, err := xr.file.ReadAt(toc, xarHeaderSize); err != nil {
		return err
	} else if uint64(n) != xr.header.tocLengthCompressed {
		return errors.New("xar toc size mismatch")
	}

	br := bytes.NewReader(toc)
	zr, err := zlib.NewReader(br)
	if err != nil {
		return err
	}

	return xml.NewDecoder(zr).Decode(&xr.toc)
}

func (xr *XarReader) heapOffset() (int64, error) {
	if xr.header == nil {
		return -1, errors.New("cannot read toc from nil header")
	}

	return xarHeaderSize + int64(xr.header.tocLengthCompressed), nil
}

func (xr *XarReader) readHeap() error {

	if xr.file == nil {
		return errors.New("cannot read heap from nil file")
	}

	if xr.header == nil {
		return errors.New("cannot read heap from nil header")
	}

	heapOffset, err := xr.heapOffset()
	if err != nil {
		return nil
	}

	xr.heap = io.NewSectionReader(xr.file, heapOffset, int64(xr.header.size)-heapOffset)

	return nil
}

func (xr *XarReader) readCertificates() error {

	if xr.toc == nil {
		return errors.New("cannot check signatures from nil toc")
	}

	for _, encCert := range xr.toc.Toc.Signature.KeyInfo.X509Data.X509Certificate {

		decCert, err := base64.StdEncoding.DecodeString(encCert)
		if err != nil {
			return err
		}

		if cert, err := x509.ParseCertificate(decCert); err != nil {
			return err
		} else {
			xr.certs = append(xr.certs, cert)
		}
	}

	return nil
}

func (xr *XarReader) Verify() error {
	if len(xr.certs) == 0 {
		if err := xr.readCertificates(); err != nil {
			return err
		}
	}

	//TODO: implement cert verification

	return nil
}

func (xr *XarReader) indexFiles() error {

	if xr.toc == nil {
		return errors.New("cannot index files from nil toc")
	}

	files := xr.toc.Toc.File

	xr.filesOrder = make([]string, 0, len(files))
	xr.filesIndex = make(map[string]*xarFile, len(files))

	for _, file := range files {
		fo, fi := file.indexFiles()
		for _, fof := range fo {
			xr.filesOrder = append(xr.filesOrder, fof)
			xr.filesIndex[fof] = fi[fof]
		}
	}

	return nil
}

func (xr *XarReader) ReadFiles() ([]string, error) {

	if xr.filesOrder == nil {
		if err := xr.indexFiles(); err != nil {
			return nil, err
		}
	}

	if xr.filesOrder == nil {
		return nil, errors.New("xar files have not been ordered")
	}

	return xr.filesOrder, nil
}

func (xr *XarReader) Open(name string) (io.Reader, error) {

	if xr.filesIndex == nil {
		if err := xr.indexFiles(); err != nil {
			return nil, err
		}
	}

	if xr.filesIndex == nil {
		return nil, errors.New("xar files have not been indexed")
	}

	if xf, ok := xr.filesIndex[name]; ok {
		return xr.openFile(xf)
	} else {
		return nil, errors.New("xar file not found: " + name)
	}
}

func (xr *XarReader) openFile(xf *xarFile) (io.Reader, error) {

	if xf == nil {
		return nil, errors.New("cannot open nil xar file")
	}

	sectionReader := io.NewSectionReader(xr.heap, xf.Data.Offset, xf.Data.Length)

	enc := xf.Data.Encoding.Style
	switch enc {
	case applicationOctetStreamMimeType:
		return sectionReader, nil
	case applicationGzipStreamMimeType:

		return zlib.NewReader(sectionReader)
	default:
		return nil, errors.New("unknown xar file encoding: " + enc)
	}
}

func (xr *XarReader) Close() error {
	return xr.file.Close()
}

func NewReader(path string) (*XarReader, error) {

	xf, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	xr := &XarReader{
		file: xf,
	}

	// xar file = Header + TOC + Heap
	// we initialize all three below

	if err := xr.readHeader(); err != nil {
		return nil, err
	}

	if err := xr.readToc(); err != nil {
		return nil, err
	}

	if err := xr.readHeap(); err != nil {
		return nil, err
	}

	return xr, nil
}
