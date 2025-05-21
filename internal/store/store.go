package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"ppAuthService/internal/config"
	"ppAuthService/internal/entity"

	repoDto "ppAuthService/internal/repository/dto"
	repoErr "ppAuthService/internal/repository/err"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	addUserQuery = `
INSERT INTO "user" (login,password) 
VALUES ($1, $2) RETURNING user_id;`
	getUserByLoginQuery = `
SELECT * FROM "user" 
WHERE login=$1;`
	updateUserQuery = `
UPDATE "user" SET 
login = CASE WHEN $2::character varying IS NULL THEN login ELSE $2 END,
password = CASE WHEN $3::character varying IS NULL THEN password ELSE $3 END
WHERE user_id=$1
RETURNING user_id;`
	removeUserQuery = `
DELETE FROM "user" WHERE user_id=$1 RETURNING user_id;`
	addRefreshTokenWithRefreshTokenIdQuery = `
INSERT INTO refresh_token VALUES ($1,$2,$3,$4,$5);`
	getRefreshTokenQuery = `
SELECT * FROM refresh_token 
WHERE refresh_token_id=$1;`
	revokeRefreshTokensByUserIdAndDeviceCodeQuery = `
UPDATE refresh_token 
SET is_revoke=true
WHERE user_id=$1 AND ($2::character varying IS NULL OR device_code=$2);`
	revokeRefreshTokenByRefreshTokenIdQuery = `
UPDATE refresh_token 
SET is_revoke=true
WHERE refresh_token_id = $1;`
	removeRefreshTokensByExpirationAtQuery = `
DELETE FROM refresh_token
WHERE expiration_at < $1;`
)

type Store struct {
	pool *pgxpool.Pool
	lg   *slog.Logger
}

func MustNew(lg *slog.Logger, cfg *config.Store) *Store {
	connString := fmt.Sprintf(
		`user=%s password=%s host=%s port=%d dbname=%s sslmode=%s pool_max_conns=%d pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s`,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
		cfg.PoolMaxConns,
		cfg.PoolMaxConnLifetime.String(),
		cfg.PoolMaxConnIdleTime.String(),
	)
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatalf("failed to initialize store: %v\n", err)
	}
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("failed to initialize store: %v\n", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("failed to initialize store: %v\n", err)
	}
	return &Store{
		pool,
		lg,
	}
}

func (s *Store) AddUser(dto *repoDto.AddUser) (*uuid.UUID, error) {
	userId := new(uuid.UUID)
	err := s.pool.QueryRow(context.Background(), addUserQuery, dto.Login, dto.Password).Scan(userId)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "store.AddUser"))
		if pgError, ok := err.(*pgconn.PgError); ok && pgError.Code == "23505" {
			return nil, repoErr.ErrUniqueViolation
		}
		return nil, repoErr.ErrInternalServerError

	}
	return userId, nil
}
func (s *Store) GetUserByLogin(login string) (*entity.User, error) {
	user := new(entity.User)
	err := s.pool.QueryRow(context.Background(), getUserByLoginQuery, login).Scan(&user.UserId, &user.Login, &user.Password)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "store.GetUserByLogin"))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repoErr.ErrRecordNotFound
		}
		return nil, repoErr.ErrInternalServerError
	}
	return user, err
}
func (s *Store) UpdateUser(dto *repoDto.UpdateUser) error {
	userId := new(uuid.UUID)
	err := s.pool.QueryRow(context.Background(), updateUserQuery, dto.UserId, dto.Login, dto.Password).Scan(userId)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "store.UpdateUser"))
		if errors.Is(err, sql.ErrNoRows) {
			return repoErr.ErrRecordNotFound
		}
		return repoErr.ErrInternalServerError
	}
	return nil
}
func (s *Store) RemoveUser(userId *uuid.UUID) error {
	err := s.pool.QueryRow(context.Background(), removeUserQuery, userId).Scan(userId)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "store.RemoveUser"))
		if errors.Is(err, sql.ErrNoRows) {
			return repoErr.ErrRecordNotFound
		}
		return repoErr.ErrInternalServerError
	}
	return nil
}

func (s *Store) AddRefreshTokenWithRefreshTokenId(dto *repoDto.AddRefreshTokenWithRefreshTokenId) error {
	_, err := s.pool.Exec(context.Background(), addRefreshTokenWithRefreshTokenIdQuery, dto.RefreshTokenId, dto.UserId, dto.DeviceCode, dto.ExpirationAt, dto.IsRevoke)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "store.AddRefreshTokenWithRefreshTokenId"))
		return repoErr.ErrInternalServerError
	}
	return nil
}
func (s *Store) GetRefreshToken(refreshTokenId *uuid.UUID) (*entity.RefreshToken, error) {
	refreshToken := new(entity.RefreshToken)
	err := s.pool.QueryRow(context.Background(), getRefreshTokenQuery, refreshTokenId).Scan(&refreshToken.RefreshTokenId, &refreshToken.UserId, &refreshToken.DeviceCode, &refreshToken.ExpirationAt, &refreshToken.IsRevoke)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "store.GetRefreshToken"))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repoErr.ErrRecordNotFound
		}
		return nil, repoErr.ErrInternalServerError
	}
	return refreshToken, nil
}
func (s *Store) RevokeRefreshTokenByRefreshTokenId(refreshTokenId *uuid.UUID) error {
	_, err := s.pool.Exec(context.Background(), revokeRefreshTokenByRefreshTokenIdQuery, refreshTokenId)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "store.RevokeRefreshTokenByRefreshTokenId"))
		return repoErr.ErrInternalServerError
	}
	return nil
}
func (s *Store) RevokeRefreshTokensByUserIdAndDeviceCode(dto *repoDto.RevokeRefreshTokensByUserIdAndDeviceCode) error {
	_, err := s.pool.Exec(context.Background(), revokeRefreshTokensByUserIdAndDeviceCodeQuery, dto.UserId, dto.DeviceCode)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "store.RevokeRefreshTokensByUserIdAndDeviceCode"))
		return repoErr.ErrInternalServerError
	}
	return nil
}
func (s *Store) RemoveRefreshTokensByExpirationAt(now time.Time) (int64, error) {
	result, err := s.pool.Exec(context.Background(), removeRefreshTokensByExpirationAtQuery, now)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "store.RemoveRefreshTokensByExpirationAt"))
		return -1, repoErr.ErrInternalServerError
	}
	return result.RowsAffected(), nil
}
