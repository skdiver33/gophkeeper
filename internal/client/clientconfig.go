package client

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	Version   string
	BuildDate string
)

type KeeperClientConfig struct {
	ServerAddr string `json:"serveraddr" env:"SRV_ADDRESS"`
	CertFile   string `json:"certfile" env:"CERT_FILE_PATH"`
}

func NewKeeperClientConfig() (*KeeperClientConfig, error) {
	config := KeeperClientConfig{}
	clienFlags := flag.NewFlagSet("Clientt config flags", 0)
	clienFlags.StringVar(&config.ServerAddr, "s", "https://localhost:4443", "keeper server address in form https://ip:port.")
	clienFlags.StringVar(&config.CertFile, "cert_file", "../../tls/server.crt", "path to file with TLS certificate")
	clienFlags.BoolFunc("version", "show information about client", func(s string) error {
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Build date: %s\n", BuildDate)
		return nil
	})
	clienFlags.Parse(os.Args[1:])
	cleanenv.ReadEnv(&config)

	if config.CertFile == "" {
		return nil, fmt.Errorf("client configure error. Setup cert file for TLS server")
	}

	return &config, nil
}
