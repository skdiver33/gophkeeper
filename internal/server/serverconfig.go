package server

import (
	"flag"
	"fmt"
	"os"
)

type KeeperServerConfig struct {
	ListenAddr string `json:"address" env:"ADDRESS"`
	KeyPath    string `json:"keypath" env:"SECRET_KEY_PATH"`
	CertFile   string `json:"certfile" env:"CERT_FILE_PATH"`
	DBAddress  string `json:"dbaddress" env:"DB_ADDRESS"`
}

func NewKeeperServerConfig() (*KeeperServerConfig, error) {
	config := KeeperServerConfig{}

	serverFlags := flag.NewFlagSet("Server config flags", 0)
	serverFlags.StringVar(&config.ListenAddr, "a", "localhost:4443", "adress for start server in form ip:port. default localhost:4443")
	serverFlags.StringVar(&config.KeyPath, "secr_key", "../../tls/server.key", "path to file with private key for TLS. default key in current dirrectory")
	serverFlags.StringVar(&config.CertFile, "cert_path", "../../tls/server.crt", "path to file with server certificate for TLS. default cert in current dirrectory")
	serverFlags.StringVar(&config.DBAddress, "db_addr", "postgres://keeper:secret@localhost:5432/keepermd?sslmode=disable", "DB connection address. Default empty.")
	serverFlags.Parse(os.Args[1:])
	//cleanenv.ReadEnv(&config)

	if config.CertFile == "" || config.KeyPath == "" {
		return nil, fmt.Errorf("server configure error. Setup cert file and key file for TLS server")
	}

	return &config, nil
}
