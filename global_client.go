package memcache_locks

import (
	"context"
	memcache_storage "github.com/storage-lock/go-memcache-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
	storage_lock_factory "github.com/storage-lock/go-storage-lock-factory"
	"github.com/bradfitz/gomemcache/memcache"
)

var clientStorageLockFactoryBeanFactory *storage_lock_factory.StorageLockFactoryBeanFactory[*memcache.Client, *memcache.Client] = storage_lock_factory.NewStorageLockFactoryBeanFactory[*memcache.Client, *memcache.Client]()

func NewMemcacheLockByClient(ctx context.Context, client *memcache.Client, lockId string) (*storage_lock.StorageLock, error) {
	factory, err := GetMemcacheLockFactoryByClient(ctx, client)
	if err != nil {
		return nil, err
	}
	return factory.CreateLock(lockId)
}

func NewMemcacheLockByClientWithOptions(ctx context.Context, client *memcache.Client, options *storage_lock.StorageLockOptions) (*storage_lock.StorageLock, error) {
	factory, err := GetMemcacheLockFactoryByClient(ctx, client)
	if err != nil {
		return nil, err
	}
	return factory.CreateLockWithOptions(options)
}

func GetMemcacheLockFactoryByClient(ctx context.Context, client *memcache.Client) (*storage_lock_factory.StorageLockFactory[*memcache.Client], error) {
	return clientStorageLockFactoryBeanFactory.GetOrInit(ctx, client, func(ctx context.Context) (*storage_lock_factory.StorageLockFactory[*memcache.Client], error) {
		connectionManager := memcache_storage.NewMemcacheConnectionManagerFromClient(client)
		options := memcache_storage.NewMemcacheStorageOptions().SetConnectionManager(connectionManager)
		storage, err := memcache_storage.NewMemcacheStorage(ctx, options)
		if err != nil {
			return nil, err
		}
		factory := storage_lock_factory.NewStorageLockFactory(storage, options.ConnectionManager)
		return factory, nil
	})
}
