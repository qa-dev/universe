package storage

import (
	"testing"

	"github.com/qa-dev/universe/config"
	"github.com/stretchr/testify/assert"
)

func init() {
	config.SetTestDitectory()
}

func TestNewStorage(t *testing.T) {
	storage := NewStorage()
	assert.NotNil(t, storage.Data)
	assert.NotNil(t, storage.Mutex)
}
