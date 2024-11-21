package Config

type App struct {
	Model         string `json:"model"`
	Debug         bool   `json:"debug"`
	Port          int    `json:"port"`
	IpRegistryKey string `json:"ip_registry_key"`
	Database      Database
	Client        Client
	Oss           Oss
	Mq            Mq
	HeadImgUrl    string    `json:"head_img_url"`
	Manager       Manager   `json:"manager"`
	HttpHost      string    `json:"http_host"`
	LiveAppKey    string    `json:"live_app_key"`
	LiveAppSecret string    `json:"live_app_secret"`
	LiveAppHost   string    `json:"live_app_host"`
	AesKey        string    `json:"aes_key"`
	PayUrl        string    `json:"pay_url"`
	PayConfig     PayConfig `json:"pay_config"`
}
