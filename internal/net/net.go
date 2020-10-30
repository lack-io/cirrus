package net

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	ErrInvalidProxy = errors.New("invalid proxy ip")
	ErrNetwork = errors.New("bad network")
)

const UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36"

type PublicIP struct {
	// 公网 IP
	IP string `json:"ip,omitempty"`

	// 国家
	Country string `json:"country,omitempty"`

	// 地区
	Area string `json:"area,omitempty"`

	// 省份
	Province string `json:"province,omitempty"`

	// 城市
	City string `json:"city,omitempty"`

	// 运营商
	ISP string `json:"isp,omitempty"`

	// 时间戳
	Timestamp string `json:"timestamp,omitempty"`
}

// GetPublicIP 获取本机的公网 IP
func GetPublicIP() (*PublicIP, error) {
	return getMyIP("")
}

func getMyIP(proxy string) (*PublicIP, error) {
	client := &http.Client{}

	if proxy != "" {
		urli := url.URL{}
		urlproxy, err := urli.Parse(proxy)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidProxy, err)
		}

		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(urlproxy),
		}
	}

	rqt, err := http.NewRequest("GET", "http://myip.top", nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNetwork, err)
	}
	rqt.Header.Add("User-Agent", UserAgent)
	response, _ := client.Do(rqt)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	ip := &PublicIP{}
	err = json.Unmarshal(body, ip)
	return ip, nil
}

// CheckProxyIP 检测 proxy 是否有效
func CheckProxyIP(proxy string) (bool, error) {
	ip, port, err := ParseProxy(proxy)
	if err != nil {
		return false, err
	}

	pip, err := getMyIP(fmt.Sprintf("http://%s:%s", ip, port))
	if err != nil {
		return false, err
	}

	return pip.IP == ip, nil
}

// ParseProxy 解析 ip:port 格式的字符串
func ParseProxy(proxy string) (string, string, error) {
	sp := strings.SplitN(proxy, ":", 2)
	if len(sp) != 2 {
		return "", "", ErrInvalidProxy
	}
	return sp[0], sp[1], nil
}