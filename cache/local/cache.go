package local

import (
	"fmt"
	"sync"
	"time"
)

const DefaultExpired = 0

type Cache struct {
	DefaultExpired     	time.Duration  //默认存儲过期时间
	items             	map[string]Item //缓存数据项存储在map中
	gcInterval        	time.Duration  //过期数据项清理周期
	stopGc           	chan bool
	mapmu         	   	sync.RWMutex    		//读写锁
}

// 删除
func (c *Cache) del(k string) {
	c.mapmu.Lock()
	defer c.mapmu.Unlock()
    if item, ok := c.items[k]; ok {
        item.Del()
    }
    delete(c.items, k)
}


//设置
func (c *Cache) set(k string, v interface{}, d time.Duration) {
	if !IsOpen {return}
	c.mapmu.Lock()
	defer c.mapmu.Unlock()
	item,ok := c.items[k]
	if ok {
		item.Del()
	}
    if d != DefaultExpired {
        e := time.Now().Add(d).UnixNano()
        item.Set(v,e)
    } else {
    	e := NONEEXPIRATION
    	if e != int64(c.DefaultExpired) {
			e = time.Now().Add(c.DefaultExpired).UnixNano()
		}
		item.Set(v,e)
    }
	c.items[k] = item
}

//读取
func (c *Cache) get(k string) (interface{}, error) {
	if !IsOpen {return nil,fmt.Errorf("cache not open")}
	c.mapmu.RLock()
	defer c.mapmu.RUnlock()
    item, found := c.items[k]
    if !found {
        return nil, fmt.Errorf("item not exist")
    }
    if item.Expired() {
        return nil, fmt.Errorf("item is expired")
    }

    return item.Get(), nil
}

//判断
func (c *Cache) exist(k string) bool {
	if !IsOpen {return false}
	c.mapmu.RLock()
	defer c.mapmu.RUnlock()
	item,exist := c.items[k]
	if !exist {
		return false
	}
	if item.Expired() {
		return false
	}
	return true
}

// 刷新
func (c *Cache) flush() {
	c.mapmu.Lock()
	defer c.mapmu.Unlock()
	for k,item := range c.items {
		item.Del()
		delete(c.items, k)
	}
	c.items = make(map[string]Item)
}

// 存储数量
func (c *Cache) len() int {
	c.mapmu.RLock()
	defer c.mapmu.RUnlock()
	return len(c.items)
}

//读取所有的Key
//func (c *Cache) keys() []string {
//	c.mapmu.RLock()
//	defer c.mapmu.RUnlock()
//	var keys []string
//	for k,_ := range c.items {
//		keys = append(keys,k)
//	}
//	return keys
//}

// Set 若存在则直接修改，若不存在则增加
func (c *Cache) Set(k string, v interface{}, d time.Duration) {
	c.set(k,v,d)
}

// Add 若不存在则增加，若存在则报错
func (c *Cache) Add(k string, v interface{}, d time.Duration) error {
	if c.exist(k) {
		return fmt.Errorf("item already exist")
	}
	c.set(k, v, d)
	return nil
}

// Get 若存在则返回 v,true， 若不存在则返回 n,false
func (c *Cache) Get(k string) (interface{}, bool) {
	obj,err := c.get(k)
	if err != nil {
		return nil,false
	} else {
		return obj,true
	}
}

// 若存在则返回 v,nil, 若不存在则返回 n,err
func (c *Cache) Read(k string) (interface{}, error) {
	return c.get(k)
}

func (c *Cache) Delete(k string) {
	c.del(k)
}

func (c *Cache) Exit(k string) bool {
	return c.exist(k)
}

func (c *Cache) Replace(k string, v interface{}, d time.Duration) error {
	if _, err := c.get(k); err != nil {
		return err
	}
	c.set(k, v, d)
	return nil
}



//func (c *Cache) Save(w io.Writer) (err error) {
//	c.mapmu.RLock()
//	defer c.mapmu.RUnlock()
//
//	enc := gob.NewEncoder(w)
//	defer func() {
//		if x := recover(); x != nil {
//			err = fmt.Errorf("Error registering item types with Gob library")
//		}
//	}()
//
//	for _, v := range c.items {
//		gob.Register(v.Get())
//	}
//	err = enc.Encode(&c.items)
//	return
//}

//func (c *Cache) SaveToFile(file string) error {
//	f, err := os.Create(file)
//	if err != nil {
//		return err
//	}
//	if err = c.Save(f); err != nil {
//		f.Close()
//		return err
//	}
//	return f.Close()
//}

//func (c *Cache) Load(r io.Reader) error {
//	c.mapmu.Lock()
//	defer c.mapmu.Unlock()
//
//	dec := gob.NewDecoder(r)
//	items := map[string]Item{}
//	err := dec.Decode(&items)
//	if err == nil {
//		for k, v := range items {
//			if item,ok := c.items[k]; ok {
//				item.Set(v.Get(),v.Expiration)
//			} else {
//                c.items[k] = v
//            }
//		}
//	}
//	return err
//}

//func (c *Cache) LoadFile(file string) error {
//	f, err := os.Open(file)
//	if err != nil {
//		return err
//	}
//	if err = c.Load(f); err != nil {
//		f.Close()
//		return err
//	}
//	return f.Close()
//}

func (c *Cache) Count() int {
	return c.len()
}

func (c *Cache) Flush() {
	c.flush()
}

func (c *Cache) StopGc() {
	c.stopGc <- true
}

// GcLoop 过期缓存数据项清理
func (c *Cache) GcLoop() {
    ticker := time.NewTicker(c.gcInterval)
    for {
        select {
        case <-ticker.C:
            c.DeleteExpired()
        case <-c.stopGc:
            ticker.Stop()
            return
        }
    }
}

func (c *Cache) DeleteExpired() {
	c.mapmu.Lock()
	defer c.mapmu.Unlock()
    for k, item := range c.items {
    	if item.Expired() {
			item.Del()
			delete(c.items, k)
		}
    }
}


var (
	XCache *Cache
	IsOpen 			bool 	= false
)

// NewCache 初始化一个新的cache服务 公开出去
func NewCache(expired time.Duration) *Cache {
	c := &Cache{
		DefaultExpired:    expired,
		gcInterval:        5 * time.Second,		// 5秒进行一次清理
		items:             make(map[string]Item),
		stopGc:            make(chan bool),
	}
	go c.GcLoop()
	return c
}

func Initialize(expired time.Duration) {
    XCache = NewCache(expired)
	IsOpen = true
}
func UnInitialize() {
	IsOpen = false
	XCache.StopGc()
}