package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	watchUrl        string
	goodsUrl        string
	notifyWechatKey string
)

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
	res, err := _request(watchUrl)
	if err != nil {
		log.Println("监控mbp失败:", err.Error())
		notifyWechat("监控mbp失败", err.Error()+"\n"+watchUrl)
		return
	}
	if !strings.Contains(res, "200") {
		msg := "监控mbp失败:返回结果不是200"
		log.Println(msg)
		notifyWechat("监控mbp结果非200失败", watchUrl)
		return
	}
	if strings.Contains(res, "脱销") {
		log.Println("脱销状态:", goodsUrl)
	} else {
		log.Println("请及时购买:", goodsUrl)
		notifyWechat("请及时购买", goodsUrl)
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

func notifyWechat(title, content string) {
	if notifyWechatKey == "" {
		return
	}
	content += fmt.Sprintf("%%0D%%0A%%0D%%0A%d请及时关闭监控服务，否则会继续发送通知消息", time.Now().Unix()) // 换行加随机串是因为server酱不让发送相同内容
	url := fmt.Sprintf("https://sc.ftqq.com/%s.send?text=%s&desp=%s", notifyWechatKey, title, content)
	res, err := _request(url)
	if err != nil {
		log.Println("通知微信失败:", err.Error())
		return
	}
	log.Println("通知server酱:", res)
}

type configS struct {
	ServerKey string
	WatchUrl  string
	GoodsUrl  string
}

func initConfig() {
	bytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic("读取配置文件失败: " + err.Error())
	}
	c := new(configS)
	if err := json.Unmarshal(bytes, c); err != nil {
		panic("配置文件格式错误: " + err.Error())
	}
	watchUrl = c.WatchUrl
	goodsUrl = c.GoodsUrl
	notifyWechatKey = c.ServerKey
	if watchUrl == "" || goodsUrl == "" {
		panic("商品监控链接或者商品链接不能为空")
	}
	if notifyWechatKey == "" {
		fmt.Println("因为server酱key为空，所以不会通知到微信")
	}
}
