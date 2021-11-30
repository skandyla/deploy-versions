package transport

import (
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type CtxValue int

const (
	ctxUserID CtxValue = iota
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method": r.Method,
			"uri":    r.RequestURI,
		}).Info()
		//log.Printf("%s: [%s] - %s ", time.Now().Format(time.RFC3339), r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := h.sessionsStore.Get(r, "cookie-name")

		// Check if user is authenticated
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			//resp := map[string]interface{}{
			//	"code":  403,
			//	"error": "access denied",
			//}
			//respondWithJSON(w, http.StatusForbidden, resp)
			handleError401(w, "authMiddleware", "access denied", errors.New("cookie is not set"))

			return
		}
		next.ServeHTTP(w, r)
	})
}
