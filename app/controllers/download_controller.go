package controllers

import (
	httpUtils "teagans-web-api/app/utilities/http"
	"teagans-web-api/app/services/dtos"
	"teagans-web-api/app/services"
	"net/http"
	"strconv"
	"fmt"
	"io"
)

func DownloadVideo(w http.ResponseWriter, r *http.Request) {
	var data dtos.YoutubeDataDTO
	var dataErr error

	data.Url = r.URL.Query().Get("url")
	data.AudioOnly, dataErr = strconv.ParseBool(r.URL.Query().Get("audioOnly"))
	if dataErr != nil {
		httpUtils.RenderErrorJSON(w, r, dtos.CreateErrorDTO(dataErr, 400, false))
		return
	}

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.DownloaderService{BaseService: baseService}

	dto, errDTO := service.DownloadVideo(data)
	if errDTO.Exists() {
		httpUtils.RenderErrorJSON(w, r, errDTO)
		return
	}
	defer dto.Filereader.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", dto.Filename))
	w.Header().Set("Content-Type", dto.ContentType)

	io.Copy(w, dto.Filereader)
}
