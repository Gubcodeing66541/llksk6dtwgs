package Service

type ServiceMessageReplyReq struct {
	Id      int    `json:"id"`
	Title   string `json:"title" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Content string `json:"content" binding:"required"`
	MegType string `json:"meg_type" binding:"required"`
}
