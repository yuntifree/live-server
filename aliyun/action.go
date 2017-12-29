package aliyun

import (
	"fmt"
	"log"
	"net/url"
	"sort"
	"time"

	"github.com/yuntifree/components/httputil"
	"github.com/yuntifree/components/strutil"
)

//GenAuthKey generate auth key
func GenAuthKey(uri, key string, timestamp, rand, uid int64) string {
	str := fmt.Sprintf("%s-%d-%d-%d-%s", uri, timestamp, rand,
		uid, key)
	return strutil.MD5(str)
}

//GenPushURL generate push url
func GenPushURL(name string, uid int64) string {
	stream := genAuthStream(name, uid)
	return fmt.Sprintf("%s/%s/%s", pushHost, appname, stream)
}

func genAuthStream(name string, uid int64) string {
	uri := "/" + appname + "/" + name
	ts := time.Now().Unix() + 3600
	auth := GenAuthKey(uri, authKey, ts, 0, uid)
	return fmt.Sprintf("%s?vhost=%s&auth_key=%d-0-%d-%s",
		name, vhost, ts, uid, auth)
}

func getResolutionID(resolution int64) string {
	var id string
	switch resolution {
	case 0:
		id = "lld"
	case 1:
		id = "lsd"
	case 2:
		id = "lhd"
	case 3:
		id = "lud"
	}
	return id
}

//GenLiveHLS generate live hls url
func GenLiveHLS(name string) string {
	ts := time.Now().Unix()
	str := fmt.Sprintf("/%s/%s.m3u8-%d-0-0-%s", appname, name, ts, authKey)
	sign := strutil.MD5(str)
	return fmt.Sprintf("http://%s/%s/%s.m3u8?auth_key=%d-0-0-%s", vhost, appname, name,
		ts, sign)
}

//GenLiveRTMP generate live rtmp
func GenLiveRTMP(name string, resolution int64) string {
	ts := time.Now().Unix()
	id := getResolutionID(resolution)
	str := fmt.Sprintf("/%s/%s_%s-%d-0-0-%s", appname, name, id, ts, authKey)
	sign := strutil.MD5(str)
	return fmt.Sprintf("rtmp://%s/%s_%s?auth_key=%d-0-0-%s", vhost, name, id, ts, sign)
}

//GenLiveFLV generate live flv
func GenLiveFLV(name string, resolution int64) string {
	ts := time.Now().Unix()
	id := getResolutionID(resolution)
	str := fmt.Sprintf("/%s/%s_%s.flv-%d-0-0-%s", appname, name, id, ts, authKey)
	sign := strutil.MD5(str)
	return fmt.Sprintf("http://%s/%s_%s.flv?auth_key=%d-0-0-%s", vhost, name, id, ts, sign)
}

func genSign(m map[string]string, key string) string {
	var sortedKeys []string
	for k := range m {
		sortedKeys = append(sortedKeys, k)
	}

	sort.Strings(sortedKeys)
	prefix := "GET&%2F&"
	var signStr string
	for _, k := range sortedKeys {
		log.Printf("%v -- %v", k, m[k])
		value := fmt.Sprintf("%v", m[k])
		if value != "" {
			signStr += k + "=" + url.QueryEscape(value) + "&"
		}
	}
	l := len(signStr)
	signStr = signStr[:l-1]
	signStr = url.QueryEscape(signStr)
	signStr = prefix + signStr
	fmt.Printf("signStr:%s\n", signStr)

	return strutil.HmacSha1(signStr, key+"&")
}

func getTimeStr() string {
	ts := time.Now().UTC()
	return fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02dZ",
		ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), ts.Second())
}

func fillCommParams() map[string]string {
	m := make(map[string]string)
	m["Format"] = format
	m["Version"] = version
	m["AccessKeyId"] = accessKeyID
	m["SignatureMethod"] = signMethod
	m["Timestamp"] = getTimeStr()
	m["SignatureVersion"] = signVersion
	m["SignatureNonce"] = strutil.GenUUID()
	return m
}

func toString(m map[string]string) string {
	var query string
	for k, v := range m {
		query += k + "=" + v + "&"
	}
	return query
}

//DescribeLiveStreamPublishList list stream
func DescribeLiveStreamPublishList(startTime, endTime string) string {
	m := fillCommParams()
	m["Action"] = "DescribeLiveStreamsPublishList"
	m["DomainName"] = domain
	m["StartTime"] = startTime
	m["EndTime"] = endTime
	sign := genSign(m, accessKeySecret)
	query := toString(m)
	url := apiHost + "?" + query + "Signature=" + sign
	rsp, err := httputil.Request(url, "")
	if err != nil {
		log.Printf("DescribeLiveStreamPublishList failed:%s %v", url, err)
		return ""
	}
	return rsp
}

//DescribeLiveStreamOnlineUserNum list online user num
func DescribeLiveStreamOnlineUserNum() string {
	m := fillCommParams()
	m["Action"] = "DescribeLiveStreamOnlineUserNum"
	m["DomainName"] = domain
	sign := genSign(m, accessKeySecret)
	query := toString(m)
	url := apiHost + "?" + query + "Signature=" + sign
	rsp, err := httputil.Request(url, "")
	if err != nil {
		log.Printf("DescribeLiveStreamOnlineUserNum failed:%s %v", url, err)
		return ""
	}
	return rsp
}

//ForbidLiveStream forbid live stream
func ForbidLiveStream(stream string) string {
	m := fillCommParams()
	m["Action"] = "ForbidLiveStream"
	m["DomainName"] = domain
	m["AppName"] = appname
	m["StreamName"] = stream
	m["LiveStreamType"] = "publisher"
	sign := genSign(m, accessKeySecret)
	query := toString(m)
	url := apiHost + "?" + query + "Signature=" + sign
	rsp, err := httputil.Request(url, "")
	if err != nil {
		log.Printf("ForbidLiveStream failed:%s %v", url, err)
		return ""
	}
	return rsp
}

//ResumeLiveStream resume live stream
func ResumeLiveStream(stream string) string {
	m := fillCommParams()
	m["Action"] = "ResumeLiveStream"
	m["DomainName"] = domain
	m["AppName"] = appname
	m["StreamName"] = stream
	m["LiveStreamType"] = "publisher"
	sign := genSign(m, accessKeySecret)
	query := toString(m)
	url := apiHost + "?" + query + "Signature=" + sign
	rsp, err := httputil.Request(url, "")
	if err != nil {
		log.Printf("ResumeLiveStream failed:%s %v", url, err)
		return ""
	}
	return rsp
}

//DescribeLiveStreamsFrameRateAndBitRateData stream frame rate and bitrate
func DescribeLiveStreamsFrameRateAndBitRateData() string {
	m := fillCommParams()
	m["Action"] = "DescribeLiveStreamsFrameRateAndBitRateData"
	m["DomainName"] = domain
	sign := genSign(m, accessKeySecret)
	query := toString(m)
	url := apiHost + "?" + query + "Signature=" + sign
	rsp, err := httputil.Request(url, "")
	if err != nil {
		log.Printf("ResumeLiveStream failed:%s %v", url, err)
		return ""
	}
	return rsp
}
