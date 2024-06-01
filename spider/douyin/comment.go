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

const cookieStr = "ttwid=1%7CFVE2kMPBqFvkO4qA0LratVDDT_BsgZcjxH9Yrdmf2q8%7C1684170786%7Ce6aea5165f82b8a6b1d623a26d2e7317e31341ef4dd16ab07710d73a5a434f9d; LOGIN_STATUS=1; store-region=cn-ah; store-region-src=uid; my_rd=1; bd_ticket_guard_client_web_domain=2; dy_swidth=1512; dy_sheight=982; live_use_vvc=%22false%22; SEARCH_RESULT_LIST_TYPE=%22single%22; xgplayer_device_id=29315359380; passport_assist_user=CjyogSxaLwNt7PeDD6g5bMxENqNErkVSzOnN6fpG5ZyOLKDP4QBsYcCH3HJGOtFEA3-fAJ0spjt_VyEvIgIaSgo8xghSo7ZW1S7BW6V1bkziJlpiNCzxAaPzXw1n5FHmSqegN_FJIVhuyP3Q_S-W82r5Eiirx3cmjSIQepzFEKrgzg0Yia_WVCABIgED1cI8XA%3D%3D; n_mh=5sxwI9YZgUwJ1h3k7SsM8ZO0MNR7APLsq56Z8YZk9nQ; sso_uid_tt=9e7ecb8a01d3aeac26e4be4f34abd603; sso_uid_tt_ss=9e7ecb8a01d3aeac26e4be4f34abd603; toutiao_sso_user=05be5032ed2dd33dd3ea2d34576beae2; toutiao_sso_user_ss=05be5032ed2dd33dd3ea2d34576beae2; uid_tt=d1ad0bf67a15f064259b37a07105851f; uid_tt_ss=d1ad0bf67a15f064259b37a07105851f; sid_tt=11627b6b93e061f75da62c2db102c463; sessionid=11627b6b93e061f75da62c2db102c463; sessionid_ss=11627b6b93e061f75da62c2db102c463; _bd_ticket_crypt_doamin=2; _bd_ticket_crypt_cookie=42e301000c4ea541972f2cf8a23e192c; __security_server_data_status=1; douyin.com; device_web_cpu_core=8; device_web_memory_size=8; csrf_session_id=ff5a5854443708de6da8386180701fbc; passport_fe_beating_status=true; sid_ucp_sso_v1=1.0.0-KGFlM2I2NjcxMzc2ODQzNGFmNTZhNjI4MzkwNDNhZDUyMmJlY2RjMWEKHQjYidvXygIQjv_zsQYY7zEgDDC827_TBTgGQPQHGgJsZiIgMDViZTUwMzJlZDJkZDMzZGQzZWEyZDM0NTc2YmVhZTI; ssid_ucp_sso_v1=1.0.0-KGFlM2I2NjcxMzc2ODQzNGFmNTZhNjI4MzkwNDNhZDUyMmJlY2RjMWEKHQjYidvXygIQjv_zsQYY7zEgDDC827_TBTgGQPQHGgJsZiIgMDViZTUwMzJlZDJkZDMzZGQzZWEyZDM0NTc2YmVhZTI; sid_guard=11627b6b93e061f75da62c2db102c463%7C1715273614%7C5184000%7CMon%2C+08-Jul-2024+16%3A53%3A34+GMT; sid_ucp_v1=1.0.0-KDZjYTg2ZmJmMmU4ZTY2YzEwN2Q3NGQ2ZTZlM2E1YjgzZmExY2NmNzYKGQjYidvXygIQjv_zsQYY7zEgDDgGQPQHSAQaAmhsIiAxMTYyN2I2YjkzZTA2MWY3NWRhNjJjMmRiMTAyYzQ2Mw; ssid_ucp_v1=1.0.0-KDZjYTg2ZmJmMmU4ZTY2YzEwN2Q3NGQ2ZTZlM2E1YjgzZmExY2NmNzYKGQjYidvXygIQjv_zsQYY7zEgDDgGQPQHSAQaAmhsIiAxMTYyN2I2YjkzZTA2MWY3NWRhNjJjMmRiMTAyYzQ2Mw; s_v_web_id=verify_lw53dly9_Sfb80Jjv_JIMT_4YKv_BVVq_vYe0bPggLOVK; passport_csrf_token=ae7d474be5a5ee98db090f0bccc213ed; passport_csrf_token_default=ae7d474be5a5ee98db090f0bccc213ed; publish_badge_show_info=%220%2C0%2C0%2C1716393925183%22; download_guide=%223%2F20240518%2F1%22; webcast_paid_live_duration=%7B%227371372792128241727%22%3A1%7D; pwa2=%220%7C0%7C3%7C0%22; __live_version__=%221.1.2.555%22; webcast_leading_last_show_time=1716917660303; webcast_leading_total_show_times=8; xg_device_score=7.379649772277164; strategyABtestKey=%221716971708.932%22; xgplayer_user_id=965101663182; live_can_add_dy_2_desktop=%221%22; stream_recommend_feed_params=%22%7B%5C%22cookie_enabled%5C%22%3Atrue%2C%5C%22screen_width%5C%22%3A1512%2C%5C%22screen_height%5C%22%3A982%2C%5C%22browser_online%5C%22%3Atrue%2C%5C%22cpu_core_num%5C%22%3A8%2C%5C%22device_memory%5C%22%3A8%2C%5C%22downlink%5C%22%3A0.15%2C%5C%22effective_type%5C%22%3A%5C%22slow-2g%5C%22%2C%5C%22round_trip_time%5C%22%3A3000%7D%22; volume_info=%7B%22isUserMute%22%3Afalse%2C%22isMute%22%3Atrue%2C%22volume%22%3A1%7D; WallpaperGuide=%7B%22showTime%22%3A1716974114778%2C%22closeTime%22%3A0%2C%22showCount%22%3A5%2C%22cursor1%22%3A99%2C%22cursor2%22%3A0%2C%22hoverTime%22%3A1715954055532%7D; stream_player_status_params=%22%7B%5C%22is_auto_play%5C%22%3A0%2C%5C%22is_full_screen%5C%22%3A0%2C%5C%22is_full_webscreen%5C%22%3A0%2C%5C%22is_mute%5C%22%3A0%2C%5C%22is_speed%5C%22%3A1%2C%5C%22is_visible%5C%22%3A0%7D%22; FOLLOW_NUMBER_YELLOW_POINT_INFO=%22MS4wLjABAAAAnHay_u3dMGKaG09nbzbp3Wp1TxUUp8sE-5IvtSG3rzw%2F1716998400000%2F0%2F1716977060800%2F0%22; __ac_nonce=0665707de00be3f6f1302; __ac_signature=_02B4Z6wo00f01akBNRgAAIDC-.4VAKyFN9WpITGAAAwYQpL5uFoNO6j7GuU1n9lYJMVjwfPAVqiDEwWBTMeyc8KP14ndD8hqIjDa.T.Z8MLxtQ1LrH51-GhRRL77aduERMhHMHZlGIbq.7Yo9b; IsDouyinActive=true; home_can_add_dy_2_desktop=%221%22; FOLLOW_LIVE_POINT_INFO=%22MS4wLjABAAAAnHay_u3dMGKaG09nbzbp3Wp1TxUUp8sE-5IvtSG3rzw%2F1716998400000%2F0%2F1716979681464%2F0%22; bd_ticket_guard_client_data=eyJiZC10aWNrZXQtZ3VhcmQtdmVyc2lvbiI6MiwiYmQtdGlja2V0LWd1YXJkLWl0ZXJhdGlvbi12ZXJzaW9uIjoxLCJiZC10aWNrZXQtZ3VhcmQtcmVlLXB1YmxpYy1rZXkiOiJCT0V4aEh3WStKZHBhWm96azYrN1ZWT21MaEdSSTVrMURmVzhRWTRxcGVzK1U3U1UvZXZObjhDQTlRZUIrWDNtSEZzbWs4dDhZekMyeEhhYjNJMjR3T3c9IiwiYmQtdGlja2V0LWd1YXJkLXdlYi12ZXJzaW9uIjoxfQ%3D%3D; odin_tt=8cf084543cbd0e5de180dba6652e23dd3eebb7f507758d5d34351bdd2cf299c055cb45f727ec87b0aeb67e74a11e4c8b; msToken=Gle8DPe0y_reUZwUK-5aPEYdgFL91l0EvRZ99nVf7iZWJ-ihNHhfXPJ8b1JlcruCFZdOa0R-tYysR1kyCn0nfE-kZnPXi_FuiNtBU0ebEq-r-nu8PuqdC9AfrF8OoVB6BQ=="

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

	var replyParams = "device_platform=webapp&aid=6383&channel=channel_pc_web&item_id=" + item_id + "&comment_id=" + comment_id + "&cut_version=1&cursor=" + strconv.Itoa(cursor) + "&count=" + strconv.Itoa(count) + "&item_type=0&update_version_code=170400&pc_client_type=1&version_code=170400&version_name=17.4.0&cookie_enabled=true&screen_width=1512&screen_height=982&browser_language=zh-CN&browser_platform=MacIntel&browser_name=Chrome&browser_version=124.0.0.0&browser_online=true&engine_name=Blink&engine_version=124.0.0.0&os_name=Mac+OS&os_version=10.15.7&cpu_core_num=8&device_memory=8&platform=PC&downlink=1.4&effective_type=3g&round_trip_time=600&webid=7233458391619110455&msToken=" + utils.RandString(128) + "&verifyFp=verify_lw53dly9_Sfb80Jjv_JIMT_4YKv_BVVq_vYe0bPggLOVK&fp=verify_lw53dly9_Sfb80Jjv_JIMT_4YKv_BVVq_vYe0bPggLOVK"
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
func GetAwemeComment(aweme_id string) {

	file := csv.NewFile(aweme_id + "-" + strconv.FormatInt(time.Now().Unix(), 10) + ".csv")
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
