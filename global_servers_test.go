package memcache_locks

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewMemcacheLockByServers(t *testing.T) {
	envName := "STORAGE_LOCK_MEMCACHE_SERVERS"
	servers := os.Getenv(envName)
	if servers == "" {
		t.Skip("STORAGE_LOCK_MEMCACHE_SERVERS not set, skipping memcache test")
	}

	ctx := context.Background()
	lockId := "test-memcache-lock"

	factory, err := GetMemcacheLockFactoryByServers(ctx, []string{servers})
	assert.Nil(t, err)
	assert.NotNil(t, factory)

	lock, err := NewMemcacheLockByServers(ctx, []string{servers}, lockId)
	assert.Nil(t, err)
	assert.NotNil(t, lock)

	// 加锁
	ownerId := "test-owner"
	err = lock.Lock(ctx, ownerId)
	assert.Nil(t, err)

	// 解锁
	err = lock.UnLock(ctx, ownerId)
	assert.Nil(t, err)
}
