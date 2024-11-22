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

func (Count) CountWeek(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"count": string(resStr)})
}

func (Count) CountMonth(c *gin.Context) {
	member := c.Param("code")

	sql := `
        SELECT count(*) AS ip_cnt, sum(cnt) AS user_cnt, times 
        FROM (
            SELECT count(*) AS cnt, ip, DATE_FORMAT(create_time, '%Y-%m-%d') AS times 
            FROM service_room_details 
            WHERE service_id IN (
                SELECT service_id FROM services WHERE code = ?
            ) 
            AND create_time >= DATE_FORMAT(CURDATE(), '%Y-%m-01')  -- 当月第一天
            AND create_time < DATE_FORMAT(CURDATE() + INTERVAL 1 MONTH, '%Y-%m-01')  -- 下月第一天
            GROUP BY ip, times
        ) t 
        GROUP BY times 
        ORDER BY times DESC
    `

	type response struct {
		IpCnt   int64  `json:"ip_cnt"`
		UserCnt int64  `json:"user_cnt"`
		Times   string `json:"times"`
	}

	var res []response
	// 执行 SQL 查询
	Base.MysqlConn.Raw(sql, member).Scan(&res)

	// 将查询结果转换为 JSON 格式
	resStr, err := json.Marshal(&res)
	if err != nil {
		Common.ApiResponse{}.Error(c, err.Error(), gin.H{})
		return
	}

	// 返回查询结果
	c.JSON(http.StatusOK, gin.H{"count": string(resStr)})
}

func (Count) CountDay(c *gin.Context) {
	member := c.Param("code")

	sql := `
        SELECT count(*) AS ip_cnt, sum(cnt) AS user_cnt, times 
        FROM (
            SELECT count(*) AS cnt, ip, DATE_FORMAT(create_time, '%Y-%m-%d') AS times 
            FROM service_room_details 
            WHERE service_id IN (
                SELECT service_id FROM services WHERE code = ?
            ) 
            AND DATE(create_time) = CURDATE()  -- 当前日期
            GROUP BY ip, times
        ) t 
        GROUP BY times 
        ORDER BY times DESC
    `

	type response struct {
		IpCnt   int64  `json:"ip_cnt"`
		UserCnt int64  `json:"user_cnt"`
		Times   string `json:"times"`
	}

	var res []response
	// 执行 SQL 查询
	Base.MysqlConn.Raw(sql, member).Scan(&res)

	// 将查询结果转换为 JSON 格式
	resStr, err := json.Marshal(&res)
	if err != nil {
		Common.ApiResponse{}.Error(c, err.Error(), gin.H{})
		return
	}

	// 返回查询结果
	c.JSON(http.StatusOK, gin.H{"count": string(resStr)})
}
