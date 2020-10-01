package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var timeout *time.Timer
var once sync.Once

const DOWNLOAD = "./download"

func init() {
	timeout = time.NewTimer(2 * time.Second)
	_, err := os.Stat(DOWNLOAD)
	if os.IsNotExist(err) {
		os.Mkdir(DOWNLOAD, 0666)
	}
}

type info struct {
	link string
	name string
}

func main() {
	var urlChan = make(chan string, 10)
	// PPT详情页，中间有点击下载按钮
	var pptPageChan = make(chan string, 100)
	// PPT下载页
	var hrefPageChan = make(chan info, 100)
	// 下载链接
	var linkLPageChan = make(chan info, 100)
	// 是否已经完成，这里使用一个超时来判断
	// 严格来说并不是很好
	// 因为网络比较差的时候也会被判断为完成
	// 但是设置时间长一点也没有很大关系
	// 如果在这么长的时间里没有完成，那么肯定是全部下载完了
	var doneChan = make(chan bool, 1)
	var lastPage = make(chan int, 1)

	urlChan <- "http://www.ypppt.com/moban/"
	go GetHTML(urlChan, pptPageChan, lastPage)
	go generateUrl(lastPage, urlChan)
	go GetDownloadPageHref(pptPageChan, hrefPageChan)
	go GetDownloadLink(hrefPageChan, linkLPageChan)
	go downloadFile(linkLPageChan, doneChan)

	for {
		select {
		case <-doneChan:
			fmt.Println("Done")
			return
		}
	}
}

func generateUrl(lastPage <-chan int, urlChan chan<- string) {
	page := <-lastPage
	for i := 2; i <= page; i++ {
		urlChan <- fmt.Sprintf("http://www.ypppt.com/moban/list-%d.html", i)
	}
	close(urlChan)
}

// 获取每一个PPT的详情界面
// 但是该界面并不是下载PPT的界面
// 下载PPT的界面需要到另外一个界面中
// 这个界面有一个下载按钮通往另一个界面中
func GetHTML(urlChan <-chan string, c chan<- string, lastPage chan<- int) {
	for {
		select {
		case url := <-urlChan:
			fmt.Println(url)
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("Err: ", err)
			}
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Err: ", err)
			}
			resp.Body.Close()
			ulPattern := regexp.MustCompile(`(?s)<ul.*?class="posts clear">(.*?)</ul>`)
			list := ulPattern.FindSubmatch(data)

			hrefPattern := regexp.MustCompile(`(?s)<a href="(.*?)".*?class="p-title".*?</a>`)
			page := regexp.MustCompile(`下一页</a><a.*?href="list-(.*?).html">末页</a>`)
			matches := hrefPattern.FindAllSubmatch(list[1], -1)
			// 获取末尾页
			once.Do(func() {
				p, _ := strconv.Atoi(string(page.FindSubmatch(data)[1]))
				lastPage <- p
				close(lastPage)
			})
			for _, value := range matches {
				c <- "http://www.ypppt.com" + string(value[1])
			}
		}
	}
}

// 获取下载文件的那个界面
// 这个界面中存在文件的下载连接
// 这个函数返回的就是点击 《点击下载》按钮之后跳转的url界面
func GetDownloadPageHref(pptPageChan <-chan string, hrefChan chan<- info) {
	for {
		select {
		case page := <-pptPageChan:
			resp, err := http.Get(page)
			if err != nil {
				fmt.Println("Err: ", err)
				continue
			}
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Err: ", err)
				continue
			}
			resp.Body.Close()

			hrefPattern := regexp.MustCompile(`<a href="(.*?)" rel="nofollow" class="down-button".*?</a>`)
			namePattern := regexp.MustCompile(`(?s)<div class="infoss">.*?<h1>(.*?)</h1>`)
			hrefChan <- info{
				link: "http://www.ypppt.com" + string(hrefPattern.FindSubmatch(data)[1]),
				name: string(namePattern.FindSubmatch(data)[1]) + ".zip",
			}
		}
	}
}

// 然后我们需要从这个下载界面中抽取出来下载的连接
func GetDownloadLink(hrefChan <-chan info, linkChan chan<- info) {
	for {
		select {
		case href := <-hrefChan:
			resp, err := http.Get(href.link)
			if err != nil {
				fmt.Println("Err: ", err)
				continue
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Err: ", err)
				continue
			}
			resp.Body.Close()

			// 我们只需要取第一个url即可
			downloadLinkPattern := regexp.MustCompile(`(?s)<ul class="down clear">.*?<li><a href="(.*?)".*?>.*?</ul>`)

			link := string(downloadLinkPattern.FindSubmatch(data)[1])
			if !strings.HasPrefix(link, "http") && len(link) > 8 {
				link = "http://www.youpinppt.com" + link[8:]
			} else if len(link) < 8 {
				continue
			}
			linkChan <- info{
				link: link,
				name: href.name,
			}

		}
	}
}

// 下载文件
func downloadFile(linkChan <-chan info, done chan<- bool) {
	for {
		select {
		case linkInfo := <-linkChan:
			// 下载文件执行
			err := _download(linkInfo)
			timeout.Reset(10 * time.Second)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case <-timeout.C:
			done <- true
		}
	}
}

func _download(info info) error {
	resp, err := http.Get(info.link)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Printf("开始下载：%s, link: %s \n", info.name, info.link)
	file, err := os.Create(filepath.Join(DOWNLOAD, info.name))
	file.Write(data)
	defer file.Close()
	if err != nil {
		return err
	}
	return nil
}
