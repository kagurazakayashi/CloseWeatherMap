package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	xlsxFilePath     string
	titles           []string
	datas            [][]string
	dataLen          int
	listenHost       string
	appids           string
	uri              string
	baseDay          string // YYYYMMDD
	baseDayDate      time.Time
	forceReload      bool
	lockTimeStr      string
	lockTime         time.Time
	timeLayout       string = "2006-01-02 15:04:05"
	verbose          bool
	reverseDirection bool
)

func main() {
	log.Println("XLSWeather 1.1.0  " + time.Now().Format(timeLayout))
	flag.StringVar(&xlsxFilePath, "f", "", "XLSX 文件路径。")
	flag.StringVar(&baseDay, "d", "", "基准日期(YYYYMMDD)，为空则为当前日期。")
	flag.StringVar(&listenHost, "l", "127.0.0.1:80", "HTTP 接口所使用的 <IP>:<端口号>，不提供 IP 则允许所有 IP。")
	flag.StringVar(&uri, "u", "/data/2.5/weather", "HTTP 接口的 URI。")
	flag.StringVar(&appids, "a", "", "限制只有指定的几个 APPID 才能访问，使用英文逗号分隔。留空则不限制。")
	flag.BoolVar(&forceReload, "r", false, "强制重新加载 XLSX 文件。")
	flag.StringVar(&lockTimeStr, "t", "", "强制按指定时间提供数据，格式示例: 2006-01-02 15:04:05")
	flag.BoolVar(&verbose, "v", false, "显示详细信息用于调试。")
	flag.BoolVar(&reverseDirection, "rd", false, "反转风向数据。")
	flag.Parse()

	if len(xlsxFilePath) < 6 {
		log.Println("你必须使用 -f <文件.xlsx> 指定一个 XLSX 文件。")
		return
	}

	if len(lockTimeStr) > 0 {
		lockTimeN, err := time.ParseInLocation(timeLayout, lockTimeStr, time.Local)
		if err != nil {
			log.Println("错误：时间格式不正确，使用当前时间。")
		}
		lockTime = lockTimeN
	}

	if !genBaseDay("") {
		return
	}

	fmt.Println(strings.Join([]string{
		"XLSX 文件路径: " + xlsxFilePath,
		"基准日期: " + baseDayDate.Format(timeLayout),
		"强制时间: " + lockTime.Format(timeLayout),
		"HTTP 接口地址: " + listenHost + uri,
		"APPID 限制: " + appids,
		"强制重新加载: " + strconv.FormatBool(forceReload),
		"详细信息: " + strconv.FormatBool(verbose),
		"反转风向数据: " + strconv.FormatBool(reverseDirection),
	}, "\n"))

	reloadXLSX()

	if !initweb() {
		return
	}
}

func nowTime() time.Time {
	if len(lockTimeStr) > 0 {
		return lockTime
	} else {
		return time.Now()
	}
}

func genBaseDay(baseDayI string) bool { // ->baseDayDate
	if len(baseDayI) == 0 {
		baseDayI = baseDay
	}
	if len(baseDayI) == 0 {
		var nowTime time.Time = nowTime()
		baseDayDate = time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 0, 0, 0, 0, time.Local)
		return true
	}
	if len(baseDayI) == 8 {
		year, err := strconv.Atoi(baseDay[0:4])
		if err != nil {
			log.Println("错误：基础年份输入不正确。")
			return false
		}
		month, err := strconv.Atoi(baseDay[4:6])
		if err != nil {
			log.Println("错误：基础月份输入不正确。")
			return false
		}
		day, err := strconv.Atoi(baseDay[6:8])
		if err != nil {
			log.Println("错误：基础日期输入不正确。")
			return false
		}
		baseDayDate = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
		return true
	}
	log.Println("错误：基础日期输入不正确。")
	return false
}

func daysApart(t1, t2 time.Time) int {
	startOfDay1 := time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, t1.Location())
	startOfDay2 := time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, t2.Location())
	days := int(startOfDay2.Sub(startOfDay1).Hours() / 24)
	return days
}

func reloadXLSX() {
	titles, datas = loadXLSX(true)
	if len(datas) == 0 || len(titles) == 0 {
		return
	}
	datas = reverseDirectionDatas(datas)
	dataLen = len(datas)
	log.Println("读取文件:", xlsxFilePath, "完成，数据量:", dataLen)
}
