package client

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type KeeperClientConfig struct {
	ServerAddr string `json:"serveraddr" env:"SRV_ADDRESS"`
	CertFile   string `json:"certfile" env:"CERT_FILE_PATH"`
}

func NewKeeperClientConfig() (*KeeperClientConfig, error) {
	config := KeeperClientConfig{}
	clienFlags := flag.NewFlagSet("Clientt config flags", 0)
	clienFlags.StringVar(&config.ServerAddr, "s", "https://localhost:4443", "keeper server address in form https://ip:port. default https://localhost:4443")
	clienFlags.StringVar(&config.CertFile, "cert_file", "../../tls/server.crt", "path to file with TLS certificate")

	clienFlags.Parse(os.Args[1:])
	cleanenv.ReadEnv(&config)

	if config.CertFile == "" {
		return nil, fmt.Errorf("client configure error. Setup cert file for TLS server")
	}

	return &config, nil
}
