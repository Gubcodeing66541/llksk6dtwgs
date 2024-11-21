package Admin

type ServiceCreateReq struct {
	Number int    `json:"number"  form:"number" uri:"number" xml:"number" binding:"required"`
	Type   string `json:"type"  form:"type" uri:"type" xml:"type" binding:"required"`
	Day    int    `json:"day"  form:"day" uri:"day" xml:"day"`
	Price  int    `json:"price"  form:"price" uri:"price" xml:"price" binding:"required"`
}

type ServiceRenewReq struct {
	Member []string `json:"member"  form:"member" uri:"member" xml:"member" binding:"required"`
	Day    int      `json:"day"  form:"day" uri:"day" xml:"day" binding:"required"`
	Type   string   `json:"type"  form:"type" uri:"type" xml:"type" binding:"required"`
	Price  int      `json:"price"  form:"price" uri:"price" xml:"price" binding:"required"`
}
