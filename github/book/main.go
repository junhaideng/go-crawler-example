package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)



var ReadmeURL = "https://raw.githubusercontent.com/guanpengchn/awesome-books/master/README.md"
var LinkPattern = regexp.MustCompile(`\((https://github.*?.pdf)\)`)
var dir = "pdf"

func init(){
	os.Mkdir(dir, 0775)
}

func main() {
	client := http.Client{}
	data, err := getHtml(client, ReadmeURL)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	// 用来限制线程数目
	var limitChan = make(chan struct{}, 10)
	for _, v := range LinkPattern.FindAllSubmatch(data, -1)[1:] {
		limitChan <- struct{}{}
		pdf_url := string(v[1])
		url, _ := url.Parse(pdf_url)
		subseq := strings.Split(url.Path, "/")
		filename := subseq[len(subseq)-1]
		go download(client, pdf_url, filename, limitChan)
	}
}

// 获取网页源码
func getHtml(client http.Client, url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// 下载文件
func download(client http.Client, url, filename string, limitChan chan struct{}) {
	now := time.Now()
	fmt.Println("正在下载: ", filename)
	defer func() {
		fmt.Printf("下载 %s 完成，用时: %.2f s\n", filename, float64(time.Since(now))/float64(time.Second))
		<- limitChan
	}()
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}	
	defer resp.Body.Close()
	file, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		fmt.Println("create file error: ", err)
		return
	}
	defer file.Close()
	var chuck [1024]byte
	for {
		n, err := resp.Body.Read(chuck[:])
		if err != nil {
			if err == io.EOF {
				file.Write(chuck[:n])
				break
			}
		}
		file.Write(chuck[:n])
	}
}
