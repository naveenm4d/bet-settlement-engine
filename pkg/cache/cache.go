package cache

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/allegro/bigcache"
)

type Cache interface {
	Set(key string, value any) error
	Get(key string) (any, error)
}

type cache struct {
	cache *bigcache.BigCache
}

func NewCache(
	ctx context.Context,
	config bigcache.Config,
) (Cache, error) {
	bc, err := bigcache.NewBigCache(config)
	if err != nil {
		return nil, err
	}

	cacheService := &cache{
		cache: bc,
	}

	return cacheService, nil
}

func (c *cache) Set(key string, value any) error {
	// Serialize the value into bytes.
	valueBytes, err := serialize(value)
	if err != nil {
		return err
	}

	return c.cache.Set(key, valueBytes)
}

func (c *cache) Get(key string) (any, error) {
	// Get the value in the byte format it is stored in.
	valueBytes, err := c.cache.Get(key)
	if err != nil {
		// Entry not found in cache
		return nil, err
	}

	// Deserialize the bytes of the value.
	value, err := deserialize(valueBytes)
	if err != nil {
		return nil, err
	}

	return value, nil
}

// serialize encode the struct to byte.
func serialize(value any) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	gob.Register(value)

	err := enc.Encode(&value)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// deserialize decode the byte to struct.
func deserialize(valueBytes []byte) (any, error) {
	var value interface{}

	buf := bytes.NewBuffer(valueBytes)
	dec := gob.NewDecoder(buf)

	err := dec.Decode(&value)
	if err != nil {
		return nil, err
	}

	return value, nil
}
