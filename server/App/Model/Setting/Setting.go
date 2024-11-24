package Setting

import "time"

type ServiceRoomDetail struct {
	Id         int    `gorm:"primary_key;AUTO_INCREMENT"`
	ServiceId  int    `gorm:"index:service_id_idx"`
	Scan       string `json:"scan"`        //扫码引粉  enable 启用、  un_enable 不启用 、  stop 暂停引入新粉
	ScanDrive  string `json:"scan_drive"`  //扫码用户  wx 仅限微信	 、default 不限制	、ie 仅限浏览器
	ScanFitter string `json:"scan_fitter"` //扫码过滤 un_enable 不启用 、 room  过滤机房、not_china 过滤非大陆  、room_and_not_china  过滤 机房 和非大陆

	ScanChange string `json:"scan_change"` //扫码用户验证    enable 启用、  un_enable 不启用 、

	ToUrl     string `json:"to_url"`      // 网址跳转   enable 启用、  un_enable 不启用
	Banner    string `json:"banner"`      // 滚动条
	ScanToUrl string `json:"scan_to_url"` // 扫码跳转下载页 必须http起头

	ChangeQrSound string `json:"change_qr"`    // 扫码提醒 enable 启用、  un_enable 不启用
	MessageSound  string `json:"message"`      // 消息提示音 enable 启用、  un_enable 不启用
	OnlineSound   string `json:"online_sound"` // 上线提示音 enable 启用、  un_enable 不启用
	CreateTime    time.Time
}
