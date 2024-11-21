package Common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var host = "https://api.ipregistry.co"

type IpRegistry struct {
	Ip       string `json:"ip"`
	Type     string `json:"type"`
	Hostname string `json:"hostname"`
	Carrier  struct {
		Name interface{} `json:"name"`
		Mcc  interface{} `json:"mcc"`
		Mnc  interface{} `json:"mnc"`
	} `json:"carrier"`
	Company struct {
		Domain string `json:"domain"`
		Name   string `json:"name"`
		Type   string `json:"type"`
	} `json:"company"`
	Connection struct {
		Asn          int    `json:"asn"`
		Domain       string `json:"domain"`
		Organization string `json:"organization"`
		Route        string `json:"route"`
		Type         string `json:"type"`
	} `json:"connection"`
	Currency struct {
		Code         string `json:"code"`
		Name         string `json:"name"`
		NameNative   string `json:"name_native"`
		Plural       string `json:"plural"`
		PluralNative string `json:"plural_native"`
		Symbol       string `json:"symbol"`
		SymbolNative string `json:"symbol_native"`
		Format       struct {
			Negative struct {
				Prefix string `json:"prefix"`
				Suffix string `json:"suffix"`
			} `json:"negative"`
			Positive struct {
				Prefix string `json:"prefix"`
				Suffix string `json:"suffix"`
			} `json:"positive"`
		} `json:"format"`
	} `json:"currency"`
	Location struct {
		Continent struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"continent"`
		Country struct {
			Area              int           `json:"area"`
			Borders           []interface{} `json:"borders"`
			CallingCode       string        `json:"calling_code"`
			Capital           string        `json:"capital"`
			Code              string        `json:"code"`
			Name              string        `json:"name"`
			Population        int           `json:"population"`
			PopulationDensity float64       `json:"population_density"`
			Flag              struct {
				Emoji        string `json:"emoji"`
				EmojiUnicode string `json:"emoji_unicode"`
				Emojitwo     string `json:"emojitwo"`
				Noto         string `json:"noto"`
				Twemoji      string `json:"twemoji"`
				Wikimedia    string `json:"wikimedia"`
			} `json:"flag"`
			Languages []struct {
				Code   string `json:"code"`
				Name   string `json:"name"`
				Native string `json:"native"`
			} `json:"languages"`
			Tld string `json:"tld"`
		} `json:"country"`
		Region struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"region"`
		City      string  `json:"city"`
		Postal    string  `json:"postal"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Language  struct {
			Code   string `json:"code"`
			Name   string `json:"name"`
			Native string `json:"native"`
		} `json:"language"`
		InEu bool `json:"in_eu"`
	} `json:"location"`
	Security struct {
		IsAbuser        bool `json:"is_abuser"`
		IsAttacker      bool `json:"is_attacker"`
		IsBogon         bool `json:"is_bogon"`
		IsCloudProvider bool `json:"is_cloud_provider"`
		IsProxy         bool `json:"is_proxy"`
		IsRelay         bool `json:"is_relay"`
		IsTor           bool `json:"is_tor"`
		IsTorExit       bool `json:"is_tor_exit"`
		IsVpn           bool `json:"is_vpn"`
		IsAnonymous     bool `json:"is_anonymous"`
		IsThreat        bool `json:"is_threat"`
	} `json:"security"`
	TimeZone struct {
		Id               string    `json:"id"`
		Abbreviation     string    `json:"abbreviation"`
		CurrentTime      time.Time `json:"current_time"`
		Name             string    `json:"name"`
		Offset           int       `json:"offset"`
		InDaylightSaving bool      `json:"in_daylight_saving"`
	} `json:"time_zone"`
	UserAgent struct {
		Header       string `json:"header"`
		Name         string `json:"name"`
		Type         string `json:"type"`
		Version      string `json:"version"`
		VersionMajor string `json:"version_major"`
		Device       struct {
			Brand interface{} `json:"brand"`
			Name  string      `json:"name"`
			Type  string      `json:"type"`
		} `json:"device"`
		Engine struct {
			Name         string `json:"name"`
			Type         string `json:"type"`
			Version      string `json:"version"`
			VersionMajor string `json:"version_major"`
		} `json:"engine"`
		Os struct {
			Name    string      `json:"name"`
			Type    string      `json:"type"`
			Version interface{} `json:"version"`
		} `json:"os"`
	} `json:"user_agent"`
}

func (i *IpRegistry) ToString() string {
	v, _ := json.Marshal(i)
	return string(v)
}

// ParseIp 解析IP
func ParseIp(IpRegistryKey string, ip string) (IpRegistry, error) {
	api := fmt.Sprintf("%s/%s?key=%s", host, ip, IpRegistryKey)
	resp, err := http.Get(api)
	if err != nil {
		fmt.Println(err)
		return IpRegistry{}, errors.New("ipregistry api error")
	}

	defer resp.Body.Close()

	// 读取数据
	v, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return IpRegistry{}, errors.New("ipregistry api error")
	}

	// 解析数据
	var ipRegistry IpRegistry
	err = json.Unmarshal(v, &ipRegistry)
	if err != nil {
		fmt.Println(err)
		return IpRegistry{}, errors.New("ipregistry api error")
	}

	return ipRegistry, nil
}

// IsPass 检测IP是否允许通过
func IsPass(IpRegistryKey string, ip string, userCountryCode string, os []string) (IpRegistry, bool, error) {
	res, err := ParseIp(IpRegistryKey, ip)
	if err != nil {
		return res, false, err
	}

	// isAbuser 是否是滥用者
	if res.Security.IsAbuser {
		return res, false, nil
	}

	// isAnonymous 是否是匿名
	if res.Security.IsAnonymous {
		return res, false, nil
	}

	// isAttacker 是否是攻击者
	if res.Security.IsAttacker {
		return res, false, nil
	}

	// isBogon 是否是保留地址
	if res.Security.IsBogon {
		return res, false, nil
	}

	// isCloudProvider 是否是云服务商
	if res.Security.IsCloudProvider {
		return res, false, nil
	}

	// isProxy 是否是代理
	if res.Security.IsProxy {
		return res, false, nil
	}

	// isRelay 是否是中继
	if res.Security.IsRelay {
		return res, false, nil
	}

	// isThreat 是否是威胁
	if res.Security.IsThreat {
		return res, false, nil
	}

	// isTor 是否是Tor
	if res.Security.IsTor {
		return res, false, nil
	}

	// isTorExit 是否是Tor出口
	if res.Security.IsTorExit {
		return res, false, nil
	}

	// isVpn 是否是VPN
	if res.Security.IsVpn {
		return res, false, nil
	}

	// 判断国家
	if userCountryCode != "" && userCountryCode != res.Location.Country.Code {
		return res, false, nil
	}

	// 判断操作系统
	if len(os) > 0 {
		for _, v := range os {
			if v == res.UserAgent.Os.Name {
				return res, true, nil
			}
		}
		return res, false, nil
	}

	return res, true, nil
}

// IsPassByIpRegistry 检测IP是否允许通过
func IsPassByIpRegistry(res IpRegistry, userCountryCode string, os []string) (IpRegistry, bool, error) {
	// isAbuser 是否是滥用者
	if res.Security.IsAbuser {
		return res, false, nil
	}

	// isAnonymous 是否是匿名
	if res.Security.IsAnonymous {
		return res, false, nil
	}

	// isAttacker 是否是攻击者
	if res.Security.IsAttacker {
		return res, false, nil
	}

	// isBogon 是否是保留地址
	if res.Security.IsBogon {
		return res, false, nil
	}

	// isCloudProvider 是否是云服务商
	if res.Security.IsCloudProvider {
		return res, false, nil
	}

	// isProxy 是否是代理
	if res.Security.IsProxy {
		return res, false, nil
	}

	// isRelay 是否是中继
	if res.Security.IsRelay {
		return res, false, nil
	}

	// isThreat 是否是威胁
	if res.Security.IsThreat {
		return res, false, nil
	}

	// isTor 是否是Tor
	if res.Security.IsTor {
		return res, false, nil
	}

	// isTorExit 是否是Tor出口
	if res.Security.IsTorExit {
		return res, false, nil
	}

	// isVpn 是否是VPN
	if res.Security.IsVpn {
		return res, false, nil
	}

	// 判断国家
	if userCountryCode != "" && userCountryCode != res.Location.Country.Code {
		return res, false, nil
	}

	// 判断操作系统
	if len(os) > 0 {
		for _, v := range os {
			if v == res.UserAgent.Os.Name {
				return res, true, nil
			}
		}
		return res, false, nil
	}

	return res, true, nil
}
