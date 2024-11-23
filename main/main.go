package main

import (
	"fmt"
	"os"
	"runtime"
	"path/filepath"
	core "github.com/gzjjjfree/hello"
)

var (
	version  = "1.0.0.0"
	build    = "Custom"
	codename = "一个学习代理转发的软件"
	intro    = "从 v2ray 借鉴学习"
)

func Version() string {
	return version
}

func VersionStatement() string {
	return fmt.Sprintf("HELLO %s (%s) %s (%s %s/%s)\n%s", Version(), codename, build, runtime.Version(), runtime.GOOS, runtime.GOARCH, intro)	
}

func startHELLO() {
	configFile, err := getConfigFilePath()
	config, err := core.Configload(configFile)
	if err != nil {
		fmt.Println("read the configFile err is: ", err)
	} else {
		fmt.Println("configFile from: ", config)
	}
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
	fmt.Println(VersionStatement())
	startHELLO()
}