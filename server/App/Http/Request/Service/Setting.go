package Service

import "github.com/go-playground/validator/v10"

type BindDiversionReq struct {
	Member string `json:"member" form:"member"`
}

func (BindDiversionReq) GetError(err error) string {
	_, ok := err.(validator.ValidationErrors)
	if !ok {
		return "参数错误"
	}
	for _, v := range err.(validator.ValidationErrors) {
		switch v.Field() {
		case "Member":
			switch v.Tag() {
			case "required":
				return "账号不能为空"
			}
		}
	}

	return "参数错误"
}
