package caching

import(
	"time"
	"context"
	
	gocache "github.com/patrickmn/go-cache"
	libcache "github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/store/go_cache/v4"
)

type Cacher struct {
	ctx context.Context
	store *go_cache.GoCacheStore
	manager *libcache.Cache[[]byte]
}

func NewCacher() *Cacher {
	cacher := &Cacher{
		ctx:context.Background(),
	}
	client := gocache.New(5 * time.Minute, 10 * time.Minute)
	cacher.store = go_cache.NewGoCache(client)
	cacher.manager = libcache.New[[]byte](cacher.store)
	
	return cacher
}

func (c *Cacher) Get(key string) (string, error) {
	value, err := c.manager.Get(c.ctx, key)
	if err != nil {
		return "", err
	}
	
	return string(value[:]), nil	
}

func (c *Cacher) Set(key string, value string, expires time.Duration) error {
	err := c.manager.Set(c.ctx, key, []byte(value))
	if err != nil {
		return err
	}
	return nil
}

func (c *Cacher) Remove(key string) error {
	err := c.manager.Delete(c.ctx, key)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cacher) Clear() error {
	err := c.manager.Clear(c.ctx)
	if err != nil {
		return err
	}
	return nil
}	