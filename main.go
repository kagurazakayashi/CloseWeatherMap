package main

import (
	"flag"
	"log"
)

var (
	xlsxFilePath string
	datas        [][]string
	listenHost   string
)

func main() {
	log.Println("XLSWeather 1.0.0")
	flag.StringVar(&xlsxFilePath, "f", "", "XLSX 文件路径")
	flag.StringVar(&listenHost, "l", "127.0.0.1:80", "监听地址")
	flag.Parse()

	if len(xlsxFilePath) < 6 {
		log.Println("你必须使用 -f <文件.xlsx> 指定一个 XLSX 文件。")
		return
	}

	log.Println("正在加载数据文件:", xlsxFilePath)
	datas = loadXLSX()
	if len(datas) == 0 {
		return
	}
	log.Println("数据量:", len(datas))

	initweb()
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
