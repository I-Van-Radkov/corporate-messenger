package adapter

import (
	"sync"

	"github.com/google/uuid"
)

type OfflineStorage struct {
	messages map[uuid.UUID][][]byte
	mu       sync.Mutex
}

func NewOfflineStorage() *OfflineStorage {
	return &OfflineStorage{
		messages: make(map[uuid.UUID][][]byte),
	}
}

func (s *OfflineStorage) SendMessage(userId uuid.UUID, data []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.messages[userId] = append(s.messages[userId], data)
}

func (s *OfflineStorage) GetMessages(userId uuid.UUID) ([][]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	messages, exists := s.messages[userId]
	if !exists {
		return [][]byte{}, nil
	}

	delete(s.messages, userId)

	return messages, nil
}
