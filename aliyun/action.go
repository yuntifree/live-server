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
	url := host + "?" + query + "Signature=" + sign
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
	url := host + "?" + query + "Signature=" + sign
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
	url := host + "?" + query + "Signature=" + sign
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
	url := host + "?" + query + "Signature=" + sign
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
	url := host + "?" + query + "Signature=" + sign
	rsp, err := httputil.Request(url, "")
	if err != nil {
		log.Printf("ResumeLiveStream failed:%s %v", url, err)
		return ""
	}
	return rsp
}
