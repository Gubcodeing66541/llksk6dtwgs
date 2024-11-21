package User

type UserAuthMap struct {
	CookieUid string `json:"cookie_uid"`
	UserId    int    `json:"user_id"`
	Action    string `json:"action"`
}
