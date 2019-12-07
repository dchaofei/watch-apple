# Why?
穷，想买 mbp 翻新版，刚开始有货，但是第二天打开网页的时候没货了。。。一怒之下写了个程序监控货源状态

# 配置文件 config.json
```json
{
  "serverKey": "", // Server酱通知微信key,不填则不通知微信
  "watchUrl": "https://www.apple.com.cn/cn-k12/shop/delivery-message?parts.0=G0W53CH%2FA", // 打开商品，打开控制台，刷新页面，复制异步请求的这个url,记得去掉后边的时间戳
  "goodsUrl": "https://www.apple.com.cn/cn-k12/shop/product/G0W53CH/A?fnode=82c65a2ad3bfc35c5451e4561b930167fc6b5c9206bde95e29fffc095a26ca15d3cce09504e07e247561cb9af06f3b5f4523847992aebd57d42368da886f8040ce2c1c35b52e323a206464448e2bba2f" // 浏览器打开商品，复制链接
}
```
