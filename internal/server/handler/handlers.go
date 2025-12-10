package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/skdiver33/gophkeeper/internal/usermanager"
	"github.com/skdiver33/gophkeeper/model"
)

type UserManager interface {
	UserRegister(context.Context, *model.User) (*model.User, error)
	UserAuth(context.Context, *model.User) (string, error)
}

type ServerHandler struct {
	um UserManager
}

func (handler *ServerHandler) UserRegisterHandler(rw http.ResponseWriter, request *http.Request) {
	userData := model.User{}
	if err := json.NewDecoder(request.Body).Decode(&userData); err != nil {
		log.Printf("error user register. %s", err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	//userData.CryptPasswd()
	user, err := handler.um.UserRegister(request.Context(), &userData)
	if err != nil {
		returnStatus := http.StatusInternalServerError
		if errors.Is(err, usermanager.ErrUserAlreadyExist) {
			returnStatus = http.StatusConflict
		}
		log.Printf("error user register. %s", err.Error())
		http.Error(rw, err.Error(), returnStatus)
		return
	}
	userAuthToken, err := handler.um.UserAuth(request.Context(), user)
	if err != nil {
		returnStatus := http.StatusInternalServerError
		if errors.Is(err, usermanager.ErrUserWithCredNotFound) {
			returnStatus = http.StatusUnauthorized
		}
		log.Printf("error user register. auth error. %s", err.Error())
		http.Error(rw, err.Error(), returnStatus)
		return
	}

	cookie := http.Cookie{}
	cookie.Name = "jwt"
	cookie.Value = userAuthToken
	http.SetCookie(rw, &cookie)
	rw.Header().Set("Content-Type", "application/text-plain")
	rw.Write([]byte(userAuthToken))
}

func (handler *ServerHandler) UserLoginHandler(rw http.ResponseWriter, request *http.Request) {
	userData := model.User{}
	if err := json.NewDecoder(request.Body).Decode(&userData); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	//userData.CryptPasswd()
	userAuthToken, err := handler.um.UserAuth(request.Context(), &userData)
	if err != nil {
		returnStatus := http.StatusInternalServerError
		if errors.Is(err, usermanager.ErrUserWithCredNotFound) {
			returnStatus = http.StatusUnauthorized
		}
		log.Printf("error user login. auth error. %s", err.Error())
		http.Error(rw, err.Error(), returnStatus)
		return
	}
	cookie := http.Cookie{}
	cookie.Name = "jwt"
	cookie.Value = userAuthToken
	http.SetCookie(rw, &cookie)
	rw.Header().Set("Content-Type", "application/text-plain")
	rw.Write([]byte(userAuthToken))
}
