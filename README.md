# Why?
穷，想买 mbp 翻新版，刚开始有货，但是第二天打开网页的时候没货了。。。一怒之下写了个程序监控货源状态  
2019-12-09 14:21 监控到了,已下单

# 配置文件 config.json
```json
{
  "push_plus_token": "", // pushplus 通知 token,不填则不通知微信
  
  //打开浏览器控制台，然后刷新要买的 apple 翻新产品，应该能看到会发起一个 /shop/fulfillment-messages 异步请求，把完整链接复制到这里就行
  "watchUrl": "https://www.apple.com.cn/shop/fulfillment-messages?little=false&parts.0=FU0Q2CH/A&mts.0=regular&fts=true",
  "goodsUrl": "https://www.apple.com.cn/shop/product/FU0X2CH/A?fnode=112db2e59016208309c3c507e114dc5b2d03538c81ffa60b01f08b0ebe7e4fe2e26318603270c010447d515dea3124861b6145541b0fb16faf39d9b9a2629908e387a438ae8a5e8a48e956ee9ce0001c" // 浏览器打开商品，复制链接
}
```