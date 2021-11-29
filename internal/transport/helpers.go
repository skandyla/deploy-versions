package transport

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

//-------------------------------
func handleError500(w http.ResponseWriter, logMarker string, err error) {
	log.Errorf("%s: %+v\n", logMarker, err)
	w.WriteHeader(http.StatusInternalServerError)
}

func handleError400(w http.ResponseWriter, logMarker string, respMsg string, err error) {
	log.Debugf("%s: %+v\n", logMarker, err)
	respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": respMsg})
}

func handleError401(w http.ResponseWriter, logMarker string, respMsg string, err error) {
	log.Debugf("%s: %+v\n", logMarker, err)
	respondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": respMsg})
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
