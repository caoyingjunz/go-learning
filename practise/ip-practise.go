package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type IpInfo struct {
	Code int    `json:"code"`
	Data IpSpec `json:"data"`
	Msg  string `json:"msg"`
}

type IpSpec struct {
	Area      string `json:"area"`
	Country   string `json:"country"`
	IspId     string `json:"isp_id"`
	QueryIp   string `json:"queryIp"`
	City      string `json:"city"`
	Ip        string `json:"ip"`
	Isp       string `json:"isp"`
	County    string `json:"county"`
	RegionId  string `json:"region_id"`
	AreaId    string `json:"area_id"`
	CountyId  string `json:"county_id"`
	Region    string `json:"region"`
	CountryId string `json:"country_id"`
	CityId    string `json:"city_id"`
}

func main() {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(fmt.Sprintf("http://ip.taobao.com/outGetIpInfo?ip=%s&accessKey=alibaba-inc", "183.183.183.183"))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var ipInfo IpInfo
	if err = json.Unmarshal(data, &ipInfo); err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("%+v", ipInfo))
}
