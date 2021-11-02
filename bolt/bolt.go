package bolt

import (
	"fmt"
	"github.com/boltdb/bolt"
)

var (
	BDB     *bolt.DB
	Bucket []byte
)


const dbname = "module.db"

func init() {

	NewBolt()
	//MGO = NewMgo()
}



func NewBolt(){
	dbc, err := bolt.Open(dbname, 0600, nil)
	//初始化bucket
	Bucket = []byte("demoBucket")
	if err != nil {

		panic(fmt.Errorf("connect bolt is error : %v",err))
	} else {

		BDB = dbc
	}

	BDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket(Bucket)
		return err
	})

	//创建bucket

}

//把数据插入到bolt数据库中，相当于redis中的set命令
func Insert(key, value string) {
	k := []byte(key)
	v := []byte(value)
	fmt.Println(BDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(Bucket)

		err := b.Put(k, v)
		return err
	}))
}

// Rm 删除一个指定的key中的数据
func Rm(key string) {
	k := []byte(key)
	BDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(Bucket)
		err := b.Delete(k)
		return err
	})
}

//读取一条数据
func Read(key string) string {
	k := []byte(key)
	var val []byte
	BDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(Bucket)
		val = b.Get(k)
		return nil
	})
	return string(val)
}



