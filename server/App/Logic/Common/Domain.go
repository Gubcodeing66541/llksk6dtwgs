package Common

import (
	"errors"
	"server/App/Model/Common"
	Service2 "server/App/Model/Service"
	"server/Base"
	"strings"
	"time"
)

type Domain struct{}

func (Domain) GetAction() string {
	var domain Common.Domain
	Base.MysqlConn.Find(&domain, "type = ? and status = 'enable'", "action")
	if domain.Domain != "" {
		return domain.Domain
	}
	return ""
}
func (Domain) GetLive() string {
	var domain Common.Domain
	Base.MysqlConn.Find(&domain, "type = ? and status = 'enable'", "live")
	if domain.Domain != "" {
		return domain.Domain
	}
	return ""
}

func (Domain) Get(id int) Common.Domain {
	var Model Common.Domain
	Base.MysqlConn.Find(&Model, "id = ?", id)
	return Model
}

func (Domain) GetNoUsePrivateNum() int {
	var leng int
	Base.MysqlConn.Model(&Common.Domain{}).Where("type = 'private' and bind_service_id = 0 and we_chat_ban_status = 'success' and status = 'enable'").Count(&leng)
	return leng
}

func (Domain) GetTransfer() Common.Domain {
	var domain Common.Domain
	Base.MysqlConn.Find(&domain, "type = 'transfer' and status = 'enable' and we_chat_ban_status = 'success' ")
	return domain
}

func (Domain) GetServiceBind(serviceId int) Common.Domain {
	var domain Common.Domain
	Base.MysqlConn.Find(&domain, "bind_service_id = ?", serviceId)
	return domain
}

func (Domain) Bind(serviceId int) error {
	err := ServiceIsOverdue(serviceId)
	if err != nil {
		return err
	}

	var service Service2.Service
	Base.MysqlConn.Find(&service, "service_id = ?", serviceId)

	var domain Common.Domain
	Base.MysqlConn.Find(&domain, "bind_service_id = ?  ", serviceId)
	if domain.BindServiceId != 0 {
		return errors.New("已绑定域名")
	}
	Base.MysqlConn.Find(&domain, "status = ? and  type = ?  and bind_service_id = 0", "enable", "public")
	if domain.Id == 0 {
		return errors.New("无可用分配域名")
	}

	Base.MysqlConn.Model(&domain).Where("id = ?", domain.Id).Update("bind_service_id", serviceId)
	return nil
}

func (Domain) UnEnable(domainId int) {
	var domain Common.Domain
	Base.MysqlConn.Find(&domain, "id = ?", domainId)

	if domain.Id == 0 {
		return
	}

	domain.Status = "un_enable"
	domain.BindServiceId = 0
	Base.MysqlConn.Update(&domain)
}

func (Domain) QueryById(id int) Common.Domain {
	var domain Common.Domain
	Base.MysqlConn.Find(&domain, "id = ?", id)
	return domain
}

func (Domain) Update(id int, domain string, typeEd string, status string) error {
	var domainEntity Common.Domain
	Base.MysqlConn.Find(&domainEntity, "id = ?", id)
	if domainEntity.BindServiceId != 0 {
		return errors.New("该域名已绑定客服，请解绑后再修改")
	}

	Base.MysqlConn.Model(&domainEntity).Where("id = ?", id).Updates(Common.Domain{Domain: domain, Type: typeEd, Status: status})
	return nil
}

func (Domain) Delete(id int) error {
	if id == 0 {
		return errors.New("参数不全")
	}

	var domain Common.Domain
	Base.MysqlConn.Find(&domain, "id = ?", id)

	bindDomain := domain.BindServiceId
	Base.MysqlConn.Delete(&domain, "id = ?", id)
	if bindDomain != 0 {
		err := Domain{}.Bind(bindDomain)
		return err
	}
	return nil
}

func (Domain) Create(domain string, typeEd string, status string) {
	domainList := strings.Split(domain, "\n")
	for _, domain = range domainList {
		var domaineEntity Common.Domain
		domaineEntity.Domain = strings.TrimSpace(domain)
		domaineEntity.CreateTime = time.Now()
		domaineEntity.UpdateTime = time.Now()
		domaineEntity.Status = status
		domaineEntity.WeChatBanStatus = "success"
		if typeEd != "" {
			domaineEntity.Type = typeEd
		} else {
			domaineEntity.Type = "private"
		}
		Base.MysqlConn.Create(&domaineEntity)
	}
}

func ServiceIsOverdue(serviceId int) error {
	var seviceAuth Service2.Service
	Base.MysqlConn.Find(&seviceAuth, "service_id = ?", serviceId)
	if seviceAuth.ServiceId == 0 {
		return errors.New("客服不存在")
	}

	now := time.Now()
	if seviceAuth.TimeOut.Before(now) {
		return errors.New("客服过期")
	}
	return nil
}

func (Domain) EnableDisable(id int, status string) {
	var domain Common.Domain
	Base.MysqlConn.Find(&domain, "id = ?", id)

	if domain.BindServiceId > 0 && status == "un_enable" {
		bindServiceId := domain.BindServiceId

		Base.MysqlConn.Model(&domain).Where("id = ?", domain.Id).Updates(&Common.Domain{Status: "un_enable", BindServiceId: 0})

		var domainStruct Domain
		domainStruct.Bind(bindServiceId)
	} else {
		Base.MysqlConn.Model(&domain).Where("id = ?", id).Update(&Common.Domain{Status: "un_enable", BindServiceId: 0})
	}
}

func (d Domain) CheckBindDomain(service Service2.Service) {
	// 检测账号是否过期
	if service.TimeOut.Unix()-time.Now().Unix() <= 0 {
		return
	}

	// 检测域名是否绑定
	Domain{}.Bind(service.ServiceId)
}
