package internal

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/skandyla/deploy-versions/models"
)

type VersionHandler struct {
	storage VersionStorage
}

func NewVersionHandler(storage VersionStorage) VersionHandler {
	return VersionHandler{
		storage: storage,
	}
}

// endpoints
func (h VersionHandler) Info(w http.ResponseWriter, req *http.Request) {
	err := h.storage.Health()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	type healthResponse struct {
		Status string `json:"status"`
	}
	resp := &healthResponse{
		Status: "ok",
	}
	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h VersionHandler) GetAllVersions(w http.ResponseWriter, req *http.Request) {
	versions, err := h.storage.GetAllVersions()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var versionResponse []models.VersionResponse
	for _, v := range versions {
		versionResponse = append(versionResponse, versionResponseFromDBModel(v))
	}

	err = json.NewEncoder(w).Encode(versionResponse)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h VersionHandler) GetVersion(w http.ResponseWriter, req *http.Request) {
}

func (h VersionHandler) PostVersion(w http.ResponseWriter, req *http.Request) {
}

func (h VersionHandler) GetVersionByID(w http.ResponseWriter, req *http.Request) {
}

func (h VersionHandler) PutVersionByID(w http.ResponseWriter, req *http.Request) {
}

func (h VersionHandler) DeleteVersionByID(w http.ResponseWriter, req *http.Request) {
}
