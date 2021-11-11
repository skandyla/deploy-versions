package internal

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
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
func (h VersionHandler) Info(w http.ResponseWriter, r *http.Request) {
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

func (h VersionHandler) GetAllVersions(w http.ResponseWriter, r *http.Request) {
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

func (h VersionHandler) GetVersionByID(w http.ResponseWriter, r *http.Request) {
	buildID, err := strconv.Atoi(chi.URLParam(r, "buildID"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resp, err := h.storage.GetVersionByID(buildID)
	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			//respondWithError(w, http.StatusNotFound, fmt.Sprintf("buildID:%v not found", buildID))
			respondWithError(w, http.StatusNotFound, "not found")

		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	respondWithJSON(w, 200, resp)
}

func (h VersionHandler) PostVersion(w http.ResponseWriter, r *http.Request) {
	var createVersionRequest models.CreateVersionRequest
	err := json.NewDecoder(r.Body).Decode(&createVersionRequest)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	validate := validator.New()
	err = validate.Struct(createVersionRequest)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	version := versionDBModelFromCreateRequest(createVersionRequest)

	//in current behaviour entities not required to be unique
	err = h.storage.PostVersion(version)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

//TBD
func (h VersionHandler) PutVersionByID(w http.ResponseWriter, req *http.Request) {
}

func (h VersionHandler) DeleteVersionByID(w http.ResponseWriter, req *http.Request) {
}

//sugar
func respondWithError(w http.ResponseWriter, code int, message string) {
	log.Println(message)
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
