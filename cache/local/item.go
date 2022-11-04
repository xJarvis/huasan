package local

import (
	"time"
)

const (
	// ISEXPIRATION 绝对过期
	ISEXPIRATION 	int64 = -1

	// NONEEXPIRATION 永不过期
	NONEEXPIRATION 	int64 = 0
)

type Item struct {
	object     interface{} 	//真正的数据项
	Expiration int64       		//生存时间
}

func (item *Item) Expired() bool {
	if NONEEXPIRATION == item.Expiration {
		return false	//永不过期
	}
	return time.Now().UnixNano() > item.Expiration
}

func (item *Item) Set(obj interface{},time int64) {
	item.object = obj
	item.Expiration = time
}

func (item *Item) Get() interface{} {
	return item.object
}

func (item *Item) Del() {
	item.object = nil
	item.Expiration = ISEXPIRATION
}
