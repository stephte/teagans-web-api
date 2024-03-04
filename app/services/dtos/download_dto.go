package dtos

import (
	"io"
)

type DownloadDTO struct {
	Filename	string
	Filebytes	[]byte
	Filereader	io.ReadCloser
	Filesize	int64
	ContentType	string
}


type YoutubeDataDTO struct {
	Url			string
	AudioOnly	bool
}
