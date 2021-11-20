package transport

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"github.com/skandyla/deploy-versions/models"
)

type Versions interface {
	Create(ctx context.Context, version models.CreateVersionRequest) error
	GetByID(ctx context.Context, id int) (models.VersionDBModel, error)
	GetAll(ctx context.Context) ([]models.VersionResponse, error)
	Delete(ctx context.Context, id int) error
	Health(ctx context.Context) error
	//Update(ctx context.Context, id int64, req models.UpdateVersionRequest) error
}

type Handler struct {
	versionsService Versions
}

func NewHandler(versions Versions) *Handler {
	return &Handler{
		versionsService: versions,
	}
}

func (h *Handler) InitRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(loggingMiddleware) //test our own middleware implementation

	r.Route("/info", func(r chi.Router) {
		r.Get("/", h.info)
	})

	r.Route("/versions", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Get("/", h.getAllVersions)
	})

	r.Route("/version", func(r chi.Router) {
		r.Post("/", h.createVersion)
		r.Route("/{buildID}", func(r chi.Router) {
			r.Get("/", h.getVersionByID)
			r.Put("/", h.updateVersionByID) //update entity
			r.Delete("/", h.deleteVersionByID)
		})
	})

	return r
}

//-------------------------------
//endpoints
func (h Handler) info(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := h.versionsService.Health(ctx)
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

func (h Handler) getAllVersions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	versions, err := h.versionsService.GetAll(ctx)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(versions)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) getVersionByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	buildID, err := strconv.Atoi(chi.URLParam(r, "buildID"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resp, err := h.versionsService.GetByID(ctx, buildID)
	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "not found")

		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	respondWithJSON(w, 200, resp)
}

func (h Handler) createVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	//in current behaviour entities not required to be unique
	err = h.versionsService.Create(ctx, createVersionRequest)
	if err != nil {
		log.Println("createVersion() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

//TBD
func (h Handler) updateVersionByID(w http.ResponseWriter, req *http.Request) {
}

func (h Handler) deleteVersionByID(w http.ResponseWriter, req *http.Request) {
}

//-------------------------------
func respondWithError(w http.ResponseWriter, code int, message string) {
	log.Println(message)
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(response)
	if err != nil {
		log.Println(err)
	}
}
