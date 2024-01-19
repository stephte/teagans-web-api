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
	unescapedUrl, escErr := url.QueryUnescape(data.Url)
	if escErr != nil {
		return dtos.DownloadDTO{}, dtos.CreateErrorDTO(escErr, 0, false)
	}

	u, err := url.Parse(unescapedUrl)
	if err != nil {
		return dtos.DownloadDTO{}, dtos.CreateErrorDTO(err, 0, false)
	}

	qmap := u.Query()
	var video_id string
	if strings.Contains(u.Hostname(), "youtube.com") {
		video_id = qmap.Get("v")
	} else if strings.Contains(u.Hostname(), "youtu.be") {
		video_id = strings.Replace(u.Path, "/", "", 1)
	} else {
		return dtos.DownloadDTO{}, dtos.CreateErrorDTO(errors.New("must be a youtube url"), 0, false)
	}

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
