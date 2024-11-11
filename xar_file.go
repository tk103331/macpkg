package xargon

import (
	"encoding/binary"
	"errors"
	"os"
)

const (
	xarHeaderMagic   = 0x78617221 /* 'xar!' */
	xarHeaderVersion = 1          // Currently there is only version 1.
	xarHeaderSize    = 28         /* (32 + 16 + 16 + 64 + 64 + 32) / 8 */
)

type xarHeader struct {
	/* This should always equal 'xar!' */
	magic                 uint32 // File signature used to identify the file format as Xar.
	size                  uint16 // Header size
	version               uint16 // Version of Xar format to use.
	tocLengthCompressed   uint64 // Length of the TOC compressed data.
	tocLengthUncompressed uint64 // Length of the TOC uncompressed data.
	/* Checksum algorithm:
	0 = none
	1 = SHA1
	2 = MD5
	3 = SHA-256
	4 = SHA-512 */
	checksumAlgorithm uint32
	/* A nul-terminated, zero-padded to multiple of 4, message digest name
	 * appears here if checksumAlgorithm is 3 which must not be empty ("") or "none".
	 */
}

type XarReader struct {
	file   *os.File
	header *xarHeader
}

func (xr *XarReader) ReadHeader() error {

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

func (xr *XarReader) ReadTOC() error {

	if xr.file == nil {
		return errors.New("cannot read toc from nil file")
	}

	if xr.header == nil {
		return errors.New("cannot read toc from nil header")
	}

	return nil
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

	return xr, nil
}
