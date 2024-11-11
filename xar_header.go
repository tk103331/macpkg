package xargon

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
