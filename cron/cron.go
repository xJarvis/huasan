package cron

import (
	"huashan/exerror"
	"huashan/logger"
	"fmt"
	"github.com/jakecoffman/cron"
	"strconv"
	"sync"
)

// 定时任务调度管理器
var scheduled *cron.Cron

// 任务列表
var TaskMap sync.Map

/*任务接口*/
type ITask interface {
	Initialize(pData *map[string]interface{}) error
	Update(pData *map[string]interface{}) error
	GetID() 	int
	GetName() 	string
	GetSpec() 	string
	Run() 		error
	IsActive()	bool
	Release() 	error
}

// 初始化任务
func Initialize() {
	//初始化任务调度器
	scheduled = cron.New()
	scheduled.Start()
}

// 结束计划任务
func UnInitialize() {
	exerror.Catch()
	ClearJob()
	if scheduled != nil {
		scheduled.Stop()
	}
	logger.Info("schedule task release finish!")
}


// 创建计划任务
func createJob(iTask ITask) (int,string,cron.FuncJob) {
	logger.Debug(fmt.Sprintf("add task(%v-%v)to schedule,spec:%v",iTask.GetID(),iTask.GetName(),iTask.GetSpec()))
	taskID := iTask.GetID()
	taskSpec := iTask.GetSpec()
	taskFunc := func() {
		defer exerror.Catch() //捕获异常
		logger.Debug(fmt.Sprintf("task begin:(%v-%v)", iTask.GetID(), iTask.GetName()))
		err := iTask.Run()
		if err != nil {
			logger.Error(fmt.Sprintf("task running failed:(%v-%v)", iTask.GetID(), iTask.GetName()))
		} else {
			logger.Debug(fmt.Sprintf("task   end:(%v-%v)", iTask.GetID(), iTask.GetName()))
			logger.Debug("***********************  next ************************")
		}
	}
	return taskID,taskSpec,taskFunc
}

// 添加计划任务
func AddJob(iTask ITask) {
	defer exerror.Catch() //捕获异常

	// 添加计划任务
	taskID,taskSpec,taskFunc := createJob(iTask)
	scheduled.AddFunc(taskSpec, taskFunc, strconv.Itoa(taskID))

	// 存储已加载任务
	TaskMap.Store(iTask.GetID(),iTask)
}

// 删除计划任务
func RemoveJob(taskID int) {
	scheduled.RemoveJob(strconv.Itoa(taskID))
	if value,ok :=TaskMap.Load(taskID); ok {
		iTask := value.(ITask)
		iTask.Release()
		TaskMap.Delete(taskID)
	} else {
		logger.Error(fmt.Sprintf("job(%v) already exist!",taskID))
	}
}

/*清理非活动任务*/
func ClearUnActiveJob() {
	exerror.Catch()
	cFunc := func(key interface{},value interface{})bool {
		iTask := value.(ITask)
		if !iTask.IsActive() {
			scheduled.RemoveJob(strconv.Itoa(key.(int)))
			iTask.Release()
			TaskMap.Delete(key.(int))
		}
		return true
	}
	//遍历MAP
	TaskMap.Range(cFunc)
}

/*清理job*/
func ClearJob() {
	exerror.Catch()
	cFunc := func(key interface{},value interface{})bool {
		scheduled.RemoveJob(strconv.Itoa(key.(int)))
		iTask := value.(ITask)
		iTask.Release()
		TaskMap.Delete(key.(int))
		return true
	}

	//遍历MAP
	TaskMap.Range(cFunc)
}
