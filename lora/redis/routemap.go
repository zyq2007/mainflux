//
// Copyright (c) 2018
// Mainflux
//
// SPDX-License-Identifier: Apache-2.0
//

package redis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/mainflux/mainflux/lora"
)

const (
	loraMapPrefix = "lora:mfx"
	mfxMapPrefix  = "mfx:lora"
)

var _ lora.RouteMapRepository = (*routerMap)(nil)

type routerMap struct {
	client *redis.Client
}

// NewRouteMapRepository returns redis thing cache implementation.
func NewRouteMapRepository(client *redis.Client) lora.RouteMapRepository {
	return &routerMap{
		client: client,
	}
}

func (mr *routerMap) Save(key string, val string) error {
	println("ROUTEMAP SAVE")
	tkey := fmt.Sprintf("%s:%s", loraMapPrefix, key)
	if err := mr.client.Set(tkey, val, 0).Err(); err != nil {
		return err
	}
	tkey = fmt.Sprintf("%s:%s", mfxMapPrefix, val)
	if err := mr.client.Set(tkey, key, 0).Err(); err != nil {
		return err
	}

	val, err := mr.client.Get(tkey).Result()
	if err != nil {
		return err
	}
	println("ROUTEMAP SAVED: ", tkey, val)

	return nil
}

func (mr *routerMap) Get(key string) (string, error) {
	laKey := fmt.Sprintf("%s:%s", loraMapPrefix, key)
	val, err := mr.client.Get(laKey).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}

func (mr *routerMap) Remove(key string) error {
	tid := fmt.Sprintf("%s:%s", loraMapPrefix, key)
	key, err := mr.client.Get(tid).Result()
	if err != nil {
		return err
	}

	tkey := fmt.Sprintf("%s:%s", loraMapPrefix, key)
	return mr.client.Del(tkey, tid).Err()
}