package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	url2 "net/url"
	"strings"
	"time"
)

var conf *configS

func init() {
	initConfig()
}

func main() {
	for {
		watch()
		time.Sleep(3 * time.Second)
	}
}

func watch() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("recover panic: ", err)
			return
		}
	}()
	res, err := _request(conf.WatchUrl)
	if err != nil {
		log.Println("监控商品失败:", err.Error())
		notifyWechat("监控商品失败", err.Error()+"\n"+conf.WatchUrl)
		return
	}
	if !strings.Contains(res, "200") {
		msg := "监控商品失败:返回结果不是200"
		log.Println(msg)
		notifyWechat("监控商品结果非200失败", conf.WatchUrl)
		return
	}
	if strings.Contains(res, "脱销") {
		log.Println("脱销状态:", conf.GoodsUrl)
	} else {
		log.Println("请及时购买:", conf.GoodsUrl)
		notifyWechat("请及时购买", conf.GoodsUrl)
	}
}

func _request(url string) (string, error) {
	timeoutCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	request, err := http.NewRequestWithContext(timeoutCtx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("构建请求失败: %s", err.Error())
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %s", err.Error())
	}
	resBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("读取结果失败: %s", err.Error())
	}
	_ = response.Body.Close()
	return string(resBytes), nil
}

var canNotifyChan = make(chan struct{}, 1)

func init() {
	go func() {
		for {
			canNotifyChan <- struct{}{}
			time.Sleep(12 * time.Second)
		}
	}()
}

// https://www.pushplus.plus/push1.html 推送
func notifyWechat(title, content string) {
	select {
	case <-canNotifyChan:
		log.Println("每12秒钟只能推送一次, 因为push_plus限制")
		return
	default:

	}
	if conf.PushPlusToken == "" {
		return
	}
	//content += fmt.Sprintf("%%0D%%0A%%0D%%0A%d请及时关闭监控服务，否则会继续发送通知消息", time.Now().Unix()) // 换行加随机串是因为server酱不让发送相同内容
	content += fmt.Sprintf("\n\n%%0D%%0A%%0D%%0A%d请及时关闭监控服务，否则会继续发送通知消息", time.Now().Unix()) // 换行加随机串是因为server酱不让发送相同内容
	//content += "\n\n请及时关闭监控服务，否则会继续发送通知消息"
	content = url2.QueryEscape(content)
	url := fmt.Sprintf("http://www.pushplus.plus/send?token=%s&title=%s&content=%s&template=txt",
		conf.PushPlusToken, url2.QueryEscape(title), content)

	res, err := _request(url)
	if err != nil {
		log.Println("通知微信失败:", err.Error())
		return
	}
	log.Println("通知 push_plus:", res)
}

type configS struct {
	PushPlusToken string `json:"push_plus_token"`
	WatchUrl      string
	GoodsUrl      string
}

func initConfig() {
	bytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic("读取配置文件失败: " + err.Error())
	}
	conf = new(configS)
	if err := json.Unmarshal(bytes, conf); err != nil {
		panic("配置文件格式错误: " + err.Error())
	}
	if conf.WatchUrl == "" || conf.GoodsUrl == "" {
		panic("商品监控链接或者商品链接不能为空")
	}
	if conf.PushPlusToken == "" {
		fmt.Println("因为 pushplus token 为空，所以不会通知到微信")
	}
}
