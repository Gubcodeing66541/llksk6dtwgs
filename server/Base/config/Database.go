package Config

type Database struct {
	Mysql Mysql
	Redis Redis
}

type Mysql struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

type Redis struct {
	Host     string
	Port     int
	Password string
}
