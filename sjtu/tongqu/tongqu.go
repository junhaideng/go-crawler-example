package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
)

type conf struct{
	Data Info `json:"data"`
	Cookie string `json:"cookie"`
}

// Info .
type Info struct{
	ActID string `json:"act_id"`
	UserSignInfo json.RawMessage `json:"user_sign_info"`
	AttachID string `json:"attach_id"`
}

// 请求返回的消息
type response struct {
	Version string `json:"version"`
	Error uint `json:"error"`
	Msg string `json:"msg"`
	Code uint `json:"code"`
}


func main(){

	data, _ := ioutil.ReadFile("config.json")
	var conf = &conf{}
	json.Unmarshal(data, conf)
	//fmt.Println(string(conf.Data.UserSignInfo))
	url := "https://tongqu.sjtu.edu.cn/api/act/sign"


	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	writer.WriteField("act_id", conf.Data.ActID)
	writer.WriteField("user_sign_info", string(conf.Data.UserSignInfo))
	writer.WriteField("attach_id", conf.Data.AttachID)
	writer.Close()


	req, _ := http.NewRequest("POST", url, &buffer)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Safari/537.36 Edg/87.0.664.41")
	req.Header.Add("Origin", "https://tongqu.sjtu.edu.cn")
	req.Header.Add("Cookie", conf.Cookie)
	req.Header.Add("Host", "tongqu.sjtu.edu.cn")
	req.Header.Add("Content-Type", writer.FormDataContentType())


	var client http.Client
	counter := 1
	now := time.Now()
	for{
		fmt.Printf("正在进行第%d次尝试..\n", counter)
		counter ++
		res, err := client.Do(req)
		if err != nil {
			fmt.Println("err: ", err)
			continue
		}


		var response = &response{}
		resJSON, err := ioutil.ReadAll(res.Body)
		if err != nil{
			fmt.Println("读取返回数据失败")
			continue
		}
		res.Body.Close()

		err = json.Unmarshal(resJSON, response)
		if err != nil{
			fmt.Println("解析返回的json数据失败")
			continue
		}
		if response.Error == 0 {
			fmt.Println(response.Msg)
			break
		}else{
			fmt.Println("返回消息: ", response.Msg)
		}
	}
	fmt.Printf("共用时: %d ms", time.Since(now)/time.Millisecond)
}