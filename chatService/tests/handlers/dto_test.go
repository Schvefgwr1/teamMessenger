package handlers

import (
	"encoding/json"
	"testing"

	"chatService/internal/handlers/dto"
	fc "common/contracts/file-contracts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"time"
)

func TestGetChatMessage_MarshalJSON(t *testing.T) {
	fileID1 := 1
	fileID2 := 2
	files := []*fc.File{
		{ID: fileID1, Name: "file1.jpg"},
		{ID: fileID2, Name: "file2.jpg"},
	}
	updatedAt := time.Now()

	msg := dto.GetChatMessage{
		ID:        uuid.New(),
		ChatID:    uuid.New(),
		SenderID:  func() *uuid.UUID { id := uuid.New(); return &id }(),
		Content:   "Test message",
		CreatedAt: time.Now(),
		UpdatedAt: &updatedAt,
		Files:     &files,
	}

	data, err := json.Marshal(msg)
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	var unmarshaled dto.GetChatMessage
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	assert.Equal(t, msg.ID, unmarshaled.ID)
	assert.Equal(t, msg.Content, unmarshaled.Content)
}

func TestGetChatMessage_MarshalJSON_NilFiles(t *testing.T) {
	updatedAt := time.Now()
	msg := dto.GetChatMessage{
		ID:        uuid.New(),
		ChatID:    uuid.New(),
		SenderID:  func() *uuid.UUID { id := uuid.New(); return &id }(),
		Content:   "Test message",
		CreatedAt: time.Now(),
		UpdatedAt: &updatedAt,
		Files:     nil,
	}

	data, err := json.Marshal(msg)
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	var unmarshaled dto.GetChatMessage
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	assert.Equal(t, msg.ID, unmarshaled.ID)
}

func TestChatResponse_MarshalJSON(t *testing.T) {
	description := "Test description"
	avatarFileID := 1
	avatarFile := &fc.File{ID: avatarFileID, Name: "avatar.jpg"}

	chat := dto.ChatResponse{
		ID:           uuid.New(),
		Name:         "Test Chat",
		IsGroup:      true,
		Description:  &description,
		AvatarFileID: &avatarFileID,
		AvatarFile:   avatarFile,
		CreatedAt:    time.Now(),
	}

	data, err := json.Marshal(chat)
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	var unmarshaled dto.ChatResponse
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	assert.Equal(t, chat.ID, unmarshaled.ID)
	assert.Equal(t, chat.Name, unmarshaled.Name)
}

func TestChatResponse_MarshalJSON_WithoutAvatar(t *testing.T) {
	chat := dto.ChatResponse{
		ID:        uuid.New(),
		Name:      "Test Chat",
		IsGroup:   false,
		CreatedAt: time.Now(),
	}

	data, err := json.Marshal(chat)
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	var unmarshaled dto.ChatResponse
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	assert.Equal(t, chat.ID, unmarshaled.ID)
}
