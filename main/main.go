package main

import (
	"fmt"
	"os"
	"path/filepath"
	"errors"
	"runtime"
	"os/signal"
	"syscall"

	core "github.com/gzjjjfree/hello"
	_ "github.com/gzjjjfree/hello/proxy/vmess/inbound"
	_ "github.com/gzjjjfree/hello/proxy/vmess/outbound"
	_ "github.com/gzjjjfree/hello/app/dns"
	_ "github.com/gzjjjfree/hello/app/routing"
	_ "github.com/gzjjjfree/hello/transport/internet/tcp"
	_ "github.com/gzjjjfree/hello/common/bytespool"
)

func startHello() (core.Server, error) {
	configFile, err := getConfigFilePath()
	config, err := core.Configload(configFile)
	if err != nil {
		return nil, err
	} 

	server, err := core.New(config)
	if server == nil {
		return nil, errors.New("server is error")
	}
	return server, nil
}

func getConfigFilePath() (string, error) {
	workingDir, err := os.Getwd()
	if err != nil { // Getwd 返回当前目录对应的根路径名。如果当前目录可以通过多条路径到达（由于符号链接），Getwd 可能会返回其中的任意一条。
		// 当没有预设路径时，读取当前目录根路径的 config.json			
		return fmt.Sprintln("workingDir is err"), err		
	}
	configFile := filepath.Join(workingDir, "config.json") // Join 将任意数量的路径元素合并为一个路径，并使用特定于操作系统的 [Separator] 将它们分隔开。空元素将被忽略。
		// 结果为 Cleaned。但是，如果参数列表为空或其所有元素都为空，则 Join 将返回一个空字符串。在 Windows 上，如果第一个非空元素是 UNC 路径，则结果将仅为 UNC 路径
	return configFile, err
}
func main () {
	fmt.Println(core.VersionStatement())
	server, err := startHello()
	if err != nil {
		fmt.Println(err)
		os.Exit(23)
	}
	if err := server.Start(); err != nil {
		fmt.Println("Failed to start", err)
		os.Exit(-1)
	}
	defer server.Close()

	runtime.GC()

	{
		// osSignals: 声明了一个无缓冲的 channel，其类型为 os.Signal，用于接收操作系统信号
		// 创建一个容量为 1 的 channel。容量为 1 表示该 channel 最多只能存储一个信号
		osSignals := make(chan os.Signal, 1)
		// signal.Notify: 这个函数的作用是将指定的信号注册到指定的 channel 上
		// osSignals: 上面创建的 channel，用于接收信号
		// os.Interrupt 和 syscall.SIGTERM: 两个常见的操作系统信号，分别表示用户中断 (Ctrl+C) 和终止进程
		signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
		// <-osSignals: 从 osSignals channel 接收一个信号。当程序执行到这一行时，会阻塞，直到有信号被发送到该 channel
		<-osSignals
	}
}