package orm

import "gorm.io/gorm"

type MQuery struct {
	query *gorm.DB
}

func (p *MQuery)Where(param string, tag string,value interface{}) *MQuery {
	p.query = p.query.Where("%v %v %v",[]interface{}{param,tag,value})
	return p
}

func (p *MQuery)In(param string, value []interface{}) *MQuery {
	p.query = p.query.Where(param,value)
	return p
}

func (p *MQuery)NotIn(param string, value []interface{}) *MQuery {
	p.query = p.query.Not(param,value)
	return p
}

func (p *MQuery)Order(order string) *MQuery {
	p.query = p.query.Order(order)
	return p
}

func (p *MQuery)Limit(limit int) *MQuery {
	p.query = p.query.Limit(limit)
	return p
}

func (p *MQuery)Offset(index int) *MQuery {
	p.query = p.query.Offset(index)
	return p
}

func (p *MQuery)Get(out interface{}) error {
	return p.query.Find(out).Error
}

func (p *MQuery)Count() (int64,error) {
	var count int64
	err := p.query.Count(&count).Error
	return count,err
}

func Query(model interface{}) *MQuery {
	p := new(MQuery)
	if model != nil {
		p.query = DB.Model(model)
	} else {
		p.query = DB
	}
	return p
}
