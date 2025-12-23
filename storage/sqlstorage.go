package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/skdiver33/gophkeeper/model"

	"github.com/golang-migrate/migrate/v4"
)

type SQLStorage struct {
	config *SQLStorageConfig
	db     *sql.DB
}

type SQLStorageConfig struct {
	DBAddress string
}

func NewSQLStorageConfig(address string) *SQLStorageConfig {
	storageConfig := SQLStorageConfig{DBAddress: address}
	return &storageConfig
}

func NewSQLStorage(address string) (*SQLStorage, error) {
	newStorage := SQLStorage{}
	newStorage.config = NewSQLStorageConfig(address)
	err := newStorage.InitializeConnection()
	if err != nil {
		return nil, err
	}

	migrator, err := NewMigrator(newStorage.db)
	if err != nil {
		return nil, err
	}
	err = migrator.ApplyMigrations("up")
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}
	err = newStorage.InitializeConnection()
	if err != nil {
		return nil, err
	}
	return &newStorage, nil
}

func (storage *SQLStorage) InitializeConnection() error {
	db, err := sql.Open("pgx", storage.config.DBAddress)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	if err = db.PingContext(ctx); err != nil {
		cancel()
		return err
	}
	defer cancel()
	storage.db = db
	return nil
}

func (storage *SQLStorage) CloseAndClean() error {
	migrator, err := NewMigrator(storage.db)
	if err != nil {
		return err
	}
	err = migrator.ApplyMigrations("down")
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	storage.db.Close()
	return nil
}

func (storage *SQLStorage) CloseStorage() {
	storage.db.Close()
	slog.Info("Close postgresql connection")
}

func (storage *SQLStorage) AddUser(ctx context.Context, user *model.User) (int, error) {
	id := -1
	err := storage.db.QueryRowContext(ctx, "INSERT INTO users (login, password) VALUES ($1, $2) RETURNING user_id", user.Login, user.Password).Scan((&id))
	if err != nil {
		return -1, errors.New("error get inserted id")
	}
	return int(id), nil
}
func (storage *SQLStorage) GetUser(ctx context.Context, login string, password string) (*model.User, error) {
	user := model.User{}
	row := storage.db.QueryRowContext(ctx, "SELECT * FROM users WHERE login=$1", login)
	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (storage *SQLStorage) InsertData(ctx context.Context, md model.Metadata, data []byte, userID int) error {
	md_id := -1
	tx, err := storage.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, "INSERT INTO metadata (user_id,data_type,descript,md_hash,upload_data) VALUES ($1,$2,$3,$4,$5) RETURNING md_id", userID, md.UploadType, md.Description, md.Hash, md.UploadDate.Format(time.RFC3339)).Scan(&md_id)
	if err != nil {
		return fmt.Errorf("error insert new metada %w", err)
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO users_data (md_id,user_data) VALUES ($1, $2)", md_id, data)
	if err != nil {
		return fmt.Errorf("error insert new binary data %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error insert new binary data %w", err)
	}

	return nil
}

func (storage *SQLStorage) GetData(ctx context.Context, md model.Metadata, userID int) ([]byte, error) {
	md_id := -1
	row := storage.db.QueryRowContext(ctx, "SELECT md_id FROM metadata WHERE user_id=$1 and md_hash=$2", userID, md.Hash)
	err := row.Scan(&md_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	data := make([]byte, 0)
	row = storage.db.QueryRowContext(ctx, "SELECT user_data FROM users_data WHERE md_id=$1", md_id)
	err = row.Scan(&data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return data, nil
}

func (storage *SQLStorage) GetAllData(ctx context.Context, userID int) (*[]model.Metadata, error) {

	rows, err := storage.db.QueryContext(ctx, "SELECT data_type,descript,md_hash,upload_data FROM metadata where user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("error get all metadata for user %w", err)
	}
	defer rows.Close()

	result := make([]model.Metadata, 0)

	for rows.Next() {
		curMD := model.Metadata{}
		var upload string
		if err := rows.Scan(&curMD.UploadType, &curMD.Description, &curMD.Hash, &upload); err != nil {

			return nil, fmt.Errorf("error parse result from DB %w", err)
		}
		curMD.UploadDate, err = time.Parse(time.RFC3339, upload)
		if err != nil {
			return nil, errors.New("error parse data string from  DB")
		}
		result = append(result, curMD)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error load data from DB %s", err.Error())
	}

	return &result, nil
}

func (storage *SQLStorage) GetMetaData(ctx context.Context, hash string, userID int) (*model.Metadata, error) {
	md := model.Metadata{}
	var upload string
	row := storage.db.QueryRowContext(ctx, "SELECT data_type,descript,md_hash,upload_data FROM metadata where user_id = $1 and md_hash=$2", userID, hash)

	err := row.Scan(&md.UploadType, &md.Description, &md.Hash, &upload)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	md.UploadDate, err = time.Parse(time.RFC3339, upload)
	if err != nil {
		return nil, errors.New("error parse data string from  DB")
	}

	return &md, nil
}

func (storage *SQLStorage) DeleteData(ctx context.Context, md model.Metadata, userID int) error {
	md_id := -1
	tx, err := storage.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, "DELETE FROM metadata WHERE  user_id=$1 AND data_type=$2 AND md_hash=$3 RETURNING md_id", userID, md.UploadType, md.Hash).Scan(&md_id)
	if err != nil {
		return fmt.Errorf("error delete metada %w", err)
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM users_data WHERE md_id=$1 ", md_id)
	if err != nil {
		return fmt.Errorf("error delete data %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error commit deleting data %w", err)
	}

	return nil
}
