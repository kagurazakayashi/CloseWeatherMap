package main

import (
	"flag"
	"log"
)

var (
	xlsxFilePath string
	titles       []string
	datas        [][]string
	dataLen      int
	listenHost   string
	appids       string
	uri          string
)

func main() {
	log.Println("XLSWeather 1.0.0")
	flag.StringVar(&xlsxFilePath, "f", "", "XLSX 文件路径。")
	flag.StringVar(&listenHost, "l", "127.0.0.1:80", "HTTP 接口所使用的 <IP>:<端口号>，不提供 IP 则允许所有 IP。")
	flag.StringVar(&uri, "u", "/data/2.5/weather", "HTTP 接口的 URI。")
	flag.StringVar(&appids, "a", "", "限制只有指定的几个 APPID 才能访问，使用英文逗号分隔。留空则不限制。")
	flag.Parse()

	if len(xlsxFilePath) < 6 {
		log.Println("你必须使用 -f <文件.xlsx> 指定一个 XLSX 文件。")
		return
	}

	reloadXLSX()

	if !initweb() {
		return
	}
}

func reloadXLSX() {
	titles, datas = loadXLSX(true)
	if len(datas) == 0 || len(titles) == 0 {
		return
	}
	dataLen = len(datas)
	log.Println("读取文件:", xlsxFilePath, "完成，数据量:", dataLen)
}
