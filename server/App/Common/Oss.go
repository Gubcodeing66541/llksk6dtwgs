package Common

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"hash"
	"io"
)

type Oss struct{}

// 获取阿里云oss
func (Oss) GetAliToken(OssAliRegion string, OssAliAccessKeyId string, OssAliAccessKeySecret string) (string, error) {
	//构建一个阿里云客户端, 用于发起请求。
	//设置调用者（RAM用户或RAM角色）的AccessKey ID和AccessKey Secret。
	//第一个参数就是bucket所在位置，可查看oss对象储存控制台的概况获取
	//第二个参数就是步骤一获取的AccessKey ID
	//第三个参数就是步骤一获取的AccessKey Secret
	client, err := sts.NewClientWithAccessKey(OssAliRegion, OssAliAccessKeyId, OssAliAccessKeySecret)

	//构建请求对象。
	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"

	//发起请求，并得到响应。
	response, err := client.AssumeRole(request)
	if err != nil {
		fmt.Print(err.Error())
	}

	// 结构体
	type Policy struct {
		Expiration string          `json:"expiration"`
		Conditions [][]interface{} `json:"conditions"`
	}
	// 生成签名代码
	var policy Policy
	policy.Expiration = "9999-12-31T12:00:00.000Z"
	var conditions []interface{}
	conditions = append(conditions, "content-length-range")
	conditions = append(conditions, 0)
	conditions = append(conditions, 1048576000)
	policy.Conditions = append(policy.Conditions, conditions)
	policyByte, err := json.Marshal(policy)
	if err != nil {
		return "", errors.New("序列化失败")
	}
	policyBase64 := base64.StdEncoding.EncodeToString(policyByte)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(response.Credentials.AccessKeySecret))
	io.WriteString(h, policyBase64)
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	// 主要拿的就是下面这两个玩意
	fmt.Println("policyBase64：", policyBase64)
	fmt.Println("signature", signature)
	// 将这两个与前面获取的临时授权参数一起返回就好了
	return signature, errors.New("序列化失败")

}
