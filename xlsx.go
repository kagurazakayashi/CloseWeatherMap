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

	for i, row := range rows {
		for j, cell := range row {
			rows[i][j] = trimExtraWhitespace(cell)
		}
	}

	return rows[0], reCalcDays(rows[1:])
}

func reverseDirectionAndTimeDatas(rows [][]string) [][]string {
	for i, row := range rows {
		// 補充相對日期
		var timeStr = row[1]
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
		var baseDayStr string = row[0]
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
		direction, err := strconv.ParseFloat(strings.ReplaceAll(row[6], ",", ""), 64)
		if err != nil {
			log.Println("错误: 无法解析风向数据:", row[6])
			continue
		}
		if reverseDirection {
			// 小於180時+180，大於180時-180
			if direction < 180 {
				direction += 180
			} else if direction > 180 {
				direction -= 180
			}
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
		var days string = row[0]
		var time string = row[1]
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

func genTime(timeData string, iTimezone *time.Location) (bool, time.Time) {
	var tz *time.Location = iTimezone
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

func isCurrentTimeInRange(currentTime, startTime, endTime time.Time) int8 {
	// fmt.Println(currentTime.Location(), startTime.Location(), endTime.Location())
	// 開始時間晚於結束時間
	if startTime.After(endTime) {
		return -2
	}
	// 當前時間早於開始時間
	if currentTime.Before(startTime) {
		return -1
	}
	// 當前時間晚於結束時間
	if currentTime.After(endTime) {
		return 1
	}
	return 0
}

func nowTimeData(uTime time.Time, iTimezone *time.Location) []string {
	// var realTime time.Time = nowTime()
	// var daysApart int = daysApart(baseDayDate, realTime) + 1
	for i, row := range datas {
		// rowDay, err := strconv.Atoi(row[0])
		// if err != nil {
		// 	continue
		// }
		if verbose {
			fmt.Println("\n" + viewRow(i+2, row))
		}

		// if rowDay > daysApart {
		// 	break
		// }
		// if rowDay != daysApart {
		// 	continue
		// }
		var useTimeZone *time.Location = iTimezone
		if convertTimeZone {
			useTimeZone = time.Local
		}
		isOK, startTime := genTime(row[1], useTimeZone)
		if !isOK {
			continue
		}
		if convertTimeZone {
			startTime = startTime.In(iTimezone)
		}
		if i == len(datas)-1 {
			fmt.Println("警告：达到数据末尾，返回最后的数据。")
			return row
		} else {
			isOK, endTime := genTime(datas[i+1][1], useTimeZone)
			if !isOK {
				continue
			}
			if convertTimeZone {
				endTime = endTime.In(iTimezone)
			}
			endTime = endTime.Add(-1 * time.Second)
			var timeRange int8 = isCurrentTimeInRange(uTime, startTime, endTime)
			isOK = timeRange == 0
			if verbose {
				var isOKs string = "否"
				if isOK {
					isOKs = "是"
				}
				// fmt.Println("本地时间", uTime.Local(), "在", startTime.Local(), "～", endTime.Local(), "区间？", isOKs)
				fmt.Println(iTimezone, "时间", uTime, "在", startTime, "～", endTime, "区间？", isOKs)
			} else {
				log.Println(viewRow(i+2, row))
				log.Println(iTimezone, "时区 UTC +", getUTCOffset(iTimezone), "当地时间", uTime.Format(timeLayout), "数据时间", startTime.Format(timeLayout))
			}
			if i == 0 && timeRange == -1 {
				log.Println("警告：未到数据开始时间，返回了第一条数据。")
				return row
			}
			if isOK {
				return row
			}
		}
	}
	return []string{}
}
