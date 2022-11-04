package influxdb

import (
	"github.com/xjarvis/huashan/app/config"
	"github.com/xjarvis/huashan/nosql/influxdb/client"
	"github.com/xjarvis/huashan/log/logger"
	"sync"
	"time"
)

const PRECISION 		string = "ns"
const CHANBUFFER 		int = 20480
const MAXPOINTS		int = 256

const WAITSECOND		= 3			//秒

var (
	InfluxConn client.Client
	Influx     InfluxDB
)

func NewInfluxConn() client.Client {
	conn, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.Get("influx.url"),
		Username: config.Get("influx.username"),
		Password: config.Get("influx.password"),
	})
	if err != nil {
		panic(err)
	}
	return conn
}

func Initialize() {
	InfluxConn = NewInfluxConn()
	// 检查连接
	_, _, err := InfluxConn.Ping(WAITSECOND * time.Second)
	if err != nil {
		panic(err)
	}

	// 初始化
	Influx = InfluxDB{}
	if err := Influx.Init(config.Get("influx.database")); err != nil {
		panic(err)
	}

	logger.Info("influx db init finish")
}

func UnInitialize() {
	Influx.Release()
	InfluxConn.Close()
	logger.Info("influx db release finish")
}

func Run() {
	go Influx.Run()
}


type PointData struct {
	Measurement string
	Tags 		map[string]string
	Fields 	 	map[string]interface{}
}

type InfluxDB struct{
	InfluxBP     	client.BatchPoints
	PointChan		chan PointData

	pointLen 		int
	pointSyn        sync.Mutex
	waitTime		int64
}

func (i *InfluxDB) Init (database string) error {
	var err error
	i.InfluxBP,err = client.NewBatchPoints(client.BatchPointsConfig{
		Database:  database,
		Precision: PRECISION,
	})
	if err != nil {
		return err
	}
	i.PointChan = make(chan PointData, CHANBUFFER)
	i.pointLen  = 0
	return nil
}

func (i *InfluxDB) Add (data PointData) error {
	i.pointSyn.Lock()
	defer i.pointSyn.Unlock()

	pt, err := client.NewPoint(data.Measurement, data.Tags, data.Fields, time.Now())
	if err != nil {
		logger.Error("influx add point failed!",err)
		//return err
	} else {
		logger.Debug("point :",pt.String())
	}
	i.InfluxBP.AddPoint(pt)
	i.pointLen++
	return nil
}

func (i *InfluxDB) Save () error {
	i.pointSyn.Lock()
	defer i.pointSyn.Unlock()
	defer func() {
		i.InfluxBP.ClearPoints()
		i.pointLen = 0
	}()

	tf := time.Now()
	if err := InfluxConn.Write(i.InfluxBP); err != nil {
		logger.Error("influx save bp failed!",err)
		return err
	}
	logger.Debug("influx db save time:",time.Now().Sub(tf))
	logger.Debug("influx db save count:",i.pointLen)
	i.waitTime = time.Now().Add(WAITSECOND * time.Second).UnixNano()

	return nil
}

func (i *InfluxDB) IsSave() bool {
	if i.pointLen >= MAXPOINTS {
		return true
	}
	if i.pointLen > 0 && time.Now().UnixNano() > i.waitTime {
		return true
	}
	return false
}

func (i *InfluxDB) Receive(data PointData) {
	i.PointChan <- data
}

func (i *InfluxDB) Release() {
	for i.pointLen > 0 {
		time.Sleep(WAITSECOND * time.Second)
	}
}

func (i *InfluxDB) Run() {
	go func() {
		for {
			if i.IsSave() {
				i.Save()
			} else {
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	go func() {
		for {
			data := <- i.PointChan
			i.Add(data)
		}
	}()
}




