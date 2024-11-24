package core

import (
	"encoding/json"
	"fmt"
	"os"
)

func Configload(input string) (*Config, error) {
	//fmt.Println("input: ", input)
	file, err := os.Open(input)
	if err != nil {
		fmt.Println("读取 config.json 文件出错: ", err)
		return nil, err
	}
	defer file.Close()

	// 创建解码器
	decoder := json.NewDecoder(file)

	// 将 JSON 数据解码到 Person 结构体
	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("解析 config.json 文件出错: ", err)
		return nil, err
	}
	//fmt.Println("config: ", config)
	return &config, nil
}
