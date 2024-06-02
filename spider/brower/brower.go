package brower

import (
	"fmt"
	"os"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"gopkg.in/ini.v1"
)

// 这里修改myini.ini文件中的路径为自己的路径

var chromeDriverPath = ""
var userdatadir = ""

const port = 4444

var service *selenium.Service

func init() {

	var err error

	cfg, err := ini.Load("./myini.ini")
	if err != nil {
		fmt.Printf("Fail to read file myini.ini: %v", err)
		os.Exit(1)
	}

	chromeDriverPath = cfg.Section("config").Key("chromedriverpath").String()
	userdatadir = cfg.Section("config").Key("userdatadir").String()

	opts := make([]selenium.ServiceOption, 0)
	service, err = selenium.NewChromeDriverService(chromeDriverPath, port, opts...)
	if err != nil {
		panic("启动浏览器驱动服务失败,err= " + err.Error())
	}

}

func StartChrome() (selenium.WebDriver, error) {
	// 设置Chrome选项

	caps := selenium.Capabilities{
		"browserName": "chrome",
	}

	chromeCaps := chrome.Capabilities{
		Path: "",
		ExcludeSwitches: []string{
			"enable-automation",
		},
		Args: []string{
			//"--headless", // 设置Chrome无头模式，在linux下运行，需要设置这个参数，否则会报错
			"--disable-gpu",
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36", // 模拟user-agent，防反爬
			"--user-data-dir=" + userdatadir, // 这里替换为自己的  在谷歌浏览器中输入 chrome://version
		},
	}
	//以上是设置浏览器参数
	caps.AddChrome(chromeCaps)
	return selenium.NewRemote(caps, "")
}

func Stop() error {
	return service.Stop()
}
