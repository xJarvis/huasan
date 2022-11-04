package cron

var schedule_id int = 10000

// 任务基类
type Schedule struct {
	ID 				int				//任务ID
	Name 			string			//任务名称
	Spec 			string			//任务运行周期
}

func (p *Schedule)Initialize(mData *map[string]interface{}) error {
	// 计划ID自增
	schedule_id++

	// 任务ID
	p.ID = schedule_id

	return nil
}

func (p *Schedule)Update(mData *map[string]interface{}) error {
	return nil
}

// 任务ID
func (p *Schedule)GetID() int {
	return p.ID
}

// 任务名称
func (p *Schedule)GetName() string {
	return p.Name
}

// 任务计划
func (p *Schedule)GetSpec() string	{
	return p.Spec
}

// 活动状态
func (p *Schedule)IsActive() bool {
	return true
}
// 执行
func (p *Schedule)Run() error {
	return nil
}

// 环境释放
func (p *Schedule)Release() error {
	return nil
}