package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var client http.Client

// https://us.trip.com/travel-guide/beijing/the-palace-museum-75595/

const template = `{
    "poiId": 75595,  // 每个地点对应的id
    "locale": "en-US",
    "pageIndex": %d,  // 修改该值
    "pageSize": 5,
    "head": {
        "cver": "3.0",
        "cid": "",
        "extension": [
            {
                "name": "locale",
                "value": "en-US"
            },
            {
                "name": "platform",
                "value": "Online"
            },
            {
                "name": "currency",
                "value": "USD"
            }
        ]
    }
}`

const url = "https://www.trip.com/restapi/soa2/19707/getReviewSearch"

// 自定义设置，最好可以自己到响应的网站中获取
// 当然也可以爬取时获取
const PageNum = 100

type Comment struct {
	Author  string `json:"author"`
	Content string `json:"content"`
}

func main() {
	file, err := os.Create("trip.json")
	if err != nil {
		fmt.Println("Create file error: ", err)
		return
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("  ", "    ")

	var postChan = make(chan io.Reader, 10)
	var comments []Comment

	go func() {
		for i := 1; i < PageNum; i++ {
			postChan <- strings.NewReader(fmt.Sprintf(template, i))
		}
		close(postChan)
	}()

	for data := range postChan {
		req, err := http.NewRequest("POST", url, data)
		if err != nil {
			fmt.Println("NewRequest error: ", err)
			continue
		}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Do error: ", err)
			continue
		}
		d, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("ReadAll error: ", err)
			continue
		}
		resp.Body.Close()
		reader, err := simplejson.NewJson(d)
		if err != nil {
			fmt.Println("NewJson error: ", err)
			continue
		}
		reviewList := reader.Get("reviewList")
		list, err := reviewList.Array()
		if err != nil {
			//fmt.Println(string(d))
			fmt.Println("array error: ", err)
			continue
		}
		for i := 0; i < len(list); i++ {
			item := reviewList.GetIndex(i)
			comments = append(comments, Comment{
				// MustString 指定如果对应的值为空或者不存在时的默认值
				Author:  item.Get("content").MustString("no content"),
				Content: item.Get("username").MustString("no username"),
			})
		}
	}
	if err:=encoder.Encode(&comments); err!= nil{
		fmt.Println("Encode error: ", err)
	}
}
