package services

import (
	"teagans-web-api/app/services/dtos"
	"teagans-web-api/app/utilities"
	"github.com/kkdai/youtube/v2"
	"net/url"
	"strings"
	"errors"
	"fmt"
	// "os"
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
	if strings.Contains(u.Hostname(), "youtu.be") {
		video_id = strings.Replace(u.Path, "/", "", 1)
	} else {
		video_id = qmap.Get("v")
	}

	if video_id == "" {
		return dtos.DownloadDTO{}, dtos.CreateErrorDTO(errors.New("Video not found"), 0, false)
	}

	// now get the video data and the stream
	client := youtube.Client{}
	fmt.Println("client")
	fmt.Println(client)
	video, verr := client.GetVideo(video_id)
	if verr != nil {
		fmt.Println("verr")
		fmt.Println(verr)
		return dtos.DownloadDTO{}, dtos.CreateErrorDTO(verr, 0, false)
	}
	// fmt.Println("video")
	// fmt.Println(video)

	formats := video.Formats.WithAudioChannels()
	format := formats[0]
	audioStr := ""
	if data.AudioOnly {
		formats = formats.Type("audio/mp4")
		format = formats[0]
		audioStr = "_audio"
	}

	fmt.Println("format.URL")
	fmt.Println(format.URL)

	// fmt.Println("format")
	// fmt.Println(format)

	// fmt.Println("Formats:")
	// for i, format := range formats {
	// 	println(i, format.MimeType, format.Quality)
	// }

	var rv dtos.DownloadDTO
	rv.Filename = fmt.Sprintf("%s%s.mp4", utilities.ConvertToFilename(video.Title), audioStr)
	rv.ContentType = format.MimeType
	rv.Filereader, rv.Filesize, err = client.GetStream(video, &format)
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
		return dtos.DownloadDTO{}, dtos.CreateErrorDTO(err, 0, false)
	}

	// var byts []byte
	// rv.Filereader.Read(byts)

	// fmt.Println("rv")
	// fmt.Println(rv)

	return rv, dtos.ErrorDTO{}
}
