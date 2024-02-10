package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

func initweb() bool {
	http.HandleFunc(uri, handlerRoot)
	log.Println("启动 HTTP 服务: http://" + listenHost + uri)
	err := http.ListenAndServe(listenHost, nil)
	if err != nil {
		log.Println("错误：无法启动 HTTP 服务:", err.Error())
		return false
	}
	return true
}

func handlerRoot(w http.ResponseWriter, r *http.Request) {
	// ?lat=48.0061&lon=0.1996&APPID=32bit&mode=xml&units=metric
	log.Println("<- " + r.URL.Path + "?" + r.URL.RawQuery)
	var info string = ""
	var appid string = r.FormValue("APPID")
	if !chkAPPID(appid) {
		w.WriteHeader(403)
		info = "错误：无效的 APPID: " + appid
		log.Println(info)
		w.Write([]byte(info))
		return
	}
	var mode string = r.FormValue("mode")
	if mode != "xml" {
		info = "错误：不支持的输出格式: " + mode
		log.Println(info)
		w.WriteHeader(400)
		w.Write([]byte(info))
		return
	}
	var units string = r.FormValue("units")
	if units != "metric" {
		info = "错误：不支持的单位格式: " + units
		log.Println(info)
		w.WriteHeader(400)
		w.Write([]byte(info))
		return
	}
	var lat string = r.FormValue("lat")
	latN, err := strconv.ParseFloat(lat, 64)
	if err != nil || latN < -90 || latN > 90 {
		info = "错误：无效的纬度: " + lat
		log.Println(info)
		w.WriteHeader(400)
		w.Write([]byte(info))
		return
	}
	var lon string = r.FormValue("lon")
	lonN, err := strconv.ParseFloat(lon, 64)
	if err != nil || lonN < -180 || lonN > 180 {
		info = "错误：无效的经度: " + lon
		log.Println(info)
		w.WriteHeader(400)
		w.Write([]byte(info))
		return
	}
	var date string = r.FormValue("date")
	if len(date) != 0 && len(date) != 8 {
		info = "错误：无效的日期: " + date
		log.Println(info)
		w.WriteHeader(400)
		w.Write([]byte(info))
		return
	}
	if len(date) == 8 {
		genBaseDay(date)
	}
	if forceReload {
		reloadXLSX()
	}
	var row []string = nowTimeData(time.Now())
	if len(row) == 0 {
		info = "错误：没有查询到数据。"
		log.Println(info)
		w.WriteHeader(404)
		w.Write([]byte(info))
		return
	}
	var response string = genXML(genDic(row))
	log.Println("->", response)
	w.Write([]byte(response))
}

func chkAPPID(appid string) bool {
	var appidLen int = len(appids)
	if appidLen == 0 {
		return true
	}
	if appidLen != 32 {
		return false
	}
	if checkIfStringExists(appids, appid) {
		return true
	}
	return false
}
