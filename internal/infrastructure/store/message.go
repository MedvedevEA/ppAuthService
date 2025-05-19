package store

import (
	"context"
	"database/sql"
	"fmt"
	"ppAuthService/internal/entity"
	repoStoreDto "ppAuthService/internal/repository/repostore/dto"
	repoStoreErr "ppAuthService/internal/repository/repostore/err"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	addMessageQuery = `
INSERT INTO message (task_id,user_id,"text") 
VALUES ($1, $2, $3) RETURNING *;`
	getMessageQuery = `
SELECT * FROM message WHERE message_id=$1;`
	getMessagesQuery = `
SELECT * FROM message 
WHERE $3::uuid IS null OR task_id = $3 
OFFSET $1 LIMIT $2;`
	updateMessageQuery = `
UPDATE message 
SET 
"text" = CASE WHEN $2::character varying IS NULL THEN "text" ELSE $2 END,
update_at = now()
WHERE message_id=$1
RETURNING *;`
	removeMessageQuery = `
DELETE FROM message 
WHERE message_id=$1;`
)

func (s *Store) AddMessage(dto *repoStoreDto.AddMessage) (*entity.Message, error) {
	message := new(entity.Message)
	err := s.pool.QueryRow(context.Background(), addMessageQuery, dto.TaskId, dto.UserId, dto.Text).Scan(&message.MessageId, &message.TaskId, &message.UserId, &message.Text, &message.CreateAt, &message.UpdateAt)
	if err != nil {
		return nil, fmt.Errorf("store.AddMessage: %w (%v)", repoStoreErr.ErrInternalServerError, err)
	}
	return message, nil
}
func (s *Store) GetMessage(messageId *uuid.UUID) (*entity.Message, error) {
	message := new(entity.Message)
	err := s.pool.QueryRow(context.Background(), getMessageQuery, messageId).Scan(&message.MessageId, &message.TaskId, &message.UserId, &message.Text, &message.CreateAt, &message.UpdateAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("store.GetMessage: %w (%v)", repoStoreErr.ErrRecordNotFound, err)
		}
		return nil, fmt.Errorf("store.GetMessage: %w (%v)", repoStoreErr.ErrInternalServerError, err)
	}
	return message, nil
}
func (s *Store) GetMessages(dto *repoStoreDto.GetMessages) ([]*entity.Message, error) {
	rows, err := s.pool.Query(context.Background(), getMessagesQuery, dto.Offset, dto.Limit, dto.TaskId)
	if err != nil {
		return nil, fmt.Errorf("store.GetMessages: %w (%v)", repoStoreErr.ErrInternalServerError, err)
	}
	defer rows.Close()
	var messages []*entity.Message
	for rows.Next() {
		message := new(entity.Message)
		err := rows.Scan(&message.MessageId, &message.TaskId, &message.UserId, &message.Text, &message.CreateAt, &message.UpdateAt)
		if err != nil {
			return nil, fmt.Errorf("store.GetMessages: %w (%v)", repoStoreErr.ErrInternalServerError, err)
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (s *Store) UpdateMessage(dto *repoStoreDto.UpdateMessage) (*entity.Message, error) {
	message := new(entity.Message)
	err := s.pool.QueryRow(context.Background(), updateMessageQuery, dto.MessageId, dto.Text).Scan(&message.MessageId, &message.TaskId, &message.UserId, &message.Text, &message.CreateAt, &message.UpdateAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("store.UpdateMessage: %w (%v)", repoStoreErr.ErrRecordNotFound, err)
		}
		return nil, fmt.Errorf("store.UpdateMessage: %w (%v)", repoStoreErr.ErrInternalServerError, err)
	}
	return message, nil
}

func (s *Store) RemoveMessage(messageId *uuid.UUID) error {
	result, err := s.pool.Exec(context.Background(), removeMessageQuery, messageId)
	if err != nil {
		return fmt.Errorf("store.RemoveMessage: %w (%v)", repoStoreErr.ErrInternalServerError, err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("store.RemoveMessage: %w (%v)", repoStoreErr.ErrRecordNotFound, err)
	}

	return nil

}
