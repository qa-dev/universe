package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStorage(t *testing.T) {
	storage := NewStorage()
	assert.NotNil(t, storage.Data)
	assert.NotNil(t, storage.Mutex)
}
