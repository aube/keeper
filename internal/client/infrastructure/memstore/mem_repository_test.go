package memstore

import (
	"context"
	"testing"

	"github.com/aube/keeper/internal/client/utils/apperrors"
	"github.com/stretchr/testify/assert"
)

func TestMemoryRepository(t *testing.T) {
	repo, err := NewMemoryRepository()
	assert.NoError(t, err)

	ctx := context.Background()

	t.Run("Save and Get", func(t *testing.T) {
		user := "testuser"
		token := "testtoken123"

		err := repo.Save(ctx, user, token)
		assert.NoError(t, err)

		retrievedToken, err := repo.Get(ctx, user)
		assert.NoError(t, err)
		assert.Equal(t, token, retrievedToken)
	})

	t.Run("Get non-existent", func(t *testing.T) {
		_, err := repo.Get(ctx, "nonexistent")
		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrTokenNotFound, err)
	})

	t.Run("Delete", func(t *testing.T) {
		user := "todelete"
		token := "todelete123"

		err := repo.Save(ctx, user, token)
		assert.NoError(t, err)

		err = repo.Delete(ctx, user)
		assert.NoError(t, err)

		_, err = repo.Get(ctx, user)
		assert.Error(t, err)
	})
}
