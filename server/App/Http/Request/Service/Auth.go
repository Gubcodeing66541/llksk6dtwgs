package Service

type Login struct {
	Member string `json:"member"  form:"member" uri:"member" xml:"member" binding:"required"`
}
