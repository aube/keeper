package memstore

import (
	"context"
	"sync"

	"github.com/aube/keeper/internal/client/utils/apperrors"
	"github.com/aube/keeper/internal/client/utils/logger"
	"github.com/rs/zerolog"
)

type MemoryRepository struct {
	tokens map[string]string
	mu     sync.RWMutex
	log    zerolog.Logger
}

func NewMemoryRepository() (*MemoryRepository, error) {
	return &MemoryRepository{
		tokens: make(map[string]string),
		log:    logger.Get().With().Str("mem", "mem_repository").Logger(),
	}, nil
}

func (r *MemoryRepository) Save(ctx context.Context, user string, token string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tokens[user] = token

	return nil
}

func (r *MemoryRepository) Delete(ctx context.Context, user string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.tokens, user)

	return nil
}

func (r *MemoryRepository) Get(ctx context.Context, user string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	token, ok := r.tokens[user]

	if !ok {
		return "", apperrors.ErrTokenNotFound
	}

	return token, nil
}

// var _ appFile.FileRepository = (*MemoryRepository)(nil)
