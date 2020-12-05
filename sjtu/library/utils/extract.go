package utils

import (
	"regexp"
	"bytes"
)

var Tr = regexp.MustCompile(`(?s)<tr valign=baseline>(.*?)</tr>`)
var Td = regexp.MustCompile(`(?s)<td.*?>(.*?)</td>`)

// 书本信息
type BookInfo struct {
	BookID string `json:"book_id,omitempty"`
	Title  Data   `json:"title,omitempty"`
	Author string `json:"author,omitempty"`
	Year   string `json:"year,omitempty"`
	Info   Data   `json:"info,omitempty"`
	Type   string `json:"type,omitempty"`
}

// 有一些字段，比如说书名以及馆名中同时含有文本和链接
type Data struct {
	URL  string `json:"url,omitempty"`
	Name string `json:"name,omitempty"`
}

// 从源代码中提取到需要的信息数据
func Extract(html []byte)[]BookInfo{
	info := make([]BookInfo, 0)
	// 找到所有的 tr
	trs := Tr.FindAllSubmatch(html, -1)
	for _, tr := range trs {
		// tr[0] -> 匹配的字符串
		// tr[1] -> 匹配的子串
		// 找到tr中所有的td，中间包含需要的内容
		td_list := Td.FindAllSubmatch(tr[1], -1)
		var data [][]byte
		for _, td := range td_list[2:8]{
			data = append(data, td[1])
		}
		info = append(info, process(data))
	}
	return info
}

// 处理td保存到结构体中
func process(tds [][]byte)BookInfo{
	var info = BookInfo{}
	info.BookID = string(bytes.TrimSpace(tds[0]))
	info.Title = getLinkAndTitle(tds[1])
	info.Author = string(bytes.TrimSpace(tds[2]))
	info.Year = string(bytes.TrimSpace(tds[3]))
	info.Info = getLinkAndTitle(tds[4])
	info.Type = string(bytes.TrimSpace(tds[5]))
	return info
}

var LinkTitle = regexp.MustCompile(`(?si)<a href=(.*?)>(.*?)</a>`)

func getLinkAndTitle(td []byte)Data{
	link := LinkTitle.FindSubmatch(td)
	data := Data{}
	if len(link) == 3{
		data.URL = string(link[1])
		data.Name = string(bytes.Replace(link[2], []byte(" "), []byte(""), -1))
	}
	return data
}