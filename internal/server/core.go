package server

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/skdiver33/gophkeeper/model"
)

type DataStorage interface {
	CreateData()
	ReadData()
	UpdateData()
	DeleteData()
}

type MetaDataStorage interface {
	AddMetaData()
	ListMetaData()
	DeleteMetaData()
	GetMetaData()
}

type KeeperServer struct {
	ServerConfig  *KeeperServerConfig
	DataStore     DataStorage
	MetaDataStore MetaDataStorage
}

func NewKeeperServer() (*KeeperServer, error) {
	config, err := NewKeeperServerConfig()
	if err != nil {
		return nil, err
	}
	return &KeeperServer{ServerConfig: config}, nil
}

func Hello(w http.ResponseWriter, r *http.Request) {
	user := model.AuthData{}
	data, _ := io.ReadAll(r.Body)
	json.Unmarshal(data, &user)

	fmt.Fprintln(w, "Hello tls")
}

func Run() {
	ks, err := NewKeeperServer()
	if err != nil {
		log.Println(err)
		return
	}

	http.HandleFunc("/", Hello)
	server := &http.Server{
		Addr:      ks.ServerConfig.ListenAddr,
		TLSConfig: &tls.Config{},
	}
	log.Println("Starting server on :443")
	err = server.ListenAndServeTLS(ks.ServerConfig.CertFile, ks.ServerConfig.KeyPath)
	if err != nil {
		fmt.Println(err)
	}

}
