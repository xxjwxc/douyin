# douyin_msg
golang 抖音弹幕消息

### 房间id通过抖音官网浏览器中获取:如https://live.douyin.com/533646158504，roomId为533646158504

### 是否用方法
```go
	tmp, err := douyin.NewRoom("533646158504")
	tmp.Start(onMessage)

	time.Sleep(1 * time.Minute)
	tmp.Stop()
	fmt.Println(tmp, err)
```