package healthz

import (
	"fmt"
	"net/http"

	"stats-of/internal/entities"
	"stats-of/internal/utils"
)

type response struct {
	Name         string `json:"name"`
	BuildVersion string `json:"build_version"`
	BuildTime    string `json:"build_time"`
	GitTag       string `json:"git_tag"`
	GitHash      string `json:"git_hash"`
}

func MakeHandler(info *entities.AppInfo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		response := &response{
			Name:         info.Name,
			BuildVersion: info.BuildVersion,
			BuildTime:    info.BuildTime,
			GitTag:       info.GitTag,
			GitHash:      info.GitHash,
		}
		err := utils.SuccessRespondWith200(w, response)
		if err != nil {
			fmt.Printf("failed to decode response = %v, error = %v\n", response, err)
		}
	}
}
