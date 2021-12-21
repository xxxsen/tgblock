package stat

import (
	"os"
)

type FileStat struct {
	IsDir bool
	MTime int64
	Mode  int64
	Name  string
	Size  int64
	Gid   int64
	Uid   int64
	CTime int64
}

func Stat(file string) (*FileStat, error) {
	st, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	fst := &FileStat{
		IsDir: st.IsDir(),
		MTime: st.ModTime().Unix(),
		Mode:  int64(st.Mode()),
		Size:  st.Size(),
		Name:  st.Name(),
	}
	buildStat(fst, st.Sys())
	return fst, nil
}
