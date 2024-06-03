package douyin

import (
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/gofish2020/MedicalSpider/csv"
	"github.com/gofish2020/MedicalSpider/logger"
	"github.com/gofish2020/MedicalSpider/utils"
	"github.com/gofish2020/gojson"
	"github.com/klauspost/compress/zstd"
)

/*
获取评论
*/

var cookieStr string

func init() {

	// cookie, err := os.ReadFile("./cookies.txt")
	// if err != nil {
	// 	fmt.Printf("Fail to read file: %v", err)
	// 	os.Exit(1)
	// }

	// cookieStr = string(cookie)
}

func SetCookie(cookie string) {
	cookieStr = cookie
}

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"

const headerStr = `Accept: application/json, text/plain, */*
Accept-Language: zh-CN,zh;q=0.9,en;q=0.8
Accept-Encoding:  gzip, deflate, br, zstd
Priority: u=1, i
Referer: https://www.douyin.com/
Sec-Ch-Ua: "Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"
Sec-Ch-Ua-Mobile: ?0
Sec-Ch-Ua-Platform: "macOS"
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: same-origin`

const commentURL = "https://www.douyin.com/aweme/v1/web/comment/list/"

const paramStr = "device_platform=webapp&aid=6383&channel=channel_pc_web&item_type=0&insert_ids=&whale_cut_token=&cut_version=1&rcFT=&update_version_code=170400&pc_client_type=1&version_code=170400&version_name=17.4.0&cookie_enabled=true&screen_width=1512&screen_height=982&browser_language=zh-CN&browser_platform=MacIntel&browser_name=Chrome&browser_version=124.0.0.0&browser_online=true&engine_name=Blink&engine_version=124.0.0.0&os_name=Mac+OS&os_version=10.15.7&cpu_core_num=8&device_memory=8&platform=PC&downlink=1.45&effective_type=3g&round_trip_time=600&webid=7233458391619110455&msToken=yOZbWIHcm8GM2ehTP7BaP50oviFNn9uX57gBlScZHqw89fh3nld-JyuRPfEp1jDckLru5MYvi6NpUh6ZZyd41LFiRHhCPO7YqrhCo2jLSiEAWlbN2HklYv3I6JnVfXZldA%3D%3D&a_bogus=mvmqQQuXmEDsfd6k56ALfY3q64M3YZ%2FQ0CPYMD2ffxVgbL39HMTE9exYXQsvQzyjLG%2FlIeSjy4hJT3eMxQVrA3vX9WEKlIOp-g00tFcQ-Izj-qjeeL80n4JO5kY3SFFB57NIxORkw7QCSYmpAdAj-kIAP62kFobyifELtIY%3D"

func getComments(awemeId string, cursor int, count int) (goJson *gojson.Json, err error) {

	// 设置请求参数
	reqUrl, err := url.Parse(commentURL)
	if err != nil {
		return
	}

	params := url.Values{}
	paramSlices := strings.Split(paramStr, "&")
	for _, param := range paramSlices {
		rawVal := strings.Split(param, "=")
		params.Set(rawVal[0], rawVal[1])
	}
	params.Set("cursor", strconv.Itoa(cursor))
	params.Set("count", strconv.Itoa(count))
	params.Set("aweme_id", awemeId)
	params.Set("msToken", utils.RandString(128))
	reqUrl.RawQuery = params.Encode()
	// 定义req对象
	client := &http.Client{}
	req, err := http.NewRequest("GET", reqUrl.String(), nil)

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
	req.Header.Set("User-Agent", userAgent)

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

const replyUrl = "https://www.douyin.com/aweme/v1/web/comment/list/reply/"

func getReply(item_id, comment_id string, cursor int, count int) (goJson *gojson.Json, err error) {

	// 设置请求参数

	var replyParams = "device_platform=webapp&aid=6383&channel=channel_pc_web&item_id=" + item_id + "&comment_id=" + comment_id + "&cut_version=1&cursor=" + strconv.Itoa(cursor) + "&count=" + strconv.Itoa(count) + "&item_type=0&update_version_code=170400&pc_client_type=1&version_code=170400&version_name=17.4.0&cookie_enabled=true&screen_width=1512&screen_height=982&browser_language=zh-CN&browser_platform=MacIntel&browser_name=Chrome&browser_version=124.0.0.0&browser_online=true&engine_name=Blink&engine_version=124.0.0.0&os_name=Mac+OS&os_version=10.15.7&cpu_core_num=8&device_memory=8&platform=PC&downlink=1.4&effective_type=3g&round_trip_time=600&webid=7233458391619110455&msToken=" + utils.RandString(128) + "&verifyFp=verify_lw53dly9_Sfb80Jjv_JIMT_4YKv_BVVq_vYe0bPggLOVK&fp=verify_lw53dly9_Sfb80Jjv_JIMT_4YKv_BVVq_vYe0bPggLOVK&a_bogus=df8wQRuvmEVkgfSv5XALfY3q6XF3YM0-0trEMD2fedfVxy39HMYY9exLvZ0vfc6jLG%2FlIebjy4heT3NMxQVrA3vX9WEKlIOp-g00tFcQ-ITj-qyeeL80n4JO5kY3SFFB57NIxORkqwAGKuRsAINe-7qvPE9jLojAYim7epr3"
	// 定义req对象
	client := &http.Client{}
	req, err := http.NewRequest("GET", replyUrl+"?"+replyParams, nil)

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
	req.Header.Set("User-Agent", userAgent)

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

var fileHeader = []string{"抖音号", "昵称", "ip归属地", "评论时间", "点赞数", "评论类型", "评论内容", "用户主页"}

type commentType int

const (
	topCommentType commentType = iota + 1
	bottomCommentType
)

// 获取一个视频的所有的评论 （一级评论 + 二级评论）
func GetAwemeComment(aweme_id string, douyinId string) {

	file := csv.NewFile(aweme_id+"-"+strconv.FormatInt(time.Now().Unix(), 10)+".csv", douyinId)
	file.Write(fileHeader)
	defer file.Close()

	var page = 0
	var pageSize = 20
	for {
		js, err := getComments(aweme_id, page, pageSize)
		if err != nil {
			logger.Error("getComments err:", err)
			break
		}

		// 解析评论信息
		comments := js.Get("comments")

		var data [][]string
		for i := 0; i < comments.ArrayLen(); i++ {
			item := comments.GetIndex(i)
			data = append(data, getOneComment(item, topCommentType))     // 主评论
			reply_comment_total := item.Get("reply_comment_total").Int() // 子评论数

			if reply_comment_total > 0 { // 子评论

				cid := item.Get("cid").String()
				replyCursor := 0
				replyCount := 3
				for {
					replyJson, err := getReply(aweme_id, cid, replyCursor, replyCount)
					if err != nil {
						logger.Error("getReply err:", err)
						break
					}

					// 记录评论信息

					replyComments := replyJson.Get("comments")

					for i := 0; i < replyComments.ArrayLen(); i++ {
						item := replyComments.GetIndex(i)
						data = append(data, getOneComment(item, bottomCommentType)) // 子评论
					}

					replyCursor = replyJson.Get("cursor").Int()

					if replyJson.Get("has_more").Int() != 1 {
						logger.Debug("视频:", aweme_id, "comment_id", cid, "子评论has_more:", js.Get("has_more").Int())
						break
					}
				}
			}
		}

		file.WriteAll(data)

		// 下一页的起始位置
		page = js.Get("cursor").Int()
		// 没有更多了
		if js.Get("has_more").Int() != 1 {
			logger.Info("视频:", aweme_id, " 全部爬取完成")
			break
		}
	}

}

func getOneComment(item *gojson.Json, cType commentType) []string {
	user := item.Get("user")
	oneComment := make([]string, 0, len(fileHeader))
	unique_id := user.Get("unique_id").String()
	if unique_id == "" {
		unique_id = user.Get("short_id").String()
	}
	oneComment = append(oneComment, unique_id)                                                                   // 抖音号
	oneComment = append(oneComment, user.Get("nickname").String())                                               // 昵称
	oneComment = append(oneComment, item.Get("ip_label").String())                                               // 地区
	oneComment = append(oneComment, time.Unix(item.Get("create_time").Int64(), 0).Format("2006-01-02 15:04:05")) // 评论时间
	oneComment = append(oneComment, item.Get("digg_count").String())                                             // 评论点赞数
	if cType == topCommentType {
		oneComment = append(oneComment, "一级评论")
	} else {
		oneComment = append(oneComment, "二级评论")
	}
	oneComment = append(oneComment, item.Get("text").String())                                   // 评论内容
	oneComment = append(oneComment, "https://www.douyin.com/user/"+user.Get("sec_uid").String()) // 用户主页

	return oneComment
}
