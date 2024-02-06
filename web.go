package main

import (
	"log"
	"net/http"
	"strconv"
)

func initweb() bool {
	http.HandleFunc(uri, handlerRoot)
	err := http.ListenAndServe(listenHost, nil)
	if err != nil {
		log.Println("错误：无法启动 HTTP 服务:", err.Error())
		return false
	}
	return true
}

func handlerRoot(w http.ResponseWriter, r *http.Request) {
	// ?lat=48.0061&lon=0.1996&APPID=32bit&mode=xml&units=metric
	log.Println(r.URL.Path + "?" + r.URL.RawQuery)
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
	w.Write([]byte("OK"))
}

func chkAPPID(appid string) bool {
	if len(appids) == 0 {
		return true
	}
	if checkIfStringExists(appids, appid) {
		return true
	}
	return false
}
