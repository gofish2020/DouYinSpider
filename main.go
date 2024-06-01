package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gofish2020/MedicalSpider/logger"
	"github.com/gofish2020/MedicalSpider/spider/brower"
	"github.com/gofish2020/MedicalSpider/spider/douyin"
	"github.com/gofish2020/MedicalSpider/utils"
)

func main() {

	flag.Parse() // ./MedicalSpider xxxxx

	destUrl := "https://www.douyin.com/user/MS4wLjABAAAARDPIUU4VMK7S5kzoT5RRER2eOgqYXPi9AOfO1NI8Xb8?vid=7371468774735629587"
	if len(os.Args) > 1 {
		destUrl = os.Args[1]
	}
	// 日志初始化
	logger.Setup(&logger.Settings{
		Path:       "./logs/",
		Name:       "medical",
		Ext:        "log",
		DateFormat: utils.DateFormat,
	})
	logger.SetLoggerLevel(logger.DEBUG)

	// 启动一个浏览器
	driver, err := brower.StartChrome()
	if err != nil {
		log.Fatal("Error1:", err)
		return
	}

	defer driver.Quit()

	err = driver.Get(destUrl)
	//err = driver.Get("https://www.baidu.com/")
	if err != nil {
		log.Fatal("Error3:", err)
		return
	}

	time.Sleep(5 * time.Second)

	html, err := driver.PageSource()
	if err != nil {
		log.Fatal("Error4:", err)
		return
	}

	query, err := goquery.NewDocumentFromReader(bytes.NewBuffer([]byte(html)))
	if err != nil {
		log.Fatal("Error5:", err)
		return
	}

	query.Find(".niBfRBgX").Each(func(i int, s *goquery.Selection) {

		url, ok := s.Find("a").Attr("href")
		if ok {
			pos := strings.LastIndex(url, "/")
			awemeid := url[pos+1:]
			douyin.GetAwemeComment(awemeid)
		}

	})

	//douyin.GetAwemeComment("7242230640475786556")

}
