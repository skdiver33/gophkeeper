package server

import (
	"errors"
	"flag"
	"os"
)

type KeeperServerConfig struct {
	ListenAddr string `json:"address" env:"ADDRESS"`
	KeyPath    string `json:"keypath" env:"SECRET_KEY_PATH"`
	CertFile   string `json:"certfile" env:"CERT_FILE_PATH"`
	DBAddress  string `json:"dbaddress" env:"DB_ADDRESS"`
	SignKey    string `json:"signkey" env:"SIGN_KEY"`
}

var (
	KeeperServerSetupTLSError = errors.New("path for cert file ok key not specified")
	KeeperServerSignKeyError  = errors.New("sign key for JWT not specified")
)

func NewKeeperServerConfig() (*KeeperServerConfig, error) {
	config := KeeperServerConfig{}

	serverFlags := flag.NewFlagSet("Server config flags", 0)
	serverFlags.StringVar(&config.ListenAddr, "a", "localhost:4443", "adress for start server in form ip:port. default localhost:4443")
	serverFlags.StringVar(&config.KeyPath, "secr_key", "../../tls/server.key", "path to file with private key for TLS. default key in current dirrectory")
	serverFlags.StringVar(&config.CertFile, "cert_path", "../../tls/server.crt", "path to file with server certificate for TLS. default cert in current dirrectory")
	serverFlags.StringVar(&config.SignKey, "sign_key", "blablabla", "sign key for creating JWT. Important default not specified and server automatic shutdown.")
	serverFlags.StringVar(&config.DBAddress, "db_addr", "postgres://keeper:secret@localhost:5432/keepermd?sslmode=disable", "DB connection address. Default empty.")
	serverFlags.Parse(os.Args[1:])
	//cleanenv.ReadEnv(&config)

	if config.CertFile == "" || config.KeyPath == "" {
		return nil, KeeperServerSetupTLSError
	}

	if config.SignKey == "" {
		return nil, KeeperServerSignKeyError
	}

	return &config, nil
}
