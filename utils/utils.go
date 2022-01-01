package utils

import (
	"fmt"
	"io"

	"github.com/xxxsen/tgblock/protos/gen/tgblock"
)

func CalcBlockSizeByFileId(meta *tgblock.FileContext, fileid string) (int64, error) {
	var index int64 = -1
	for i := 0; i < int(meta.BlockCount); i++ {
		if meta.FileList[i].GetFileId() == fileid {
			index = int64(i)
			break
		}
	}
	return CalcBlockSizeByIndex(meta, index)
}

func CalcBlockSizeByIndex(meta *tgblock.FileContext, index int64) (int64, error) {
	if index < 0 {
		return 0, fmt.Errorf("invalid index:%d", index)
	}
	if meta.BlockCount <= 1 {
		return meta.FileSize, nil
	}
	if index+1 > meta.BlockCount {
		return 0, fmt.Errorf("block index out of range, index:%d", index)
	}
	if index+1 < meta.BlockCount {
		return meta.BlockSize, nil
	}
	return meta.FileSize - meta.FileSize/meta.BlockSize*meta.BlockSize, nil
}

func CalcSeek(filesize int64, current, offset int64, whence int) (int64, error) {
	var startAt int64
	if whence == io.SeekStart {
		if offset > filesize {
			return 0, fmt.Errorf("SeekStart: offset:%d > filesize:%d", offset, filesize)
		}
		if offset < 0 {
			return 0, fmt.Errorf("invalid offset:%d using SeekStart", offset)
		}
		startAt = offset
	} else if whence == io.SeekEnd {
		if offset > 0 {
			return 0, fmt.Errorf("invalid offset:%d using SeekEnd", offset)
		}
		if -1*offset+filesize < 0 {
			return 0, fmt.Errorf("SeekEnd: offset:%d+filesize:%d<0", offset, filesize)
		}
		startAt = filesize + offset
	} else if whence == io.SeekCurrent {
		startAt = current + offset
		if startAt > filesize || startAt < 0 {
			return 0, fmt.Errorf("SeekCurrent: seek out of range, current:%d, offset:%d, filesize:%d", current, offset, filesize)
		}
	} else {
		return 0, fmt.Errorf("invalid seek:%d, offset:%d", whence, offset)
	}
	return startAt, nil
}

//ReadIndexToBlockIndexOffset 将文件的读取位置转换成块id和块内偏移
func ReadIndexToBlockIndexOffset(meta *tgblock.FileContext, readindex int64) (int64, int64, error) {
	if readindex > meta.GetFileSize() {
		return 0, 0, fmt.Errorf("index out of range, index:%d, filesize:%d", readindex, meta.GetBlockSize())
	}
	blockid := readindex / meta.BlockSize
	blockindex := readindex % meta.BlockSize
	return blockid, blockindex, nil
}
