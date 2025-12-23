package datamanager_test

import (
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/skdiver33/gophkeeper/internal/datamanager"
	"github.com/skdiver33/gophkeeper/model"
	"github.com/skdiver33/gophkeeper/protocol"

	"github.com/golang/mock/gomock"
	mocks "github.com/skdiver33/gophkeeper/mocks"
)

func TestOrderManager_LoadData(t *testing.T) {
	type testData struct {
		name    string
		testPkg *protocol.ProtocolPackage
		err     error
		wantErr bool
		getErr  error
		ID      int
	}
	data1, _ := protocol.CreateProtoPackage([]byte(`{Login:"user",Password}`), model.AuthDataType, "test data auth")
	data2, _ := protocol.CreateProtoPackage([]byte(`{CardNumber:"123456789",ExpireDate:"01.01.2026",CSVCode:123,CardHolder:"Ivanov"}`), model.BankCardType, "test bank card")
	pkgs := []testData{
		{
			name:    "positive test #1",
			testPkg: data1,
			ID:      1,
			err:     nil,
			getErr:  nil,
			wantErr: false,
		},
		{
			name:    "positive test #2",
			testPkg: data2,
			ID:      2,
			err:     nil,
			getErr:  nil,
			wantErr: false,
		},
		{
			name:    "negative test #1",
			testPkg: data1,
			ID:      1,
			err:     datamanager.ErrDataAlreadyLoad,
			getErr:  nil,
			wantErr: true,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockDataStorage(ctrl)

	for _, item := range pkgs {
		addPkg := item.testPkg
		m.EXPECT().InsertData(t.Context(), addPkg.MData, addPkg.Data, item.ID).Return(item.err).AnyTimes()

		md := &addPkg.MData
		if item.wantErr == false {
			md = nil
		}
		m.EXPECT().GetMetaData(t.Context(), addPkg.MData.Hash, item.ID).Return(md, item.getErr)
	}

	dm := datamanager.NewDataManager(m)
	for _, tt := range pkgs {
		t.Run(tt.name, func(t *testing.T) {

			gotErr := dm.LoadData(t.Context(), tt.testPkg, tt.ID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("LoadData() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("LoadData() succeeded unexpectedly")
			}
		})
	}
}

func TestOrderManager_GetData(t *testing.T) {
	type testData struct {
		name    string
		testPkg *protocol.ProtocolPackage
		err     error
		wantErr bool
		getErr  error
		ID      int
	}
	data1, _ := protocol.CreateProtoPackage([]byte(`{Login:"user",Password}`), model.AuthDataType, "test data auth")
	data2, _ := protocol.CreateProtoPackage([]byte(`{CardNumber:"123456789",ExpireDate:"01.01.2026",CSVCode:123,CardHolder:"Ivanov"}`), model.BankCardType, "test bank card")
	pkgs := []testData{
		{
			name:    "positive test #1",
			testPkg: data1,
			ID:      1,
			err:     nil,
			getErr:  nil,
			wantErr: false,
		},
		{
			name:    "positive test #2",
			testPkg: data2,
			ID:      2,
			err:     nil,
			getErr:  nil,
			wantErr: false,
		},
		{
			name:    "negative test #1",
			testPkg: data1,
			ID:      3,
			err:     datamanager.ErrDataNotFound,
			getErr:  nil,
			wantErr: true,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockDataStorage(ctrl)

	for _, item := range pkgs {
		addPkg := item.testPkg
		m.EXPECT().GetData(t.Context(), addPkg.MData, item.ID).Return(item.testPkg.Data, item.err).AnyTimes()
	}

	dm := datamanager.NewDataManager(m)
	for _, tt := range pkgs {
		t.Run(tt.name, func(t *testing.T) {

			pkg, gotErr := dm.GetData(t.Context(), tt.testPkg.MData, tt.ID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetData() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetData() succeeded unexpectedly")
			}
			assert.Equal(t, pkg, tt.testPkg)

		})
	}
}
