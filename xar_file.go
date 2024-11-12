package xargon

import "path"

type xarFile struct {
	Id   string `xml:"id,attr"`
	Data struct {
		Length            int64        `xml:"length"`
		Offset            int64        `xml:"offset"`
		Size              int64        `xml:"size"`
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

func (xf *xarFile) FileName() string {
	if len(xf.Name) > 0 {
		return xf.Name[0]
	}
	return ""
}

func (xf *xarFile) indexFiles() ([]string, map[string]*xarFile) {

	order := make([]string, 0, len(xf.File))
	index := make(map[string]*xarFile)

	xfn := xf.FileName()

	switch xf.Type {
	case typeFile:
		order = append(order, xfn)
		index[xfn] = xf
	case typeDirectory:
		for _, df := range xf.File {
			do, di := df.indexFiles()
			for _, dof := range do {
				dofn := path.Join(xfn, dof)
				order = append(order, dofn)
				index[dofn] = di[dof]
			}
		}
	default:
		panic("unknown xar file type: " + xf.Type)
	}

	return order, index
}
