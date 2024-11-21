package Common

type Rename struct {
	Id     int    `json:"id"  gorm:"primary_key;AUTO_INCREMENT"`
	Rename string `json:"rename" gorm:"CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci"`
}
