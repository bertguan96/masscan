# 一款golang操作masscan的工具库

## 1.说明

工具基于以下代码改造完成，暂时只支持简单的命令执行，欢迎pr新命令。

[1] https://github.com/zan8in/masscan/tree/main

[2] https://github.com/Ullaakut/nmap

## 2.使用方式

```json
	scanner, err := NewScanner(
		WithRoot(),                      // 设置是否是root权限
		WithBinaryPath("./bin/masscan"), // 设置执行的bin路径
		WithPort("1-65535"),             // 设置扫描的端口范围
		WithOutputJson(),                // 设置输出数据格式
		WithRetryTime(3),                // 设置重复探测次数
		WithRate(40000),                 // 设置扫描频率
		WithRandomizeHosts(),            // 设置随机主机
		WithTarget([]string{"114.55.97.220", "43.135.11.122", "101.42.164.23", "43.153.24.244"})) // 设置扫描目标
	if err != nil {
		log.Errorf("new scanner error, %s", err)
		return
	}
	res, _, err := scanner.Run()
	if err != nil {
		log.Errorf("exec error, %s", err)
		return
	}

	marshal, err := json.Marshal(res)
	if err != nil {
    log.Errorf("json parser error, %s", err)
		return
	}
	log.Info(marshal)
```

