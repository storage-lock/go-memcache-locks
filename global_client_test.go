package memcache_locks

import (
	"context"
	"testing"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	storage_lock "github.com/storage-lock/go-storage-lock"
	"github.com/stretchr/testify/assert"
)

func testMemcacheClient(t *testing.T) *memcache.Client {
	t.Helper()
	client := memcache.New("127.0.0.1:11211")
	if err := client.Ping(); err != nil {
		t.Skipf("memcache 127.0.0.1:11211 不可用，跳过: %v", err)
	}
	return client
}

func TestNewMemcacheLockByClient(t *testing.T) {
	client := testMemcacheClient(t)
	defer client.Close()

	ctx := context.Background()
	factory, err := GetMemcacheLockFactoryByClient(ctx, client)
	assert.Nil(t, err)
	assert.NotNil(t, factory)

	// 缓存命中
	factory2, err := GetMemcacheLockFactoryByClient(ctx, client)
	assert.Nil(t, err)
	assert.Equal(t, factory, factory2)

	lock, err := NewMemcacheLockByClient(ctx, client, "test-memcache-lock-client")
	assert.Nil(t, err)
	assert.NotNil(t, lock)

	assert.Nil(t, lock.Lock(ctx, "owner-1"))
	assert.Nil(t, lock.UnLock(ctx, "owner-1"))
}

func TestNewMemcacheLockByClientWithOptions(t *testing.T) {
	client := testMemcacheClient(t)
	defer client.Close()

	options := storage_lock.NewStorageLockOptionsWithLockId("test-memcache-lock-client-opt")
	lock, err := NewMemcacheLockByClientWithOptions(context.Background(), client, options)
	assert.Nil(t, err)
	assert.NotNil(t, lock)
}

func TestMemcacheLockByClient_DeadServer(t *testing.T) {
	// preset client 直返，但 Init 的 Ping 失败
	dead := memcache.New("127.0.0.1:19999")
	defer dead.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := GetMemcacheLockFactoryByClient(ctx, dead)
	assert.NotNil(t, err)
	_, err = NewMemcacheLockByClient(ctx, dead, "x")
	assert.NotNil(t, err)
	_, err = NewMemcacheLockByClientWithOptions(ctx, dead, storage_lock.NewStorageLockOptionsWithLockId("x"))
	assert.NotNil(t, err)
}

func TestNewMemcacheLockByServers_Local(t *testing.T) {
	servers := []string{"127.0.0.1:11211"}
	probe := memcache.New(servers...)
	_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := probe.Ping(); err != nil {
		t.Skipf("memcache 127.0.0.1:11211 不可用，跳过: %v", err)
	}
	probe.Close()

	ctx2 := context.Background()
	factory, err := GetMemcacheLockFactoryByServers(ctx2, servers)
	assert.Nil(t, err)
	assert.NotNil(t, factory)

	// 缓存命中（同 serversKey）
	factory2, err := GetMemcacheLockFactoryByServers(ctx2, []string{"127.0.0.1:11211"})
	assert.Nil(t, err)
	assert.Equal(t, factory, factory2)

	lock, err := NewMemcacheLockByServers(ctx2, servers, "test-memcache-lock-servers")
	assert.Nil(t, err)
	assert.NotNil(t, lock)

	options := storage_lock.NewStorageLockOptionsWithLockId("test-memcache-lock-servers-opt")
	lock2, err := NewMemcacheLockByServersWithOptions(ctx2, servers, options)
	assert.Nil(t, err)
	assert.NotNil(t, lock2)

	assert.Nil(t, lock.Lock(ctx2, "owner-1"))
	assert.Nil(t, lock.UnLock(ctx2, "owner-1"))
}

func TestMemcacheLockByServers_DeadServer(t *testing.T) {
	// Take 的 Ping 失败
	dead := []string{"127.0.0.1:19999"}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := GetMemcacheLockFactoryByServers(ctx, dead)
	assert.NotNil(t, err)
	_, err = NewMemcacheLockByServers(ctx, dead, "x")
	assert.NotNil(t, err)
	_, err = NewMemcacheLockByServersWithOptions(ctx, dead, storage_lock.NewStorageLockOptionsWithLockId("x"))
	assert.NotNil(t, err)
}

func TestServersKey(t *testing.T) {
	assert.Equal(t, "a,b", serversKey([]string{"a", "b"}))
	assert.Equal(t, "", serversKey(nil))
}
