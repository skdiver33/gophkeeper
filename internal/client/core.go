package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/skdiver33/gophkeeper/model"
	"github.com/skdiver33/gophkeeper/protocol"
)

type KeeperClient struct {
	ClientUser *model.User
	NWClient   *http.Client
	Config     *KeeperClientConfig
	JWT        string
	UserData   *[]model.Metadata
	CryptKey   []byte
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
	return &KeeperClient{ClientUser: &model.User{}, Config: conf, NWClient: nclient}, nil
}

func Run() {
	slog.Info("Starting client...")

	termCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	client, err := NewKeeperClient()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	client.RunUserShell(termCtx)
}

func (client *KeeperClient) UserAuth(regData *model.User, url string) error {
	data, err := json.Marshal(regData)
	if err != nil {
		return err
	}
	requestBody := bytes.NewBuffer(data)
	resp, err := client.NWClient.Post(client.Config.ServerAddr+url, "Content-Type: application/json", requestBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", string(body))
	}
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "jwt" {
			client.JWT = cookie.Value
			break
		}
	}
	return nil
}

func (client *KeeperClient) SendData(pkg *protocol.ProtocolPackage) error {

	data, err := json.Marshal(pkg)
	if err != nil {
		return err
	}

	var requestBody bytes.Buffer

	zw := gzip.NewWriter(&requestBody)
	if _, err := zw.Write(data); err != nil {
		return fmt.Errorf("error compress data %w", err)
	}
	if err := zw.Close(); err != nil {
		return fmt.Errorf("error close zip writer. error: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, client.Config.ServerAddr+"/data", &requestBody)
	if err != nil {
		return fmt.Errorf("error! create request. error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+client.JWT)
	req.Header.Set("Content-Encoding", "gzip")

	response, err := client.NWClient.Do(req)

	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error send data on server. Response code %d ", response.StatusCode)
	}

	return nil
}

func (client *KeeperClient) GetAllData() (*[]model.Metadata, error) {

	req, err := http.NewRequest(http.MethodGet, client.Config.ServerAddr+"/alldata", nil)
	if err != nil {
		return nil, fmt.Errorf("error! create request. error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+client.JWT)

	response, err := client.NWClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error get all data from server. Response code %d ", response.StatusCode)
	}

	var res []model.Metadata

	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (client *KeeperClient) GetData(md model.Metadata) (*protocol.ProtocolPackage, error) {

	data, err := json.Marshal(md)
	if err != nil {
		return nil, err
	}
	requestBody := bytes.NewBuffer(data)

	req, err := http.NewRequest(http.MethodGet, client.Config.ServerAddr+"/data", requestBody)
	if err != nil {
		return nil, fmt.Errorf("error! create request. error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+client.JWT)

	response, err := client.NWClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error get all data from server. Response code %d ", response.StatusCode)
	}

	var res protocol.ProtocolPackage

	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (client *KeeperClient) DeleteData(md model.Metadata) error {

	data, err := json.Marshal(md)
	if err != nil {
		return err
	}
	requestBody := bytes.NewBuffer(data)

	req, err := http.NewRequest(http.MethodDelete, client.Config.ServerAddr+"/data", requestBody)
	if err != nil {
		return fmt.Errorf("error! create request. error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+client.JWT)

	response, err := client.NWClient.Do(req)

	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error delete data from server. Response code %d ", response.StatusCode)
	}

	return nil
}
