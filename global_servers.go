package memcache_locks

import (
	"context"
	memcache_storage "github.com/storage-lock/go-memcache-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
	storage_lock_factory "github.com/storage-lock/go-storage-lock-factory"
	"github.com/bradfitz/gomemcache/memcache"
	"strings"
)

var serversStorageLockFactoryBeanFactory *storage_lock_factory.StorageLockFactoryBeanFactory[string, *memcache.Client] = storage_lock_factory.NewStorageLockFactoryBeanFactory[string, *memcache.Client]()

func NewMemcacheLockByServers(ctx context.Context, servers []string, lockId string) (*storage_lock.StorageLock, error) {
	factory, err := GetMemcacheLockFactoryByServers(ctx, servers)
	if err != nil {
		return nil, err
	}
	return factory.CreateLock(lockId)
}

func NewMemcacheLockByServersWithOptions(ctx context.Context, servers []string, options *storage_lock.StorageLockOptions) (*storage_lock.StorageLock, error) {
	factory, err := GetMemcacheLockFactoryByServers(ctx, servers)
	if err != nil {
		return nil, err
	}
	return factory.CreateLockWithOptions(options)
}

// serversKey 将服务器列表转为缓存 key
func serversKey(servers []string) string {
	return strings.Join(servers, ",")
}

func GetMemcacheLockFactoryByServers(ctx context.Context, servers []string) (*storage_lock_factory.StorageLockFactory[*memcache.Client], error) {
	key := serversKey(servers)
	return serversStorageLockFactoryBeanFactory.GetOrInit(ctx, key, func(ctx context.Context) (*storage_lock_factory.StorageLockFactory[*memcache.Client], error) {
		connectionManager := memcache_storage.NewMemcacheConnectionManagerFromServers(servers)
		options := memcache_storage.NewMemcacheStorageOptions().SetConnectionManager(connectionManager)
		storage, err := memcache_storage.NewMemcacheStorage(ctx, options)
		if err != nil {
			return nil, err
		}
		factory := storage_lock_factory.NewStorageLockFactory(storage, options.ConnectionManager)
		return factory, nil
	})
}
