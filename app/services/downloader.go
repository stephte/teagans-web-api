package services

import (
	"youtube-downloader/app/services/dtos"
	"youtube-downloader/app/utilities"
	"github.com/kkdai/youtube/v2"
	"net/url"
	"strings"
	"errors"
	"fmt"
)

type DownloaderService struct {
	*BaseService
}


func(this DownloaderService) DownloadVideo(data dtos.YoutubeDataDTO) (dtos.DownloadDTO, dtos.ErrorDTO) {
	// get id from URL:
	u, err := url.Parse(data.Url)
	if err != nil {
		return dtos.DownloadDTO{}, dtos.CreateErrorDTO(err, 0, false)
	}

	if !strings.Contains(u.Hostname(), "youtube.com") {
		return dtos.DownloadDTO{}, dtos.CreateErrorDTO(errors.New("must be a youtube url"), 0, false)
	}

	qmap := u.Query()
	video_id := qmap.Get("v")
	if video_id == "" {
		return dtos.DownloadDTO{}, dtos.CreateErrorDTO(errors.New("Video ID not found"), 0, false)
	}

	// now get the video data and the stream
	client := youtube.Client{}
	video, verr := client.GetVideo(video_id)
	if verr != nil {
		return dtos.DownloadDTO{}, dtos.CreateErrorDTO(verr, 0, false)
	}

	formats := video.Formats.WithAudioChannels()
	format := formats[0]
	audioStr := ""
	if data.AudioOnly {
		formats = formats.Type("audio/mp4")
		format = formats[0]
		audioStr = "_audio"
	}

	// fmt.Println("Formats:")
	// for i, format := range formats {
	// 	println(i, format.MimeType, format.Quality)
	// }

	var rv dtos.DownloadDTO
	rv.Filename = fmt.Sprintf("%s%s.mp4", utilities.ConvertToFilename(video.Title), audioStr)
	rv.ContentType = format.MimeType
	rv.Filereader, rv.Filesize, err = client.GetStream(video, &format)
	if err != nil {
		return dtos.DownloadDTO{}, dtos.CreateErrorDTO(err, 0, false)
	}

	return rv, dtos.ErrorDTO{}
}
