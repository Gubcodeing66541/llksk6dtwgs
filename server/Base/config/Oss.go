package Config

type Oss struct {
	Ali Ali `json:"ali"`
}

type Ali struct {
	Region          string `json:"region"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
}
