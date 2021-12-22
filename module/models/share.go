package models

type CreateShareRequest struct {
	FileId     string `json:"file_id"`
	ExpireTime int64  `json:"expire_time"`
	Key        string `json:"key"`
}

type CreateShareResponse struct {
	URL string `json:"url"`
}

type GetShareRequest struct {
	Key  string `json:"key"`
	Code string `json:"code"`
}
