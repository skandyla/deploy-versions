package transport

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	log "github.com/sirupsen/logrus"
	"github.com/skandyla/deploy-versions/models"
)

//-------------------------------
//endpoints
func (h Handler) info(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := h.versionsService.Health(ctx)
	if err != nil {
		handleError500(w, "Health", err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h Handler) getAllVersions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	versions, err := h.versionsService.GetAll(ctx)
	if err != nil {
		handleError500(w, "getAllVersions", err)
		return
	}

	err = json.NewEncoder(w).Encode(versions)
	if err != nil {
		handleError500(w, "getAllVersions", err)
		return
	}
}

func (h Handler) getVersionByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	buildID, err := strconv.Atoi(chi.URLParam(r, "buildID"))
	if err != nil {
		handleError400(w, "getVersionByID", "can't parse buildID", err)
		return
	}
	resp, err := h.versionsService.GetByID(ctx, buildID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			handleError400(w, "getVersionByID", "buildID not found", err)
			return
		}
		handleError500(w, "getVersionByID", err)
		return
	}

	//for test
	log.Debugf("ctxUserID:%+v", ctx.Value(ctxUserID))

	respondWithJSON(w, 200, resp)
}

func (h Handler) createVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var createVersionRequest models.CreateVersionRequest

	err := json.NewDecoder(r.Body).Decode(&createVersionRequest)
	if err != nil {
		handleError400(w, "createVersion", "Json decoding failed", err)
		return
	}

	validate := validator.New()
	err = validate.Struct(createVersionRequest)
	if err != nil {
		handleError400(w, "createVersion", "Json validation failed", err)
		return
	}

	//in current behaviour entities not required to be unique
	err = h.versionsService.Create(ctx, createVersionRequest)
	if err != nil {
		handleError500(w, "createVersion", err)
		return
	}

	//w.WriteHeader(http.StatusCreated)
	resp := map[string]interface{}{
		"code":      200,
		"createdId": createVersionRequest.BuildID,
	}
	respondWithJSON(w, http.StatusOK, resp)
}

//TBD
func (h Handler) updateVersionByID(w http.ResponseWriter, req *http.Request) {
}

func (h Handler) deleteVersionByID(w http.ResponseWriter, req *http.Request) {
}
