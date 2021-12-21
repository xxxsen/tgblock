package stat

import (
	"syscall"
	"time"
)

func buildStat(fst *FileStat, isys interface{}) {
	fst.Uid = -1
	fst.Gid = -1
	sys := isys.(*syscall.Win32FileAttributeData)
	fst.CTime = time.Unix(0, stat.CreationTime.Nanoseconds()).Unix()
}
