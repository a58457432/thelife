package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	url  = "https://restapi.amap.com/v3/weather/weatherInfo?"
	key  = "024ede8cfe5a40c29fdfa4c4863c13d1"
	city = "310000" // 城市编码
	ext  = "all"    // all返回预报天气，base返回实况天气
	WeatTemplateID = "hSICwW0B4F6uzpwfZhF0lHu55rsS7NW08pGBO2PrhXU" //templateid
)


type Weather struct {
	Status    string     `json:"status"返回状态`
	Count     string     `json:"count"返回结果总条数`
	Info      string     `json:"info"返回的状态信息`
	Infocode  string     `json:"infocode"返回状态说明`
	Forecasts []Forecast `json:"forecasts"预报天气信息数据`
}
type Forecast struct {
	City       string `json:"city"城市名称`
	Adcode     string `json:"adcode"城市编码`
	Province   string `json:"province"省份`
	Reporttime string `json:"reporttime"预报时间`
	Casts      []Cast `json:casts预报数据`
}
type Cast struct {
	Date         string `json:"date"日期`
	Week         string `json:"week"星期`
	Dayweather   string `json:"dayweather"白天天气`
	Nightweather string `json:"nightweather"晚上天气`
	Daytemp      string `json:"daytemp"白天温度`
	Nighttemp    string `json:"nighttemp"晚上温度`
	Daywind      string `json:"daywind"白天风向`
	Nightwind    string `json:"nightwind"晚上风向`
	Daypower     string `json:"daypower"白天风力`
	Nightpower   string `json:"nightpower"晚上风力`
}

type NowWeather struct {
	Status    string     `json:"status"返回状态`
	Count     string     `json:"count"返回结果总条数`
	Info      string     `json:"info"返回的状态信息`
	Infocode  string     `json:"infocode"返回状态说明`
	Forecasts []Forecast `json:"forecasts"预报天气信息数据`
	Lives	[]Live
}


type Live struct {
	City       string `json:"city"城市名称`
	Adcode     string `json:"adcode"城市编码`
	Province   string `json:"province"省份`
	Weather	   string `json "weather"`
	Temperature string `json "temperature" 平均温度`
	Winddirection string `json "winddirection"`
	Windpower	string	`json "windpower"`
	Humidity	string	`json "humidity"`
	Reporttime string `json:"reporttime"预报时间`
	Casts      []Cast `json:casts预报数据`
}

type Token struct {
	Access_token string `json "access_token"`
	Expires_in	int `json "expires_in"`

}

type Openidlist struct {
	Total int `json "total"`
	Count int	`json "count"`
	Data []Dat `json "data"`
}

type Dat struct {
	Openid string  `json "openid"`
	Next_Openid string `json "nextopenid"`
}

// 网络请求
func doHttpGetRequest(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}
	return string(body), nil
}

// 获取天气信息
func getWeather() (string, string, string ,error) {
	var data Weather
	var fore Forecast
	var cast Cast
	var str string
	urlInfo := url + "city=" + city + "&key=" + key+ "&extensions=" + ext
//	fmt.Println(urlInfo)
	rlt, err := doHttpGetRequest(urlInfo)
//	println(rlt)
	if err != nil {
		return "网络访问失败", "","", err
	}

	err = json.Unmarshal([]byte(rlt), &data)
	if err != nil {
		return "json数据解析失败", "","", err
	}
	fore = data.Forecasts[0]
	output := fore.City + "预报时间：" + fore.Reporttime + "\n"
	for i := 0; i < len(fore.Casts); i++ {
		cast = fore.Casts[i]
		str += "日期:" + cast.Date + "\t星期" + NumToStr(cast.Week) +
			"\n白天：【天气：" + cast.Dayweather + "\t温度：" + cast.Daytemp + "\t风向" + cast.Daywind + "\t风力：" + cast.Daypower + "】" +
			"\n夜晚：【天气：" + cast.Nightweather + "\t温度：" + cast.Nighttemp + "\t风向" + cast.Nightwind + "\t风力：" + cast.Nightpower + "】" + "\n"
	}
	subject := verify(fore.Casts[0].Dayweather, fore.Casts[0].Nightweather) + "\n" +"平均温度:   " + getnowWeather()


	fmt.Println(subject)
	fmt.Println(output)
//	fmt.Println(str)
	dtime := fore.Reporttime
	theday := getnowWeather()
	wea := verify(fore.Casts[0].Dayweather, fore.Casts[0].Nightweather)

	return wea, dtime, theday, nil
}

func getnowWeather() (string) {
	var data NowWeather
	var liv	Live
	urlInfo := url + "city=" + city + "&key=" + key
//	fmt.Println(urlInfo)
	rlt, err := doHttpGetRequest(urlInfo)
	if err != nil {
		return "网络访问失败"
	}
//	print(rlt)
	err = json.Unmarshal([]byte(rlt), &data)
	if err != nil {
		return "json数据解析失败"
	}
	liv = data.Lives[0]
	NowTemper := liv.Temperature
	return NowTemper
}


func NumToStr(str string) string {
	switch str {
	case "1":
		return "一"
	case "2":
		return "二"
	case "3":
		return "三"
	case "4":
		return "四"
	case "5":
		return "五"
	case "6":
		return "六"
	case "7":
		return "日"
	default:
		return ""
	}
}

func verify(dayweather, nightweather string) string {
	var sub string
	rain := "雨"
	snow := "雪"
	cloudy := "云"
	yin := "阴"
	sub = ""
	if strings.Contains(dayweather, rain) || strings.Contains(nightweather, rain) {
		sub = sub + "今天将降雨，出门请别忘带伞"
	}
	if strings.Contains(dayweather, snow) || strings.Contains(nightweather, snow) {
		sub = sub + "下雪了"
	}
	if strings.Contains(dayweather, cloudy) || strings.Contains(nightweather, cloudy) {
		sub = sub + "多云天气啊~"
	}
	if strings.Contains(dayweather, yin) || strings.Contains(nightweather, yin) {
		sub = sub + "阴天"
	}
	return sub
}





// 定时结算（一天发一次）
/**
func TimeSettle() {
	d := time.Duration(time.Minute)
	t := time.NewTicker(d)
	defer t.Stop()
	for {
		currentTime := time.Now()
		if currentTime.Hour() == 8 { // 8点发送
			sendinfo()
			time.Sleep(time.Hour)
		}
		<-t.C
	}
}
**/




// http://mp.weixin.qq.com/debug/cgi-bin/sandboxinfo?action=showinfo&t=sandbox/index
//https://blog.csdn.net/u012140251/article/details/89529540
func Getaccesstoken() string {
	var tok Token
	var APPID string = "wx8f6147f41f64fd7e"
	var APPSECRET string = "553492ef7d951d399a7fb18f11174269"

	urlInfo := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", APPID, APPSECRET)
	fmt.Println(urlInfo)
	rlt, err := doHttpGetRequest(urlInfo)
	if err != nil {
		return "网络访问失败"
	}

	err = json.Unmarshal([]byte(rlt), &tok)
	if err != nil {
		return "解析失败"
	}
//	fmt.Println(tok.Access_token)

	return tok.Access_token
}

//获取关注者列表

func GetFlist(access_token string) []gjson.Result {
	urlInfo := "https://api.weixin.qq.com/cgi-bin/user/get?access_token=" + access_token + "&next_openid="
	rlt, err := doHttpGetRequest(urlInfo)
	if err != nil {
		fmt.Println(err)
	}

//	fmt.Println(urlInfo)
//	fmt.Println(rlt)
	flist := gjson.Get(string(rlt), "data.openid").Array()
	return flist
}

//推送模板
func templatepost(access_token string, reqdata string, fxurl string, templateid string, openid string) {
	url := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + access_token

	reqbody := "{\"touser\":\"" + openid + "\", \"template_id\":\"" + templateid + "\", \"url\":\"" + fxurl + "\", \"data\":" + reqdata + "}"
//	reqbody := "{\"touser\":\"" + openid + "\", \"template_id\":\"" + templateid + "\", \"data\": " + reqdata + "}"
	fmt.Println("#############################")
	fmt.Println(reqbody)
	resp, err := http.Post(url,
		"application/x-www-form-urlencoded",
		strings.NewReader(string(reqbody)))
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("************************")
	fmt.Println(string(body))
}

func sendweather(access_token, openid string) {
	city := "上海市"
	air_tips := "哈哈, 阿妹小~"

	wea, day, tem, _ := getWeather()
	if city == "" || day == "" || wea == "" || tem == ""|| air_tips == "" {
		return
	}
	reqdata := "{\"city\":{\"value\":\"城市:  " + city + "\", \"color\":\"#0000CD\"}, \"day\":{\"value\":\"" + day + "\", \"color\":\"#BA55D3\"}, \"wea\":{\"value\":\"天气:  " + wea + "\", \"color\":\"#339933\"}, \"tem1\":{\"value\":\"平均温度:  " + tem + "°C" + "\", \"color\":\"#FF6666\"}, \"air_tips\":{\"value\":\"tips:  " + air_tips + "\", \"color\":\"#CCCC99\"}}"
//	fmt.Println(reqdata)
	templatepost(access_token, reqdata, "https://www.baidu.com/", WeatTemplateID, openid)

}


//https://api.weixin.qq.com/cgi-bin/user/get?access_token=39_VuB-WZCQDV-IcF7ceN1Lk8oofFYMhxgfn_EL_n0u32Slnb8GpBe_J3CxoK5sXCmdjdyAmQ-iCWw0yJrSJHmpYa_ls-TP1LpBdOKIHk9ThzoabZojtVYIECxWbh0FUbT2GHzgWAZ9G4VEJNGZAWVfAHAIED&next_openid=o3ush6MXVkQWgt5SoXSPTWfCauLY

func main() {

	getWeather()
//	getnowWeather()
//	getaccesstoken()
	accesstoken := Getaccesstoken()
	fmt.Println(accesstoken)
	flists := GetFlist(accesstoken)
	fmt.Println(flists)
	sendweather(accesstoken, "o3ush6MXVkQWgt5SoXSPTWfCauLY")
//	sendweather(accesstoken, "o3ush6B2h4fXMXbUZg2Ojrh6n9zo")
}
