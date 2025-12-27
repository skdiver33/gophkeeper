package protocol

import (
	"crypto/aes"
	"crypto/rand"
	"testing"

	"github.com/skdiver33/gophkeeper/model"
)

var (
	data      []byte
	cryptData []byte
	key       []byte
)

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func init() {
	key, _ = generateRandom(2 * aes.BlockSize)
}

func TestCreateProtoPackage(t *testing.T) {
	type args struct {
		data     []byte
		dt       model.DataTypes
		descript string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "positive test #1",
			args:    args{data: []byte(`{Login:"joe",Password:"doe"}`), dt: model.AuthDataType, descript: "auth"},
			wantErr: false,
		},
		{
			name:    "positive test #2",
			args:    args{data: []byte(`"Hello world"`), dt: model.FileType, descript: "file"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateProtoPackage(tt.args.data, tt.args.dt, tt.args.descript)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CreateProtoPackage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

		})
	}
}

func TestProtocolPackage_CryptPkgData(t *testing.T) {
	data = []byte(`{Login:"joe",Password:"doe"}`)

	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "positive test",
			data:    data,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg, _ := CreateProtoPackage(data, model.AuthDataType, "test data")
			if err := pkg.CryptPkgData(key); (err != nil) != tt.wantErr {
				t.Errorf("ProtocolPackage.CryptPkgData() error = %v, wantErr %v", err, tt.wantErr)
			}
			cryptData = pkg.Data
		})
	}
}

func TestProtocolPackage_DecryptPkgData(t *testing.T) {

	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "positive test",
			data:    cryptData,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg, _ := CreateProtoPackage(data, model.AuthDataType, "test data")
			pkg.Data = tt.data
			if err := pkg.DecryptPkgData(key); (err != nil) != tt.wantErr {
				t.Errorf("ProtocolPackage.DecryptPkgData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
