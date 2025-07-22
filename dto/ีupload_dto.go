package dto

type FileMeta struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Size int64  `json:"size"`
}

type RequestGetFile struct {
	Folder string `json:"folder"`
	File   string `json:"file"`
}
