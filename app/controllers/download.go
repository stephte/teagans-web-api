package controllers

import (
	"youtube-downloader/app/utilities/http_utils"
	"youtube-downloader/app/services/dtos"
	"youtube-downloader/app/services"
	"encoding/json"
	"net/http"
	"fmt"
	"io"
)

func DownloadVideo(w http.ResponseWriter, r *http.Request) {
	var data dtos.YoutubeDataDTO
	bindErr := json.NewDecoder(r.Body).Decode(&data)
	if bindErr != nil {
		http_utils.RenderErrorJSON(w, r, dtos.CreateErrorDTO(bindErr, 400, false))
		return
	}

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.DownloaderService{BaseService: baseService}

	dto, errDTO := service.DownloadVideo(data)
	if errDTO.Exists() {
		http_utils.RenderErrorJSON(w, r, errDTO)
		return
	}
	defer dto.Filereader.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", dto.Filename))
	w.Header().Set("Content-Type", dto.ContentType)

	io.Copy(w, dto.Filereader)
}
