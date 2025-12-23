package usermanager_test

import (
	"log"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/skdiver33/gophkeeper/internal/auth"
	"github.com/skdiver33/gophkeeper/internal/usermanager"
	"github.com/skdiver33/gophkeeper/mocks"
	"github.com/skdiver33/gophkeeper/model"
	"github.com/stretchr/testify/require"
)

func TestUserManager_UserRegister(t *testing.T) {

	type StorageData struct {
		Login string
		Paswd string
		ID    int
		err   error
	}
	data := []StorageData{
		{
			Login: "user",
			Paswd: "user",
			ID:    1,
			err:   nil,
		},
		{
			Login: "admin",
			Paswd: "admin",
			ID:    2,
			err:   nil,
		},
		{
			Login: "joe",
			Paswd: "dow",
			ID:    3,
			err:   nil,
		},
		{
			Login: "user",
			Paswd: "passwd",
			ID:    -1,
			err:   usermanager.ErrUserAlreadyExist,
		},
	}

	auth := auth.NewAuth()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockUserStorageInterface(ctrl)

	for _, val := range data {
		addUser := &model.User{Login: val.Login, Password: val.Paswd}
		m.EXPECT().AddUser(t.Context(), addUser).Return(val.ID, val.err).AnyTimes()
		var ret *model.User
		if val.err == usermanager.ErrUserAlreadyExist {
			ret = addUser
		}
		m.EXPECT().GetUser(t.Context(), val.Login, val.Paswd).Return(ret, val.err)
	}

	userManager := usermanager.NewUserManager(m, auth)

	type want struct {
		returnID  int
		wantError bool
	}
	tests := []struct {
		name     string
		testUser model.User
		result   want
	}{
		{
			name:     "positive register #1",
			testUser: model.User{Login: "user", Password: "user"},
			result:   want{returnID: 1, wantError: false},
		},
		{
			name:     "positive register #2",
			testUser: model.User{Login: "admin", Password: "admin"},
			result:   want{returnID: 2, wantError: false},
		},
		{
			name:     "positive register #3",
			testUser: model.User{Login: "joe", Password: "dow"},
			result:   want{returnID: 3, wantError: false},
		},
		{
			name:     "negative register #1",
			testUser: model.User{Login: "user", Password: "passwd"},
			result:   want{returnID: -1, wantError: true},
		},
	}

	for _, test := range tests {
		log.Println(test.name)
		val, err := userManager.UserRegister(t.Context(), &test.testUser)
		if val != nil {
			require.Equal(t, val.ID, test.result.returnID)
		}

		if test.result.wantError {
			require.Error(t, err)
		}
	}

}

func TestUserManager_UserAuth(t *testing.T) {
	type StorageData struct {
		Login string
		Paswd string
		ID    int
		err   error
	}
	data := []StorageData{
		{
			Login: "user",
			Paswd: "user",
			ID:    1,
			err:   nil,
		},
		{
			Login: "user",
			Paswd: "passwd",
			ID:    -1,
			err:   usermanager.ErrUserWithCredNotFound,
		},
	}

	auth := auth.NewAuth()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockUserStorageInterface(ctrl)

	for _, val := range data {
		addUser := &model.User{Login: val.Login, Password: val.Paswd}
		m.EXPECT().AddUser(t.Context(), addUser).Return(val.ID, val.err).AnyTimes()
		var ret *model.User
		if val.err == usermanager.ErrUserAlreadyExist {
			ret = addUser
		}
		m.EXPECT().GetUser(t.Context(), val.Login, val.Paswd).Return(ret, val.err)
	}

	type want struct {
		wantError bool
	}
	tests := []struct {
		name     string
		testUser model.User
		result   want
	}{
		{
			name:     "positive auth #1",
			testUser: model.User{Login: "user", Password: "user"},
			result:   want{wantError: false},
		},
		{
			name:     "negative auth #1",
			testUser: model.User{Login: "user", Password: "passwd"},
			result:   want{wantError: true},
		},
	}
	userManager := usermanager.NewUserManager(m, auth)

	for _, test := range tests {
		log.Println(test.name)
		val, err := userManager.UserAuth(t.Context(), &test.testUser)
		if test.result.wantError {
			require.ErrorIs(t, err, usermanager.ErrUserWithCredNotFound)
		}
		require.NotNil(t, val)
	}

}
