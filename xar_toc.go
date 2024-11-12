package xargon

import (
	"encoding/xml"
)

const (
	typeFile      = "file"
	typeDirectory = "directory"
)

type xarStyleAttr struct {
	Text  string `xml:",chardata"`
	Style string `xml:"style,attr"`
}

type xarOffsetSize struct {
	xarStyleAttr
	Offset string `xml:"offset"`
	Size   string `xml:"size"`
}

type xarSignature struct {
	xarOffsetSize
	KeyInfo struct {
		XmlNamespace string `xml:"xmlns,attr"`
		X509Data     struct {
			X509Certificate []string `xml:"X509Certificate"`
		} `xml:"X509Data"`
	} `xml:"KeyInfo"`
}

type xarToc struct {
	XMLName xml.Name `xml:"xar"`
	Toc     struct {
		CreationTime string        `xml:"creation-time"`
		Checksum     xarOffsetSize `xml:"checksum"`
		Signature    xarSignature  `xml:"signature"`
		XSignature   xarSignature  `xml:"x-signature"`
		File         []xarFile     `xml:"file"`
	} `xml:"toc"`
}
