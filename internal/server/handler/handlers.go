package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/skdiver33/gophkeeper/internal/auth"
	"github.com/skdiver33/gophkeeper/internal/datamanager"
	"github.com/skdiver33/gophkeeper/internal/usermanager"
	"github.com/skdiver33/gophkeeper/model"
	"github.com/skdiver33/gophkeeper/protocol"
)

type UserManager interface {
	UserRegister(context.Context, *model.User) (*model.User, error)
	UserAuth(context.Context, *model.User) (string, error)
}

type DataManager interface {
	LoadData(ctx context.Context, data *protocol.ProtocolPackage, userID int) error
	GetData(ctx context.Context, md model.Metadata, userID int) (*protocol.ProtocolPackage, error)
	GetAllData(ctx context.Context, userID int) (*[]model.Metadata, error)
	DeleteData(ctx context.Context, md model.Metadata, userID int) error
}

type ServerHandler struct {
	um UserManager
	dm DataManager
}

func NewServerHandler(userManager UserManager, dataManager DataManager) *ServerHandler {
	return &ServerHandler{um: userManager, dm: dataManager}
}

func (handler *ServerHandler) UserRegisterHandler(rw http.ResponseWriter, request *http.Request) {
	userData := model.User{}
	if err := json.NewDecoder(request.Body).Decode(&userData); err != nil {
		slog.Error("user register", "error", err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	userData.Password = auth.GetPasswdHash(userData.Password)
	user, err := handler.um.UserRegister(request.Context(), &userData)
	if err != nil {
		returnStatus := http.StatusInternalServerError
		if errors.Is(err, usermanager.ErrUserAlreadyExist) {
			returnStatus = http.StatusConflict
		}
		slog.Error("user register", "error", err.Error())
		http.Error(rw, err.Error(), returnStatus)
		return
	}
	userAuthToken, err := handler.um.UserAuth(request.Context(), user)
	if err != nil {
		returnStatus := http.StatusInternalServerError
		if errors.Is(err, usermanager.ErrUserWithCredNotFound) {
			returnStatus = http.StatusUnauthorized
		}
		slog.Error("user auth", "error", err.Error())
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
	userData.Password = auth.GetPasswdHash(userData.Password)
	userAuthToken, err := handler.um.UserAuth(request.Context(), &userData)
	if err != nil {
		returnStatus := http.StatusInternalServerError
		if errors.Is(err, usermanager.ErrUserWithCredNotFound) {
			returnStatus = http.StatusUnauthorized
		}
		slog.Error("user login", "error", err.Error())
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

func (handler *ServerHandler) LoadDataHandler(rw http.ResponseWriter, request *http.Request) {
	pkg := protocol.ProtocolPackage{}
	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		slog.Error("read loading data", "error", err.Error())
		return
	}
	err = json.Unmarshal(body, &pkg)
	if err != nil {
		slog.Error("unmarshal receive package", "error", err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := auth.GetUserIDFromClaims(request.Context())
	if err != nil {
		slog.Error("get user ID from JWT", "error", err.Error())
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}
	err = handler.dm.LoadData(request.Context(), &pkg, userID)
	if err != nil {
		returnStatus := http.StatusInternalServerError
		if errors.Is(err, datamanager.ErrDataAlreadyLoad) {
			returnStatus = http.StatusConflict
		}
		slog.Error("save data in storage", "error", err.Error())

		http.Error(rw, err.Error(), returnStatus)
		return
	}
	rw.WriteHeader(http.StatusOK)

}

func (handler *ServerHandler) GetDataHandler(rw http.ResponseWriter, request *http.Request) {
	userID, err := auth.GetUserIDFromClaims(request.Context())
	if err != nil {
		slog.Error("error get user ID from JWT", "error", err.Error())
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}
	md := model.Metadata{}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		slog.Error("read metadata body getdatahandler", "error", err.Error())
		return
	}
	err = json.Unmarshal(body, &md)
	if err != nil {
		slog.Error("unmarshal receive metadata", "error", err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	pkg, err := handler.dm.GetData(request.Context(), md, userID)
	if err != nil {
		slog.Error("get data from storage", "error", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(pkg)
	if err != nil {
		slog.Error("marshal response", "error", err.Error())
		http.Error(rw, "", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(resp)

}

func (handler *ServerHandler) GetAllDataHandler(rw http.ResponseWriter, request *http.Request) {
	userID, err := auth.GetUserIDFromClaims(request.Context())
	if err != nil {
		slog.Error("error get user ID from JWT", "error", err.Error())
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}
	data, err := handler.dm.GetAllData(request.Context(), userID)
	if err != nil {
		slog.Error("error get all data for user", "error", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(*data) == 0 {
		http.Error(rw, "no data for user", http.StatusNoContent)
		return
	}
	resp, err := json.Marshal(data)
	if err != nil {
		slog.Error("marshall response data", "error", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(resp)
}

func (handler *ServerHandler) DeleteDataHandler(rw http.ResponseWriter, request *http.Request) {
	userID, err := auth.GetUserIDFromClaims(request.Context())
	if err != nil {
		slog.Error("error get user ID from JWT", "error", err.Error())
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}
	md := model.Metadata{}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		slog.Error("read metadata body getdatahandler", "error", err.Error())
		return
	}
	err = json.Unmarshal(body, &md)
	if err != nil {
		slog.Error("unmarshal receive metadata", "error", err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = handler.dm.DeleteData(request.Context(), md, userID)
	if err != nil {
		slog.Error("delete data from storage", "error", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)

}
