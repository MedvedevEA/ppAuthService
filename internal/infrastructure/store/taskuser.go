package store

import (
	"context"
	"fmt"
	"ppAuthService/internal/entity"
	repoStoreDto "ppAuthService/internal/repository/repostore/dto"
	repoStoreErr "ppAuthService/internal/repository/repostore/err"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	addTaskUserQuery = `
INSERT INTO task_user (task_id, user_id) 
VALUES ($1, $2) RETURNING *;`
	getTaskUsersQuery = `
SELECT * 
FROM task_user 
WHERE ($3::uuid IS NULL OR task_id=$3) AND ($4::uuid IS null OR user_id=$4) 
OFFSET $1 LIMIT $2;`
	removeTaskUserQuery = `
DELETE FROM task_user
WHERE task_user_id=$1;`
)

func (s *Store) AddTaskUser(dto *repoStoreDto.AddTaskUser) (*entity.TaskUser, error) {
	taskUser := new(entity.TaskUser)
	err := s.pool.QueryRow(context.Background(), addTaskUserQuery, dto.TaskId, dto.UserId).Scan(&taskUser.TaskUserId, &taskUser.UserId, &taskUser.TaskId)
	if err != nil {
		if pgError, ok := err.(*pgconn.PgError); ok && pgError.Code == "23505" {
			return nil, fmt.Errorf("store.AddTaskUser: %w (%v)", repoStoreErr.ErrUniqueViolation, err)

		}
		return nil, fmt.Errorf("store.AddTaskUser: %w (%v)", repoStoreErr.ErrInternalServerError, err)

	}
	return taskUser, nil
}
func (s *Store) GetTaskUsers(dto *repoStoreDto.GetTaskUsers) ([]*entity.TaskUser, error) {
	const op = "store.GetTaskUsers"
	rows, err := s.pool.Query(context.Background(), getTaskUsersQuery, dto.Offset, dto.Limit, dto.TaskId, dto.UserId)
	if err != nil {
		return nil, fmt.Errorf("store.GetTaskUsers: %w (%v)", repoStoreErr.ErrInternalServerError, err)
	}
	defer rows.Close()
	taskUsers := make([]*entity.TaskUser, 0)
	for rows.Next() {
		taskUser := new(entity.TaskUser)
		err := rows.Scan(&taskUser.TaskUserId, &taskUser.TaskId, &taskUser.UserId)
		if err != nil {
			return nil, fmt.Errorf("store.GetTaskUsers: %w (%v)", repoStoreErr.ErrInternalServerError, err)
		}
		taskUsers = append(taskUsers, taskUser)
	}
	return taskUsers, nil
}

func (s *Store) RemoveTaskUser(TaskUserId *uuid.UUID) error {
	const op = "store.RemoveTaskUser"
	result, err := s.pool.Exec(context.Background(), removeTaskUserQuery, TaskUserId)
	if err != nil {
		return fmt.Errorf("store.RemoveTaskUser: %w (%v)", repoStoreErr.ErrInternalServerError, err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("store.RemoveTaskUser: %w (%v)", repoStoreErr.ErrRecordNotFound, err)
	}
	return nil
}
