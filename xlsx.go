package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func loadXLSX(loginfo bool) ([]string, [][]string) {
	f, err := excelize.OpenFile(xlsxFilePath)
	if err != nil {
		if loginfo {
			log.Println("错误：打开文件", xlsxFilePath, "失败:", err.Error())
		}
		return []string{}, [][]string{}
	}

	sheetList := f.GetSheetList()
	if loginfo {
		log.Println("表格列表:", strings.Join(sheetList, ","), "。正在加载数据表 Sheet1 ...")
	}
	if !stringInSlice("Sheet1", sheetList) {
		log.Println("错误：数据表 Sheet1 不存在。")
		return []string{}, [][]string{}
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		log.Println("错误：加载数据表 Sheet1 失败。")
		return []string{}, [][]string{}
	}

	// cellValue, err := f.GetCellValue("Sheet1", "B3")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Value at B3:", cellValue)

	// for i, row := range rows {
	// 	fmt.Println(i, row)
	// }

	return rows[0], reCalcDays(rows[1:])
}

func reCalcDays(rows [][]string) [][]string {
	var isFirst bool = true
	var oldHour int = -1
	var nDay int = 0
	for i, row := range rows {
		var days string = trimExtraWhitespace(row[0])
		var time string = trimExtraWhitespace(row[1])
		var timeArr []string = strings.Split(time, ":")
		day, err := strconv.Atoi(days)
		if err != nil {
			continue
		}
		hour, err := strconv.Atoi(timeArr[0])
		if err != nil {
			continue
		}
		if isFirst {
			oldHour = hour
			nDay = day
			isFirst = false
		} else {
			if hour < oldHour {
				nDay++
			}
			oldHour = hour
		}
		rows[i][0] = strconv.Itoa(nDay)
		timeArr[0] = strconv.Itoa(hour)
		rows[i][1] = strings.Join(timeArr, ":")
		fmt.Println(rows[i])
	}
	return rows
}

func genTime(timeData string) (bool, time.Time) {
	var nowTime time.Time = time.Now()
	var timeArr []string = strings.Split(timeData, ":")
	if len(timeArr) != 2 {
		return false, time.Time{}
	}
	startHour, err := strconv.Atoi(timeArr[0])
	if err != nil {
		return false, time.Time{}
	}
	startMinute, err := strconv.Atoi(timeArr[1])
	if err != nil {
		return false, time.Time{}
	}
	return true, time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), startHour, startMinute, 0, 0, time.Local)
}

func nowTimeData(nowTime time.Time) []string {
	var daysApart int = daysApart(baseDayDate, nowTime) + 1
	for i, row := range datas {
		rowDay, err := strconv.Atoi(row[0])
		if err != nil {
			continue
		}
		// fmt.Println("序号", i, "日期", daysApart)
		if rowDay != daysApart {
			continue
		}
		isOK, startTime := genTime(row[1])
		if !isOK {
			continue
		}
		if i == len(datas)-1 {
			return row
		} else {
			isOK, endTime := genTime(datas[i+1][1])
			if !isOK {
				continue
			}
			endTime = endTime.Add(-1 * time.Nanosecond)
			// fmt.Println("时间区间", startTime, "到", endTime)
			if nowTime.After(startTime) && nowTime.Before(endTime) {
				return row
			}
		}
	}
	return []string{}
}
