package Admin

type LoginReq struct {
	Member   string `json:"member"  form:"member" uri:"member" xml:"member" binding:"required"`
	Password string `json:"password"  form:"password" uri:"password" xml:"password" binding:"required"`
}
