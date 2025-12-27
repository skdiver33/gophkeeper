package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/jwtauth/v5"
	"github.com/magiconair/properties/assert"
	"github.com/skdiver33/gophkeeper/internal/auth"
	"github.com/skdiver33/gophkeeper/internal/datamanager"
	"github.com/skdiver33/gophkeeper/internal/server/handler"
	"github.com/skdiver33/gophkeeper/internal/usermanager"
	"github.com/skdiver33/gophkeeper/model"
	"github.com/skdiver33/gophkeeper/protocol"
	"github.com/skdiver33/gophkeeper/storage"
)

var (
	serverHandler *handler.ServerHandler
	store         *storage.SQLStorage
	token         []string
	au            *auth.Auth
)

func init() {
	store, _ = storage.NewSQLStorage("postgres://keeper:secret@localhost:5432/keepermd?sslmode=disable")
	au = auth.NewAuth("test_key")
	um := usermanager.NewUserManager(store, au)
	dm := datamanager.NewDataManager(store)
	serverHandler = handler.NewServerHandler(um, dm)
	token = make([]string, 0)
}

func TestServerHandler_UserRegisterHandler(t *testing.T) {

	type want struct {
		code        int
		wantCookies bool
	}
	tests := []struct {
		name        string
		requestData []byte
		want        want
	}{
		{
			name:        "positive register user",
			requestData: []byte(`{"login":"user","password":"user"}`),
			want: want{
				code:        200,
				wantCookies: true,
			},
		},
		{
			name:        "positive register user",
			requestData: []byte(`{"login":"admin","password":"admin"}`),
			want: want{
				code:        200,
				wantCookies: true,
			},
		},
		{
			name:        "user already exist",
			requestData: []byte(`{"login":"user","password":"user"}`),
			want: want{
				code:        409,
				wantCookies: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewReader(tt.requestData))
			request.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()
			serverHandler.UserRegisterHandler(w, request)

			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			defer res.Body.Close()
			cookies := res.Cookies()
			isCookieExist := len(cookies) > 0
			assert.Equal(t, tt.want.wantCookies, isCookieExist)

		})
	}

}

func TestServerHandler_UserLoginHandler(t *testing.T) {
	type want struct {
		code        int
		wantCookies bool
	}
	tests := []struct {
		name        string
		requestData []byte
		want        want
	}{
		{
			name:        "positive auth user",
			requestData: []byte(`{"login":"user","password":"user"}`),
			want: want{
				code:        200,
				wantCookies: true,
			},
		},
		{
			name:        "positive register user",
			requestData: []byte(`{"login":"admin","password":"admin"}`),
			want: want{
				code:        200,
				wantCookies: true,
			},
		},
		{
			name:        "wrong password",
			requestData: []byte(`{"login":"user","password":"123"}`),
			want: want{
				code:        401,
				wantCookies: false,
			},
		},
		{
			name:        "wrong request type",
			requestData: []byte(`{"name":"user","passd":"123}`),
			want: want{
				code:        400,
				wantCookies: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewReader(tt.requestData))
			request.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()
			serverHandler.UserLoginHandler(w, request)

			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			defer res.Body.Close()
			cookies := res.Cookies()
			isCookieExist := len(cookies) > 0
			assert.Equal(t, tt.want.wantCookies, isCookieExist)
			if len(cookies) > 0 {
				for _, item := range cookies {
					if item.Name == "jwt" {
						token = append(token, item.Value)
					}
				}
			}
		})
	}

}

func TestServerHandler_LoadDataHandler(t *testing.T) {
	testData := model.AuthData{Login: "user", Password: "secret"}
	testBytes, _ := testData.ToBinary()
	pkg, _ := protocol.CreateProtoPackage(testBytes, model.AuthDataType, "test data")
	tests := []struct {
		name        string
		requestData *protocol.ProtocolPackage
		userToken   string
		wantCode    int
	}{
		{
			name:        "positive order load",
			requestData: pkg,
			userToken:   token[0],
			wantCode:    200,
		},
		{
			name:        "data already load ",
			requestData: pkg,
			userToken:   token[0],
			wantCode:    409,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tok, err := jwtauth.VerifyToken(au.GetBaseToken(), tt.userToken)
			if err != nil {
				log.Print("error create tok")
			}
			ctx2 := context.WithValue(ctx, jwtauth.TokenCtxKey, tok)
			data, err := json.Marshal(tt.requestData)
			if err != nil {
				log.Println(err.Error())
			}
			request := httptest.NewRequestWithContext(ctx2, http.MethodPost, "/data/", bytes.NewReader([]byte(data)))
			request.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()
			serverHandler.LoadDataHandler(w, request)
			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.wantCode, res.StatusCode)

		})
	}
	if err := store.CloseAndClean(); err != nil {
		t.Error("error close DB")
	}

}
