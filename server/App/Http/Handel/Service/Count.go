package Service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"server/App/Common"
	"server/Base"
)

type Count struct{}

func (Count) Count(c *gin.Context) {
	member := c.Param("code")

	sql := `
		select count(*) as ip_cnt,sum(cnt) as user_cnt,times from (
    select count(*) as cnt,ip, DATE_FORMAT(create_time,'%Y-%m-%d') as times from service_room_details where service_id in (
        select service_id from services where code = ?
    ) group by ip,times
) t group by times order by  times desc
	`

	type response struct {
		IpCnt   int64  `json:"ip_cnt"`
		UserCnt int64  `json:"user_cnt"`
		Times   string `json:"times"`
	}

	var res []response
	Base.MysqlConn.Raw(sql, member).Scan(&res)

	resStr, err := json.Marshal(&res)
	if err != nil {
		Common.ApiResponse{}.Error(c, err.Error(), gin.H{})
		return
	}

	c.HTML(http.StatusOK, "count.html", gin.H{"count": string(resStr)})
}
