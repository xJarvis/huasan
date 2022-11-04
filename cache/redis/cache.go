package redis

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/xjarvis/huashan/app/config"
	"strconv"
	"time"
)

const CACHE_TTL  time.Duration = 3600

type Cache interface {
	CacheKey() string
}

func CacheKey(tag string,obj Cache) string {
	return config.Get("redis.prefix") + ":" + obj.CacheKey() + "_" + tag
}

func Set(id int,obj Cache) error {
	if !Open {return errors.New("redis not support!")}
	jObj,err := json.Marshal(obj)
	if err != nil {
		return err
	}
	err = Redis.Set(Redis.Context(), CacheKey(strconv.Itoa(id),obj), jObj, CACHE_TTL* time.Second).Err()
	return err
}

func Get(id int,obj Cache) error {
	if !Open {return errors.New("redis not support!")}
	data,err := Redis.Get(Redis.Context(), CacheKey(strconv.Itoa(id),obj)).Bytes()
	if err != nil {
		return err
	}
	err = json.Unmarshal(data,obj)
	return err
}

func Del(id int,obj Cache) error {
	if !Open {return errors.New("redis not support!")}
	return Redis.Del(Redis.Context(), CacheKey(strconv.Itoa(id),obj)).Err()
}

func Write(tag string,obj Cache) error {
	if !Open {return errors.New("redis not support!")}
	jObj,err := json.Marshal(obj)
	if err != nil {
		return err
	}
	err = Redis.Set(Redis.Context(), CacheKey(tag,obj), jObj, CACHE_TTL* time.Second).Err()
	return err
}

func Read(tag string,obj Cache) error {
	if !Open {return errors.New("redis not support!")}
	data,err := Redis.Get(Redis.Context(), CacheKey(tag,obj)).Bytes()
	if err != nil {
		return err
	}
	err = json.Unmarshal(data,obj)
	return err
}

func Clear(tag string,obj Cache) error {
	if !Open {return errors.New("redis not support!")}
	return Redis.Del(Redis.Context(), CacheKey(tag,obj)).Err()
}
