package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/skdiver33/gophkeeper/model"
)

type KeeperClient struct {
	ClientUser model.User
	NWClient   *http.Client
	Config     *KeeperClientConfig
}

func NewKeeperClient() (*KeeperClient, error) {
	conf, err := NewKeeperClientConfig()
	if err != nil {
		return nil, err
	}
	cert := x509.NewCertPool()
	pemData, err := os.ReadFile(conf.CertFile)
	if err != nil {
		return nil, err
	}
	if !cert.AppendCertsFromPEM(pemData) {
		return nil, fmt.Errorf("failed to append cert")
	}
	nclient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: cert,
			},
		},
	}
	return &KeeperClient{ClientUser: model.User{}, Config: conf, NWClient: nclient}, nil
}

func Run() {
	client, err := NewKeeperClient()
	if err != nil {
		log.Fatal(err)
	}
	client.RunUserShell()
	//client.RegisterNewUser(nil)
}

func (client *KeeperClient) RegisterNewUser(regData *model.User) error {

	data, err := json.Marshal(regData)
	if err != nil {
		return err
	}

	requestBody := bytes.NewBuffer(data)
	resp, err := client.NWClient.Post(client.Config.ServerAddr, "Content-Type: application/json", requestBody)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)

	}
	log.Println(string(body))
	return nil
}
