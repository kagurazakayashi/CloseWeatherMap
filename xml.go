package main

import "encoding/xml"

type Current struct {
	XMLName     xml.Name `xml:"current"`
	Temperature struct {
		Value string `xml:"value,attr"`
	} `xml:"temperature"`
	Humidity struct {
		Value string `xml:"value,attr"`
	} `xml:"humidity"`
	Pressure struct {
		Value string `xml:"value,attr"`
	} `xml:"pressure"`
	Wind struct {
		Speed struct {
			Value string `xml:"value,attr"`
		} `xml:"speed"`
		Direction struct {
			Value string `xml:"value,attr"`
		} `xml:"direction"`
	} `xml:"wind"`
	Weather struct {
		Number string `xml:"number,attr"`
	} `xml:"weather"`
}

func genDic(row []string) map[string]string {
	var dataDic map[string]string = map[string]string{}
	var ignoreKeys []string = []string{"day", "time"}
	for id, key := range titles {
		if stringInSlice(key, ignoreKeys) {
			continue
		}
		dataDic[key] = row[id]
	}
	return dataDic
}

func genXML(dic map[string]string) string {
	var current Current
	current.Temperature.Value = dic["temperature"]
	current.Humidity.Value = dic["humidity"]
	current.Pressure.Value = dic["pressure"]
	current.Wind.Speed.Value = dic["wind"]
	current.Wind.Direction.Value = dic["direction"]
	current.Weather.Number = dic["weather"]
	output, err := xml.MarshalIndent(current, "", "    ")
	if err != nil {
		return ""
	}
	return xml.Header + string(output)
}
