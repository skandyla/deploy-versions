package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/skandyla/deploy-versions/internal/domain"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		handleError400(w, "signUp:ReadAll", "Read body failed", err)
		return
	}

	var inp domain.SignUpInput
	if err = json.Unmarshal(reqBytes, &inp); err != nil {
		handleError400(w, "signUp:Unmarshal", "Json decoding failed", err)
		return
	}

	if err := inp.Validate(); err != nil {
		handleError400(w, "signUp:Validate", "Json validation failed", err)
		return
	}

	err = h.usersService.SignUp(r.Context(), inp)
	if err != nil {
		handleError500(w, "signUp:service", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError400(w, "signIn:ReadAll", "Read body failed", err)
		return
	}

	var inp domain.SignInInput
	if err = json.Unmarshal(reqBytes, &inp); err != nil {
		handleError400(w, "signIn:Unmarshal", "Json decoding failed", err)
		return
	}

	if err := inp.Validate(); err != nil {
		handleError400(w, "signIn:Validate", "Json validation failed", err)
		return
	}

	accessToken, refreshToken, err := h.usersService.SignIn(r.Context(), inp)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			handleError400(w, "signIn:Service", fmt.Sprintf("%+v", err), err)
			return
		}

		handleError500(w, "signIn:Service", err)
		return
	}

	response, err := json.Marshal(map[string]string{
		"token": accessToken,
	})
	if err != nil {
		handleError500(w, "signIn:response:Marshal", err)
		return
	}

	w.Header().Add("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))
	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh-token")
	if err != nil {
		handleError400(w, "refresh", "error get cookkie", err)
		return
	}

	log.Infof("%s", cookie.Value)

	accessToken, refreshToken, err := h.usersService.RefreshTokens(r.Context(), cookie.Value)
	if err != nil {
		handleError500(w, "signIn:RefreshTokens", err)
		return
	}

	response, err := json.Marshal(map[string]string{
		"token": accessToken,
	})
	if err != nil {
		handleError500(w, "signIn:Marshal:response", err)
		return
	}

	w.Header().Add("Set-Cookie", fmt.Sprintf("refresh-token='%s'; HttpOnly", refreshToken))
	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}
