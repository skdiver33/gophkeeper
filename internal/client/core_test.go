package client

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/skdiver33/gophkeeper/model"
	"github.com/skdiver33/gophkeeper/protocol"
)

var (
	keeperClient *KeeperClient
)

func init() {
	os.Setenv("SRV_ADDRESS", "")
	keeperClient, _ = NewKeeperClient()
}

func TestKeeperClient_UserAuth(t *testing.T) {

	tests := []struct {
		name    string
		regData *model.User
		wantErr bool
	}{
		{
			name:    "positive test",
			regData: &model.User{Login: "user", Password: "user"},
			wantErr: false,
		},
		{
			name:    "already exist test",
			regData: &model.User{Login: "already_exist", Password: "user"},
			wantErr: true,
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rd := model.User{}
		if err := json.NewDecoder(r.Body).Decode(&rd); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if rd.Login == "already_exist" {
			w.WriteHeader(http.StatusConflict)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			if err := keeperClient.UserAuth(tt.regData, ts.URL); (err != nil) != tt.wantErr {
				t.Errorf("KeeperClient.UserAuth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKeeperClient_SendData(t *testing.T) {

	tests := []struct {
		name    string
		data    []byte
		wantErr bool
		mtype   model.DataTypes
		desc    string
	}{
		{
			name:    "positive add auth data",
			data:    []byte(`{Login: "user", Password: "user"}`),
			wantErr: false,
			mtype:   model.AuthDataType,
			desc:    "new data",
		},
		{
			name:    "positive add bank card data",
			data:    []byte(`{CardNumber: "123456789",ExpireDate: "01.01.2026",CSVCode: "123",CardHolder: "User USer"}`),
			wantErr: false,
			mtype:   model.BankCardType,
			desc:    "new bak card",
		},
		{
			name:    "negative add auth data",
			data:    []byte(`{Login: "user", Password: "user"}`),
			wantErr: true,
			mtype:   model.AuthDataType,
			desc:    "already_exist",
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/data", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			slog.Error("create gzip archivator", "error", err.Error())
			return
		}
		decompressBody, err := io.ReadAll(gz)
		if err != nil {
			slog.Error("decompress body", "error", err.Error())
			return
		}
		gz.Close()

		pkg := protocol.ProtocolPackage{}

		err = json.Unmarshal(decompressBody, &pkg)
		if err != nil {
			slog.Error("unmarshal receive package", "error", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if pkg.MData.Description == "already_exist" {
			w.WriteHeader(http.StatusConflict)
		}
		w.WriteHeader(http.StatusOK)
	}))

	ts := httptest.NewServer(mux)
	defer ts.Close()
	keeperClient.Config.ServerAddr = ts.URL
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg, err := protocol.CreateProtoPackage(tt.data, tt.mtype, tt.desc)
			if err != nil {
				log.Println("error create pkg")
			}
			if err := keeperClient.SendData(pkg); (err != nil) != tt.wantErr {
				t.Errorf("KeeperClient.SendData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKeeperClient_GetData(t *testing.T) {

	pkg, _ := protocol.CreateProtoPackage([]byte(`{Login:"user",Password:"user"}`), model.AuthDataType, "auth data")

	tests := []struct {
		name    string
		args    model.Metadata
		wantErr bool
	}{
		{
			name:    "positive test",
			args:    pkg.MData,
			wantErr: false,
		},
		{
			name:    "negative test",
			args:    model.Metadata{UploadDate: time.Now(), UploadType: model.AuthDataType, Description: "not_found", Hash: " HIoXAEXxR9gFPUtdb4Mox7zJo8sjr+VuE/QabH0C0K4"},
			wantErr: true,
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/data", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		md := model.Metadata{}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			slog.Error("read metadata body getdatahandler", "error", err.Error())
			return
		}
		err = json.Unmarshal(body, &md)
		if err != nil {
			slog.Error("unmarshal receive metadata", "error", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if md.Description == "not_found" {
			w.WriteHeader(http.StatusInternalServerError)
		}

		resp, err := json.Marshal(pkg)
		if err != nil {
			slog.Error("marshal response", "error", err.Error())
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}))

	ts := httptest.NewServer(mux)
	defer ts.Close()
	keeperClient.Config.ServerAddr = ts.URL

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := keeperClient.GetData(tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("KeeperClient.GetData() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got.MData, tt.args) {
				t.Errorf("KeeperClient.GetData() = %v, want %v", got.MData, tt.args)
			}
		})
	}
}
