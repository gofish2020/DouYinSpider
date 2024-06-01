package douyin

import (
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/gofish2020/MedicalSpider/utils"
	"github.com/gofish2020/gojson"
	"github.com/klauspost/compress/zstd"
)

const userAwemeListUrl = "https://www.douyin.com/aweme/v1/web/aweme/post/"

func GetUserAwemeList(userUrl string, cursor string) (result *gojson.Json, err error) {

	userResult, err := url.Parse(strings.TrimSpace(userUrl))
	if err != nil {
		return
	}

	pos := strings.LastIndex(userResult.Path, "/")
	// 设置请求参数

	// &a_bogus=EfmMBfhDDDfikDyD56xLfY3q6IF3YpnK0trEMD2fexVpKy39HMOj9exoIb4vWnjjLG%2FlIeLjy4hSY3qMxQVrA3vX9WEKlIOp-g00tFcQ5xSSs1XHCL0gJUvqmkt5SFn2RkrUrO78oiKrFmw0A2Fe-7qvyhnFwo8sNikE
	var userAwemeParamStr = "device_platform=webapp&aid=6383&channel=channel_pc_web&sec_user_id=" + userResult.Path[pos+1:] + "&max_cursor=" + cursor + "&locate_item_id=" + userResult.Query().Get("vid") + "&locate_query=false&show_live_replay_strategy=1&need_time_list=1&time_list_query=0&whale_cut_token=&cut_version=1&count=18&publish_video_strategy_type=2&update_version_code=170400&pc_client_type=1&version_code=290100&version_name=29.1.0&cookie_enabled=true&screen_width=1512&screen_height=982&browser_language=zh-CN&browser_platform=MacIntel&browser_name=Chrome&browser_version=124.0.0.0&browser_online=true&engine_name=Blink&engine_version=124.0.0.0&os_name=Mac+OS&os_version=10.15.7&cpu_core_num=8&device_memory=8&platform=PC&downlink=1.5&effective_type=3g&round_trip_time=600&webid=7233458391619110455&verifyFp=verify_lw53dly9_Sfb80Jjv_JIMT_4YKv_BVVq_vYe0bPggLOVK&fp=verify_lw53dly9_Sfb80Jjv_JIMT_4YKv_BVVq_vYe0bPggLOVK&msToken=vexPHKca80ZUYpHNS5f3Qe3HvmHraLpgWFwASagk3-qtm34Lv_BspspgWGd9OLYkERhk60fEnTUjVOrZYWx2ig3jxjlUAG4ezuOgznBZ41xtB1F-E4dCmaeBk8KSoWktkg=="

	abogus := utils.GetAbogus(userAwemeParamStr, userAgent)

	userAwemeParamStr += "&a_bogus=" + url.QueryEscape(abogus)
	// 定义req对象
	client := &http.Client{}
	req, err := http.NewRequest("GET", userAwemeListUrl+"?"+userAwemeParamStr, nil)
	if err != nil {
		return
	}

	// 设置请求头
	headerSlices := strings.Split(headerStr, "\n")
	for _, headerStr := range headerSlices {
		rawHeader := strings.Split(headerStr, ": ")
		req.Header.Set(rawHeader[0], " "+rawHeader[1])
	}

	req.Header.Set("Host", "www.douyin.com")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", cookieStr)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	encodingStr := resp.Header.Get("Content-Encoding")

	var reader io.Reader

	if encodingStr == "br" {
		reader = brotli.NewReader(resp.Body)

	} else if encodingStr == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return
		}
	} else if encodingStr == "deflate" {
		zr := flate.NewReader(resp.Body)
		defer zr.Close()
		reader = zr
	} else if encodingStr == "zstd" {
		d, err := zstd.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer d.Close()
		reader = d
	} else {
		reader = resp.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return
	}

	var js gojson.Json
	err = js.LoadString(string(body))
	if err != nil {
		return
	}

	// 数据错误
	if js.Get("status_code").Int() != 0 {
		err = fmt.Errorf("status_code 不为0")
		return
	}

	return &js, nil
}
