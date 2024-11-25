package Request

import "server/App/Common"

type PageLimit struct {
	Limit  int    `json:"limit"  form:"limit" uri:"limit" xml:"limit"`
	Page   int    `json:"page"  form:"page" uri:"page" xml:"page"`
	Member string `json:"member"  form:"member" uri:"member" xml:"member"`
	Type   string `json:"type"  form:"type" uri:"type" xml:"type"`
	Id     int    `json:"id"  form:"id" uri:"id" xml:"id"`
	Ids    int    `json:"ids"  form:"ids" uri:"ids" xml:"ids"`
}

type Member struct {
	Member string `json:"member"  form:"member" uri:"member" xml:"member"`
}

func (p *PageLimit) Check() {
	p.Limit = Common.Tools{}.If(p.Limit > 500, 500, p.Limit).(int)
	p.Page = Common.Tools{}.If(p.Page < 1, 1, p.Page).(int)

}

func (p *PageLimit) GetOffset() int {
	p.Check()
	return (p.Page - 1) * p.Limit
}

func (p *PageLimit) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 20
	} else {
		p.Limit = Common.Tools{}.If(p.Limit > 500, 500, p.Limit).(int)
	}
	return p.Limit
}
