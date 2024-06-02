package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gofish2020/MedicalSpider/logger"
	"github.com/gofish2020/MedicalSpider/spider/brower"
	"github.com/gofish2020/MedicalSpider/spider/douyin"
	"github.com/gofish2020/MedicalSpider/utils"
	"github.com/tebeka/selenium"
)

func main() {

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
	for {

		var destUrl string
		fmt.Println("请输入用户主页url:")
		fmt.Scanln(&destUrl)
		if destUrl == "" {
			destUrl = "https://www.douyin.com/user/MS4wLjABAAAAG0itb5W4abVyFGxSNlGkKvG20un-Ix4R2UfuJUeXHo0?vid=7375750907994197286"
		} else if strings.ToLower(destUrl) == "q" {
			return
		}
		// 利用该浏览器爬取网页
		SpiderOneUser(destUrl, driver)
	}

}

func SpiderOneUser(destUrl string, driver selenium.WebDriver) (err error) {

	logger.Info("开始爬取用户主页:", destUrl)
	err = driver.Get(destUrl)
	if err != nil {
		log.Fatal("Error3:", err)
		return
	}

	logger.Info("等待扫码登录...")
	for {
		ele, err := driver.FindElement(selenium.ByClassName, "semi-button-content")
		if err != nil {
			log.Fatal("Error ele:", err)
			return err
		}

		val, _ := ele.Text()
		if !strings.Contains(val, "登录") {
			break

		}
	}

	logger.Info("扫码登录成功")

	// 获取cookies信息
	cookies, err := driver.GetCookies()
	if err != nil {
		log.Fatal("Error cookies:", err)
		return
	}

	cookieStr := []string{}
	for _, cookie := range cookies {

		if cookie.Name != "" {
			cookieStr = append(cookieStr, cookie.Name+"="+cookie.Value)
		} else {
			cookieStr = append(cookieStr, cookie.Value)
		}
	}
	cookieResult := strings.Join(cookieStr, "; ")
	douyin.SetCookie(cookieResult)
	// 模拟自动向下滚动
	for {

		ele, err := driver.FindElement(selenium.ByClassName, "B_mbw29p")
		if err == nil {
			val, _ := ele.Text()
			if strings.Contains(val, "暂时没有更多了") {
				break
			}
		}

		//fmt.Println("执行滚动")
		driver.ExecuteScriptRaw("window.scrollBy(0,1000)", nil)
		time.Sleep(500 * time.Millisecond)
	}

	logger.Info("滚动完成")

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

	logger.Info("开始视频的抓取")

	fmt.Println()

	douyinId := utils.GetDouYinId(query.Find(".TVGQz3SI").Text())

	count := 0
	query.Find(".niBfRBgX").Each(func(i int, s *goquery.Selection) {
		count++
		url, ok := s.Find("a").Attr("href")
		if ok {
			pos := strings.LastIndex(url, "/")
			awemeid := url[pos+1:]
			douyin.GetAwemeComment(awemeid, douyinId)
		}

	})
	logger.Info("完成", count, "视频的抓取")

	return nil
}
