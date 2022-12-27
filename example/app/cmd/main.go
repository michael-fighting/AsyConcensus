package main

import (
	"fmt"
	"github.com/DE-labtory/iLogger"
	"net/http"
	"os"
	"path/filepath"

	"github.com/DE-labtory/cleisthenes"
	"github.com/DE-labtory/cleisthenes/core"

	"github.com/DE-labtory/cleisthenes/config"
	"github.com/DE-labtory/cleisthenes/example/app"
	kitlog "github.com/go-kit/kit/log"
)

func main() {
	//host := flag.String("address", "127.0.0.2", "Application address")
	//port := flag.Int("port", 8000, "Application port")
	//configPath := flag.String("config", "", "User defined config path")
	//flag.Parse()
	fmt.Println(os.Args)
	host := os.Args[1]
	port := os.Args[2]
	configPath := os.Args[3]

	//初始化日志文件路径
	os.RemoveAll("./log")
	absPath, _ := filepath.Abs("./log/cleisthenes" + port + ".log")
	defer os.RemoveAll("./log")
	err := iLogger.EnableFileLogger(true, absPath)

	address := fmt.Sprintf("%s:%s", host, port)

	kitLogger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	kitLogger = kitlog.With(kitLogger, "ts", kitlog.DefaultTimestampUTC)  //DefaultTimestampUTC返回现在的世界时
	httpLogger := kitlog.With(kitLogger, "component", "http")

	//将defaultConfig的默认配置文件写到tempConfig/node/config.yml
	config.Init(configPath)

	txValidator := func(tx cleisthenes.Transaction) bool {
		// custom transaction validation logic，交易验证逻辑
		return true
	}

	node, err := core.New(txValidator)
	if err != nil {
		panic(fmt.Sprintf("Cleisthenes instantiate failed with err: %s", err))
	}

	go func() {
		for {
			result := <-node.Result()
			fmt.Printf("[DONE]epoch : %d, batch tx count : %d\n", result.Epoch, len(result.Batch))
		}
	}()

	go func() {
		httpLogger.Log("message", "hbbft started")
		node.Run()
	}()

	httpLogger.Log("message", fmt.Sprintf("http server started: %s", address))
	if err := http.ListenAndServe(address, app.NewApiHandler(node, httpLogger)); err != nil {
		httpLogger.Log("message", fmt.Sprintf("http server closed: %s", err))
	}
}
