package stat

import "syscall"

func buildStat(fst *FileStat, isys interface{}) {
	sys := isys.(*syscall.Stat_t)
	fst.CTime = sys.Ctim.Sec
	fst.Uid = int64(sys.Uid)
	fst.Gid = int64(sys.Gid)
}
