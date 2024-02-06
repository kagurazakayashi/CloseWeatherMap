package main

import (
	"log"
	"strings"

	"github.com/xuri/excelize/v2"
)

func loadXLSX(loginfo bool) [][]string {
	f, err := excelize.OpenFile(xlsxFilePath)
	if err != nil {
		if loginfo {
			log.Println("错误：打开文件", xlsxFilePath, "失败:", err.Error())
		}
		return [][]string{}
	}

	sheetList := f.GetSheetList()
	if loginfo {
		log.Println("表格列表:", strings.Join(sheetList, ","), "。正在加载数据表 Sheet1 ...")
	}
	if !stringInSlice("Sheet1", sheetList) {
		log.Println("错误：数据表 Sheet1 不存在。")
		return [][]string{}
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		log.Println("错误：加载数据表 Sheet1 失败。")
		return [][]string{}
	}

	rows = rows[1:]

	return rows

	// cellValue, err := f.GetCellValue("Sheet1", "B2")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Value at B2:", cellValue)
}
