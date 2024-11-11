package xargon

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"os"
)

type XarReader struct {
	file   *os.File
	header *xarHeader
	toc    *xarToc
}

func (xr *XarReader) readHeader() error {

	if xr.file == nil {
		return errors.New("cannot read header from nil file")
	}

	// read header

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

func (xr *XarReader) readTOC() error {

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

// TODO: Remove this after testing
func (xr *XarReader) TOC() *xarToc {
	return xr.toc
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

	if err := xr.readHeader(); err != nil {
		return nil, err
	}

	if err := xr.readTOC(); err != nil {
		return nil, err
	}

	return xr, nil
}
