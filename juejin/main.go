package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
	"github.com/tencentyun/scf-go-lib/functioncontext"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"io/ioutil"
	"net/http"
	"selfProject/auto-helper/juejin/pkg"
	"strconv"
	"strings"
)

var (
	ServerType = kingpin.Flag("server", "Set up server type, Default is tencent scf and main run, eg --server=\"tencent\" or --server=\"run\" ").Default("tencent").String()
	WeConnUrl = kingpin.Flag("url", "Set up we conn bot URL, Default is eg").Default("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxxxxxxxxxx").String()
	Cookie    = kingpin.Flag("cookie", "Set up jue jin cookie").Default("cookie").String()
)

const (
	SignIn             = "https://api.juejin.cn/growth_api/v1/check_in"
	Lottery            = "https://api.juejin.cn/growth_api/v1/lottery/draw"
	DefaultContentType = "application/json;charset=UTF-8"
	UserAgent          = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.106 Safari/537.36"
)

func Alert(content string) error {
	msg := pkg.WechatMarkdown{
		MsgType: "markdown",
		Markdown: &pkg.Markdown{
			Content: content,
		},
	}
	if _, err := Post(*WeConnUrl, msg, map[string]string{
		"Content-Type": DefaultContentType,
	}); err != nil {
		return err
	}
	return nil
}

func Post(url string, msg interface{}, headers map[string]string) ([]byte, error) {
	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(jsonBytes)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	for k, v := range getCookie() {
		request.AddCookie(&http.Cookie{Name: k, Value: v})
	}
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		
		}
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func getCookie() map[string]string {
	slice := strings.Split(*Cookie, ";")
	cookieMap := make(map[string]string)
	for _, v := range slice {
		v = strings.Trim(v, " ")
		index := strings.Index(v, "=")
		cookieMap[v[0:index]] = v[index+1:]
	}
	return cookieMap
}

func getSignIn() string {
	res, _ := Post(SignIn, nil, map[string]string{
		"Content-Type": DefaultContentType,
		"User-Agent":   UserAgent,
	})
	// marshal
	var content string
	resp := &pkg.JueJinSignInResp{}
	if err := json.Unmarshal(res, resp); err != nil {
		return ""
	}
	if resp.ErrNo == 0 {
		content = "签到信息: " + resp.ErrMsg + "新增矿石数：" + strconv.FormatInt(resp.Data.IncrPoint, 10) + ", 总计矿石数：" + strconv.FormatInt(resp.Data.SumPoint, 10)
	}
	if resp.ErrNo == 15001 {
		content = "签到信息: " + resp.ErrMsg
	}
	return content
}

func getLottery() string {
	res, _ := Post(Lottery, nil, map[string]string{
		"Content-Type": DefaultContentType,
		"User-Agent":   UserAgent,
	})
	fmt.Println(string(res))
	// marshal
	var content string
	resp := &pkg.JueJinLotteryResp{}
	if err := json.Unmarshal(res, resp); err != nil {
		return ""
	}
	if resp.ErrNo == 0 {
		content = "幸运抽奖信息: " + resp.ErrMsg + "奖品：" + resp.Data.LotteryName
	}
	if resp.ErrNo == 7003 {
		content = "幸运抽奖信息: " + resp.ErrMsg
	}
	return content
}

func hello(ctx context.Context, event pkg.DefineEvent) (string, error) {
	lc, _ := functioncontext.FromContext(ctx)
	fmt.Printf("ctx: %#v\n", lc)
	fmt.Printf("namespace: %s\n", lc.Namespace)
	fmt.Printf("function name: %s\n", lc.FunctionName)
	// SignIn
	signInContent := getSignIn()
	// Alert
	if err := Alert(signInContent); err != nil {
		return "", err
	}
	// Lottery
	lotteryContent := getLottery()
	// Alert
	if err := Alert(lotteryContent); err != nil {
		return "", err
	}
	return fmt.Sprintf("Hello!"), nil
}

func run() error {
	// SignIn
	signInContent := getSignIn()
	// Alert
	if err := Alert(signInContent); err != nil {
		return err
	}
	// Lottery
	lotteryContent := getLottery()
	// Alert
	if err := Alert(lotteryContent); err != nil {
		return err
	}
	return nil
}

func main() {
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	switch *ServerType {
	case "tencent":
		// Make the handler available for Remote Procedure Call by Cloud Function
		cloudfunction.Start(hello)
	default:
		if err := run(); err != nil {
			print(err)
		}
	}
}
