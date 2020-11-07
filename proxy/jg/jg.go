// 极光爬虫代理
package jg

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	json "github.com/json-iterator/go"

	"github.com/xingyys/cirrus/config"
	"github.com/xingyys/cirrus/internal/net"
	"github.com/xingyys/cirrus/proxy"
)

// 免费代理套餐的ID
const balanceAction = "31731"

type Result struct {
	Code interface{} `json:"code,omitempty"`

	Success bool `json:"success,omitempty"`

	Msg string `json:"msg,omitempty"`

	Data interface{} `json:"data,omitempty"`
}

// While JG 报名单信息
type White struct {
	// 白名单IP
	IP string `json:"mark_ip,omitempty"`

	// 更新时间，格式为时间戳，单位为秒
	UpdatedAt string `json:"updated_at,omitempty"`
}

type JG struct {
	// ctx 控制 JG 的停止
	ctx context.Context

	// 代理名称，默认为 "jiguang"
	Name string `json:"name,omitempty"`

	// 本机公共IP
	LocalIP string `json:"localIP,omitempty"`

	Neek string `json:"neek,omitempty"`

	APIAppKey string `json:"api_appKey,omitempty"`

	BalanceAppKey string `json:"balance_appKey,omitempty"`

	// 白名单信息
	Whites []White `json:"whites,omitempty"`
}

func GetJGProxy(ctx context.Context, cfg *config.ProxyJG) (*JG, error) {
	jg := &JG{
		ctx:  ctx,
		Name: "jiguang",
	}

	if cfg == nil {
		return nil, fmt.Errorf("%w: config is nil", proxy.ErrUnable)
	}

	jg.Neek = cfg.Neek
	jg.APIAppKey = cfg.APIAppKey
	jg.BalanceAppKey = cfg.BalanceAppKey

	return jg, nil
}

// Init implement proxy.Proxy
func (j *JG) Init() error {

	if j.ctx == nil {
		return fmt.Errorf("%w: ctx is nil", proxy.ErrUnable)
	}

	// 获取公共IP
	pip, err := net.GetPublicIP()
	if err != nil {
		return err
	}
	j.LocalIP = pip.IP
	type lists struct {
		Lists []White `json:"lists,omitempty"`
	}

	params := map[string]string{
		"neek":   j.Neek,
		"appkey": j.APIAppKey,
	}

	out := &lists{}
	err = j.Fetch(j.ctx, "http://webapi.jghttp.golangapi.com/index/index/white_list", params, out)
	if err != nil {
		return fmt.Errorf("%w: %v", proxy.ErrUnable, err)
	}

	j.Whites = out.Lists
	// 确认本机 IP 是否在代理商的白名单之内，否则就添加本地IP
	exists := false
	for _, item := range out.Lists {
		if item.IP == j.LocalIP {
			exists = true
		}
	}

	if !exists {
		// 添加本机的公共 IP 为白名单
		params1 := params
		params1["white"] = pip.IP
		err = j.Fetch(j.ctx, "http://webapi.jghttp.golangapi.com/index/index/save_white", params, nil)
		if err != nil {
			return err
		}
		j.Whites = append(j.Whites, White{
			IP:        pip.IP,
			UpdatedAt: fmt.Sprintf("%d", time.Now().Unix())},
		)
	}

	// 领取免费的代理IP
	data := map[string]string{"mid": j.Neek}
	_, _ = j.post(j.ctx, "http://webapi.jghttp.golangapi.com/index/users/get_day_free_pack", nil, data, nil)

	return nil
}

// GetEndpoints implement proxy.Proxy
func (j *JG) GetEndpoints(ctx context.Context, n int) ([]*proxy.Endpoint, error) {
	endpoints := []*proxy.Endpoint{}
	pb, err := j.GetPackageBalance(ctx)
	if err != nil {
		return nil, err
	}
	sub := pb - n
	if sub >= 0 {
		// 优先使用免费代理
		endpoints, err = j.getips(ctx, n, balanceAction)
		if err != nil {
			return j.getips(ctx, n, "")
		}
	} else {
		remaining := n - pb
		endpoints, err = j.getips(ctx, pb, balanceAction)
		if err != nil {
			remaining = n
		}
		es, err := j.getips(j.ctx, remaining, "")
		if err == nil {
			endpoints = append(endpoints, es...)
		}
	}

	return j.getips(ctx, n, "")
}

// GetBalance implement proxy.Proxy
func (j *JG) GetBalance(ctx context.Context) (*proxy.Balance, error) {

	type balanceData struct {
		Balance string `json:"balance,omitempty"`
	}

	params := map[string]string{
		"neek":   j.Neek,
		"appkey": j.APIAppKey,
	}

	bd := &balanceData{}
	err := j.Fetch(ctx, "http://webapi.jghttp.golangapi.com/index/index/get_my_balance", params, bd)
	if err != nil {
		return nil, err
	}

	i, _ := strconv.ParseFloat(bd.Balance, 64)
	return &proxy.Balance{Amount: i, Coin: proxy.RMB}, nil
}

// GetPackageBalance 返回免费代理的余额
func (j *JG) GetPackageBalance(ctx context.Context) (int, error) {
	type Out struct {
		PB int `json:"package_balance,omitempty"`
	}

	params := map[string]string{
		//"ac":     balanceAction,
		"neek":   j.Neek,
		"appkey": j.BalanceAppKey,
	}

	// 获取免费代理的可用余额
	out := Out{}
	err := j.Fetch(j.ctx, "http://webapi.jghttp.golangapi.com/index/index/get_my_pack_info", params, &out)
	return out.PB, err
}

// Fetch 请求 JG 代理接口
func (j *JG) Fetch(ctx context.Context, api string, params map[string]string, to interface{}) error {

	r, err := j.fetch(ctx, api, params, to)
	if err != nil {
		return err
	}

	code := fmt.Sprintf("%v", r.Code)
	if code != "0" {
		switch code {
		// JG 代理的接口有访问限制，2s内只能访问一次
		// code 为 111，等待 2s 后重试
		case "111":
			<-time.After(time.Second * 2)
			_, err = j.fetch(ctx, api, params, to)
			return err
		case "114":
			return fmt.Errorf("%w: %v", proxy.ErrInsufficient, r.Msg)
		case "121":
			return fmt.Errorf("%w: %v", proxy.ErrPackageExpired, r.Msg)
		default:
			return fmt.Errorf("%w: %v", proxy.ErrResultException, r.Msg)
		}
	}

	return nil
}

func (j *JG) post(ctx context.Context, url string, params, data map[string]string, to interface{}) (*Result, error) {
	cli := net.Builder().
		AddParams(params).
		SetContentType(net.FORM).
		SetHeader("Host", "webapi.jghttp.golangapi.com")

	reader := strings.NewReader("")
	if data != nil {
		b, err := json.Marshal(data)
		if err == nil {
			_, _ = reader.Read(b)
		}
	}

	out, err := cli.Post(ctx, url, reader)

	if err != nil {
		return nil, fmt.Errorf("%w: %v", proxy.ErrResultException, err)
	}

	r := &Result{Data: to}
	_ = json.Unmarshal(out, r)

	return r, nil
}

func (j *JG) fetch(ctx context.Context, url string, params map[string]string, to interface{}) (*Result, error) {

	cli := net.Builder().
		AddParams(params).
		SetContentType(net.JSON).
		SetHeader("Host", "webapi.jghttp.golangapi.com")

	data, err := cli.Get(ctx, url)

	if err != nil {
		return nil, fmt.Errorf("%w: %v", proxy.ErrResultException, err)
	}

	r := &Result{Data: to}
	_ = json.Unmarshal(data, r)

	return r, nil
}

// getips 获取
func (j *JG) getips(ctx context.Context, num int, pack string) ([]*proxy.Endpoint, error) {
	params := map[string]string{
		"num":     fmt.Sprintf("%d", num),
		"type":    "2",
		"pro":     "0",
		"city":    "0",
		"yys":     "0",
		"port":    "11",
		"ts":      "1",
		"ys":      "1",
		"cs":      "1",
		"lb":      "1",
		"sb":      "0",
		"pb":      "45",
		"mr":      "1",
		"regions": "",
	}
	if pack != "" {
		params["pack"] = pack
	}

	out := []*proxy.Endpoint{}
	err := j.Fetch(ctx, "http://d.jghttp.golangapi.com/getip", params, &out)
	return out, err
}
