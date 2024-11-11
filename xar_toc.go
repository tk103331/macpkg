package xargon

import "encoding/xml"

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

type xarFile struct {
	Id   string `xml:"id,attr"`
	Data struct {
		Length            string       `xml:"length"`
		Offset            string       `xml:"offset"`
		Size              string       `xml:"size"`
		Encoding          xarStyleAttr `xml:"encoding"`
		ExtractedChecksum xarStyleAttr `xml:"extracted-checksum"`
		ArchivedChecksum  xarStyleAttr `xml:"archived-checksum"`
	} `xml:"data"`
	Type             string   `xml:"type"`
	Name             []string `xml:"name"`
	FinderCreateTime struct {
		NanoSeconds string `xml:"nanoseconds"`
		Time        string `xml:"time"`
	} `xml:"FinderCreateTime"`
	CTime    string    `xml:"ctime"`
	MTime    string    `xml:"mtime"`
	ATime    string    `xml:"atime"`
	Group    string    `xml:"group"`
	Gid      string    `xml:"gid"`
	User     string    `xml:"user"`
	Uid      string    `xml:"uid"`
	Mode     string    `xml:"mode"`
	DeviceNo string    `xml:"deviceno"`
	Inode    string    `xml:"inode"`
	File     []xarFile `xml:"file"`
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
