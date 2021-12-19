package download

type DownloadFileRequest struct {
	FileId string `schema:"file_id"`
}

type DownloadBlockRequest struct {
	FileId     string `schema:"file_id"`
	BlockIndex int64  `schema:"block_index"`
}
