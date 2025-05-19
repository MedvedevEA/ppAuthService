package store

import (
	"context"
	"database/sql"
	"fmt"
	"ppAuthService/internal/entity"
	"ppAuthService/internal/repository/repostore/dto"
	repoStoreErr "ppAuthService/internal/repository/repostore/err"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	addStatusQuery = `
INSERT INTO status (name) 
VALUES ($1) RETURNING *;`
	getStatusQuery = `
SELECT * FROM status WHERE status_id=$1;`
	getStatusesQuery = `
SELECT * FROM status;`
	updateStatusQuery = `
UPDATE status 
SET name = CASE WHEN $2::character varying IS NULL THEN name ELSE $2 END
WHERE status_id=$1
RETURNING *;`
	removeStatusQuery = `
DELETE FROM status 
WHERE status_id=$1;`
)

func (s *Store) AddStatus(name string) (*entity.Status, error) {
	status := new(entity.Status)
	err := s.pool.QueryRow(context.Background(), addStatusQuery, name).Scan(&status.StatusId, &status.Name)
	if err != nil {
		return nil, fmt.Errorf("store.AddStatus: %w (%v)", repoStoreErr.ErrInternalServerError, err)
	}
	return status, nil
}
func (s *Store) GetStatus(statusId *uuid.UUID) (*entity.Status, error) {
	status := new(entity.Status)
	err := s.pool.QueryRow(context.Background(), getStatusQuery, statusId).Scan(&status.StatusId, &status.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("store.GetStatus: %w (%v)", repoStoreErr.ErrRecordNotFound, err)
		}
		return nil, fmt.Errorf("store.GetStatus: %w (%v)", repoStoreErr.ErrInternalServerError, err)
	}
	return status, nil
}
func (s *Store) GetStatuses() ([]*entity.Status, error) {
	rows, err := s.pool.Query(context.Background(), getStatusesQuery)
	if err != nil {
		return nil, fmt.Errorf("store.GetStatuses: %w (%v)", repoStoreErr.ErrInternalServerError, err)
	}
	defer rows.Close()
	statuses := make([]*entity.Status, 0)
	for rows.Next() {
		status := new(entity.Status)
		err := rows.Scan(&status.StatusId, &status.Name)
		if err != nil {
			return nil, fmt.Errorf("store.GetStatuses: %w (%v)", repoStoreErr.ErrInternalServerError, err)
		}
		statuses = append(statuses, status)
	}
	return statuses, nil
}
func (s *Store) UpdateStatus(dto *dto.UpdateStatus) (*entity.Status, error) {
	status := new(entity.Status)
	err := s.pool.QueryRow(context.Background(), updateStatusQuery, dto.StatusId, dto.Name).Scan(&status.StatusId, &status.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("store.UpdateStatus: %w (%v)", repoStoreErr.ErrRecordNotFound, err)
		}
		return nil, fmt.Errorf("store.UpdateStatus: %w (%v)", repoStoreErr.ErrInternalServerError, err)
	}
	return status, nil
}
func (s *Store) RemoveStatus(statusId *uuid.UUID) error {
	result, err := s.pool.Exec(context.Background(), removeStatusQuery, statusId)
	if err != nil {
		return fmt.Errorf("store.RemoveStatus: %w (%v)", repoStoreErr.ErrInternalServerError, err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("store.RemoveStatus: %w (%v)", repoStoreErr.ErrRecordNotFound, err)
	}
	return nil
}
