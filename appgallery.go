// 华为应用商店钉钉评论获取
// 只爬取一页内容，其他内容可以通过构造reqPageNum实现
// https://appgallery.huawei.com/#/app/C100137037
package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)


type Comment struct {
	Username string `json:"username"`
	Phone string `json:"phone"`
	Content string `json:"content"`
	Stars string `json:"stars"`
	Version string `json:"version"`
}

func main() {
	var query = url.Values{}
	var comments []Comment
	query.Add("method", "internal.user.commenList3")
	query.Add("serviceType", "20")
	query.Add("reqPageNum", "1")
	query.Add("maxResults", "25")
	query.Add("appid", "C100137037")
	query.Add("version", "10.0.0")
	query.Add("zone", "")
	query.Add("locale", "en_US")
	resp, err := http.Get("https://web-drcn.hispace.dbankcloud.cn/uowap/index?"+query.Encode())
	if err != nil {
		fmt.Println(err)
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		fmt.Println(err)
		return
	}
	reader, err := simplejson.NewJson(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	list, err := reader.Get("list").Array()
	if err != nil {
		fmt.Println(err)
		return
	}
	for i:=0; i< len(list); i++{
		comments = append(comments, Comment{
			Username: reader.Get("list").GetIndex(i).Get("accountName").MustString("anonymous"),
			Phone:    reader.Get("list").GetIndex(i).Get("phone").MustString("unknown"),
			Content:  reader.Get("list").GetIndex(i).Get("commentInfo").MustString("no content"),
			Stars:    reader.Get("list").GetIndex(i).Get("stars").MustString("unknown"),
			Version:  reader.Get("list").GetIndex(i).Get("versionName").MustString("unknown"),
		})
	}

	file, err := os.Create("appgallery.json")
	encoder := json.NewEncoder(file)
	encoder.SetIndent("  ", "    ")
	if err:=encoder.Encode(&comments); err!= nil{
		fmt.Println(err)
	}
}


