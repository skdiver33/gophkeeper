package datamanager

import (
	"context"
	"errors"
	"fmt"

	"github.com/skdiver33/gophkeeper/model"
	"github.com/skdiver33/gophkeeper/protocol"
)

type DataStorage interface {
	InsertData(ctx context.Context, md model.Metadata, data []byte, userID int) error
	GetData(ctx context.Context, md model.Metadata, userID int) ([]byte, error)
	GetMetaData(ctx context.Context, hash string, userID int) (*model.Metadata, error)
	GetAllData(ctx context.Context, userID int) (*[]model.Metadata, error)
	DeleteData(ctx context.Context, md model.Metadata, userID int) error
}

var (
	ErrDataAlreadyLoad = errors.New("data already load")
	ErrDataNotFound    = errors.New("data not found")
)

type DataManager struct {
	DataStore DataStorage
}

func NewDataManager(ds DataStorage) *DataManager {
	return &DataManager{DataStore: ds}
}

func (dm *DataManager) LoadData(ctx context.Context, data *protocol.ProtocolPackage, userID int) error {
	mdForAdd := data.MData
	dataForAdd := data.Data
	md, err := dm.DataStore.GetMetaData(ctx, mdForAdd.Hash, userID)
	if err != nil {
		return err
	}
	if md != nil {
		return ErrDataAlreadyLoad
	}
	err = dm.DataStore.InsertData(ctx, mdForAdd, dataForAdd, userID)
	if err != nil {
		return fmt.Errorf("error insert new data: %w", err)
	}

	return nil
}

func (dm *DataManager) GetData(ctx context.Context, md model.Metadata, userID int) (*protocol.ProtocolPackage, error) {
	returnPkg := protocol.ProtocolPackage{}
	data, err := dm.DataStore.GetData(ctx, md, userID)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, ErrDataNotFound
	}
	returnPkg.Data = data
	returnPkg.MData = md
	return &returnPkg, nil
}

func (dm *DataManager) GetAllData(ctx context.Context, userID int) (*[]model.Metadata, error) {
	allData, err := dm.DataStore.GetAllData(ctx, userID)
	if err != nil {
		return nil, err
	}
	return allData, nil
}

func (dm *DataManager) DeleteData(ctx context.Context, md model.Metadata, userID int) error {
	err := dm.DataStore.DeleteData(ctx, md, userID)
	if err != nil {
		return err
	}

	return nil
}
