package brower

import (
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const port = 4444

var service *selenium.Service

func init() {
	var err error
	opts := make([]selenium.ServiceOption, 0)
	service, err = selenium.NewChromeDriverService("/usr/local/bin/chromedriver", port, opts...)
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
		},
	}
	//以上是设置浏览器参数
	caps.AddChrome(chromeCaps)
	return selenium.NewRemote(caps, "")
}

func Stop() error {
	return service.Stop()
}
