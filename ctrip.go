package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

var client http.Client

// https://us.trip.com/travel-guide/beijing/the-palace-museum-75595/

const template = `
{
    "arg": {
        "resourceId": 229,
        "resourceType": 11,
        "pageIndex": %d,
        "pageSize": 10,
        "sortType": 3,
        "commentTagId": 0,
        "collapseType": 1,
        "channelType": 7,
        "videoImageSize": "700_392",
        "starType": 0
    },
    "head": {
        "cid": "09031113313834517342",
        "ctok": "",
        "cver": "1.0",
        "lang": "01",
        "sid": "8888",
        "syscode": "09",
        "auth": null,
        "extension": [
            {
                "name": "protocal",
                "value": "https"
            }
        ]
    },
    "contentType": "json"
}`
const url = "https://m.ctrip.com/restapi/soa2/13444/json/getCommentCollapseList"

var buff bytes.Buffer
var m sync.Mutex

type Comment struct {
	Author  string `json:"author"`
	Content string `json:"content"`
}

func main() {
	file, err := os.Create("ctrip.json")
	if err != nil {
		fmt.Println("Create file error: ", err)
		return
	}
	defer file.Close()

	var comments []Comment
	encoder := json.NewEncoder(file)
	encoder.SetIndent("  ", "    ")

	var postChan = make(chan io.Reader, 10)
	go func() {
		for i := 0; i < 10; i++ {
			m.Lock()
			buff.Reset()
			buff.WriteString(fmt.Sprintf(template, i))
			m.Unlock()
			postChan <- strings.NewReader(buff.String())
		}
		close(postChan)
	}()

	for data := range postChan {
		req, err := http.NewRequest("POST", url, data)

		if err != nil {
			fmt.Println("NewRequest: ", err)
			continue
		}
		resp, err := client.Do(req)

		if err != nil {
			fmt.Println("Do: ", err)
			continue
		}

		reader, err := simplejson.NewFromReader(resp.Body)
		resp.Body.Close()

		if err != nil {
			fmt.Println("NewFromReader: ", err)
			fmt.Println(resp.StatusCode)
			continue
		}
		items := reader.Get("result").Get("items")
		array, err := items.Array()
		if err != nil {
			fmt.Println("Get: ", err)
			continue
		}
		for i := 0; i < len(array); i++ {
			item := items.GetIndex(i)
			m.Lock()
			comments = append(comments, Comment{
				Author:  item.Get("userInfo").Get("userNick").MustString("anonymous"),
				Content: item.Get("content").MustString(""),
			})
			m.Unlock()
		}

	}
	encoder.Encode(&comments)
}
