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
	f.Close()

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

func reverseDirectionAndTimeDatas(rows [][]string) [][]string {
	for i, row := range rows {
		// 補充相對日期
		var timeStr = trimExtraWhitespace(row[1])
		var timeArr []string = strings.Split(timeStr, ":")
		if len(timeArr) == 1 {
			timeArr = append(timeArr, timeArr[0])
			timeArr[0] = "0"
		} else if len(timeArr) == 0 {
			log.Println("错误: 无法解析时间数据:", timeStr)
		}
		hour, err := strconv.Atoi(timeArr[0])
		if err != nil {
			log.Println("错误: 无法解析时间数据:", timeStr)
			continue
		}
		minute, err := strconv.Atoi(timeArr[1])
		if err != nil {
			log.Println("错误: 无法解析时间数据:", timeStr)
			continue
		}
		var baseDayStr string = trimExtraWhitespace(row[0])
		baseDay, err := strconv.Atoi(baseDayStr)
		if err != nil {
			log.Println("错误: 无法解析日期数据:", baseDayStr)
			continue
		}
		var baseDaye time.Time = time.Date(baseDayDate.Year(), baseDayDate.Month(), baseDayDate.Day(), hour, minute, 0, 0, time.Local)
		if baseDay > 1 {
			baseDaye = baseDaye.AddDate(0, 0, baseDay-1)
		}
		rows[i][1] = baseDaye.Format(timeLayout)
		// 處理風向
		direction, err := strconv.ParseFloat(strings.ReplaceAll(trimExtraWhitespace(row[6]), ",", ""), 64)
		if err != nil {
			log.Println("错误: 无法解析风向数据:", row[6])
			continue
		}
		// 小於180時+180，大於180時-180
		if direction < 180 {
			direction += 180
		} else if direction > 180 {
			direction -= 180
		}
		rows[i][6] = strconv.FormatFloat(direction, 'f', -1, 64)
	}
	return rows
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
	}
	return rows
}

func genTime(timeData string) (bool, time.Time) {
	var tz *time.Location = time.Local
	if len(tzServerLockName) > 0 {
		tz = tzServerLock
	}
	gTime, err := time.ParseInLocation(timeLayout, timeData, tz)
	if err != nil {
		return false, time.Time{}
	}
	return true, gTime
	// var nowTime time.Time = nowTime()
	// if timeData == "0" {
	// 	timeData = "0:00"
	// }
	// var timeArr []string = strings.Split(timeData, ":")
	// if len(timeArr) != 2 {
	// 	return false, time.Time{}
	// }
	// startHour, err := strconv.Atoi(timeArr[0])
	// if err != nil {
	// 	return false, time.Time{}
	// }
	// startMinute, err := strconv.Atoi(timeArr[1])
	// if err != nil {
	// 	return false, time.Time{}
	// }
	// return true, time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), startHour, startMinute, 0, 0, time.Local)
}

func isCurrentTimeInRange(currentTime, startTime, endTime time.Time) bool {
	// fmt.Println(currentTime.Location(), startTime.Location(), endTime.Location())
	if startTime.After(endTime) {
		return false
	}
	if currentTime.Before(startTime) {
		return false
	}
	if currentTime.After(endTime) {
		return false
	}
	return true
}

func nowTimeData(uTime time.Time) []string {
	// var realTime time.Time = nowTime()
	// var daysApart int = daysApart(baseDayDate, realTime) + 1
	for i, row := range datas {
		// rowDay, err := strconv.Atoi(row[0])
		// if err != nil {
		// 	continue
		// }
		if verbose {
			fmt.Println("行", i, ":", row)
		}
		// if rowDay > daysApart {
		// 	break
		// }
		// if rowDay != daysApart {
		// 	continue
		// }
		isOK, startTime := genTime(row[1])
		if !isOK {
			continue
		}
		// if rowDay-1 > 0 {
		// 	startTime = startTime.AddDate(0, 0, rowDay-1)
		// }
		// if daysApart > 0 {
		// 	startTime = startTime.AddDate(0, 0, daysApart-1)
		// }
		if i == len(datas)-1 {
			fmt.Println("警告：达到数据末尾，返回最后的数据。")
			return row
		} else {
			isOK, endTime := genTime(datas[i+1][1])
			if !isOK {
				continue
			}
			// if rowDay-1 > 0 {
			// 	endTime = endTime.AddDate(0, 0, rowDay-1)
			// }
			// if daysApart > 0 {
			// 	endTime = endTime.AddDate(0, 0, daysApart-1)
			// }
			// if datas[i][0] != datas[i+1][0] {
			// 	endTime = endTime.AddDate(0, 0, 1)
			// }
			endTime = endTime.Add(-1 * time.Second)
			isOK = isCurrentTimeInRange(uTime, startTime, endTime)
			if verbose {
				var isOKs string = "否"
				if isOK {
					isOKs = "是"
				}
				fmt.Println("时间", uTime.Format(timeLayout), "在", startTime.Format(timeLayout), "～", endTime.Format(timeLayout), "区间？", isOKs)
			}
			if isOK {
				return row
			}
		}
	}
	return []string{}
}
