//go:generate goversioninfo -icon=ico/icon.ico -manifest=main.exe.manifest
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	server           *http.Server
	xlsxFilePath     string
	titles           []string
	datas            [][]string
	dataLen          int
	dataLenLen       int
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
	tzClientLockName string
	tzClientLock     *time.Location
	tzServerLockName string
	tzServerLock     *time.Location
	hostEntry        string
)

func main() {
	log.Println("XLSWeather 1.2.0  " + time.Now().Format(timeLayout))
	fmt.Println("帮助和更新: https://github.com/kagurazakayashi/xlsweather")
	flag.StringVar(&xlsxFilePath, "f", "", "XLSX 文件路径。")
	flag.StringVar(&baseDay, "d", "", "基准日期(YYYYMMDD)，为空则为当前日期。")
	flag.StringVar(&listenHost, "l", "127.0.0.1:80", "HTTP 接口所使用的 <IP>:<端口号>，不提供 IP 则允许所有 IP。")
	flag.StringVar(&uri, "u", "/data/2.5/weather", "HTTP 接口的 URI。")
	flag.StringVar(&appids, "a", "", "限制只有指定的几个 APPID 才能访问，使用英文逗号分隔。留空则不限制。")
	flag.BoolVar(&forceReload, "r", false, "强制重新加载 XLSX 文件。")
	flag.StringVar(&lockTimeStr, "t", "", "强制按指定时间提供数据，格式示例: 2006-01-02 15:04:05")
	flag.BoolVar(&verbose, "v", false, "显示详细信息用于调试。")
	flag.BoolVar(&reverseDirection, "rd", false, "反转风向数据。")
	flag.StringVar(&tzClientLockName, "tc", "", "强制客户端时区为指定的 IANA 时区名称，例如 Europe/Paris 。")
	flag.StringVar(&tzServerLockName, "ts", "", "强制 XLSX 文件时区为指定的 IANA 时区名称，例如 Asia/Tokyo 。")
	flag.StringVar(&hostEntry, "host", "", "启动时临时添加一条项目到 hosts 文件中，结束时删除。格式: `[IP] [HOST]`")
	flag.Parse()

	if len(xlsxFilePath) < 6 {
		log.Println("你必须使用 -f <文件.xlsx> 指定一个 XLSX 文件。")
		return
	}

	if len(tzClientLockName) > 0 {
		if tzClientLockName == "Local" {
			tzClientLock = time.Local
		} else {
			tzClientLock = nameToTimezone(tzClientLockName)
			if tzClientLock == nil {
				log.Println("警告：客户端锁定时区名称不正确，使用客户端经纬度时区。")
				tzClientLockName = ""
			}
		}
	}
	if len(tzServerLockName) > 0 {
		if tzServerLockName == "Local" {
			tzServerLock = time.Local
		} else {
			tzServerLock = nameToTimezone(tzServerLockName)
			if tzServerLock == nil {
				log.Println("警告：XLSX 文件锁定时区名称不正确，使用客户端经纬度时区。")
				tzServerLockName = ""
			}
		}
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
		"锁定客户端时区: " + tzClientLockName,
		"锁定 XLSX 文件时区: " + tzServerLockName,
		"HTTP 接口地址: " + listenHost + uri,
		"APPID 限制: " + appids,
		"强制重新加载: " + strconv.FormatBool(forceReload),
		"详细信息: " + strconv.FormatBool(verbose),
		"反转风向数据: " + strconv.FormatBool(reverseDirection),
	}, "\n"))

	reloadXLSX()

	hostsAdd()

	initweb()

	log.Println("按 Ctrl+C 退出。")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT)
	<-quit

	exit()
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
	datas = reverseDirectionAndTimeDatas(datas)
	dataLen = len(datas)
	dataLenLen = len(strconv.Itoa(dataLen))
	if verbose {
		fmt.Println("行", fmt.Sprintf("%0*d", dataLenLen, 1), ": |", strings.Join(titles, " | "), "|")
		for i, row := range datas {
			fmt.Println(viewRow(i+2, row))
		}
	}
	log.Println("读取文件:", xlsxFilePath, "完成，数据量:", dataLen)
}

func viewRow(line int, row []string) string {
	return fmt.Sprintf("行 %0*d: | %s |", dataLenLen, line, strings.Join(row, " | "))
}

func exit() {
	log.Println("正在停止...")
	hostsRm()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("警告：强制结束服务器: %v", err)
	}
	log.Println("退出。")
	os.Exit(0)
}
