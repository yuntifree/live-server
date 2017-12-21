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
	log.Printf("url:%s\n", url)
	rsp, err := httputil.Request(url, "")
	if err != nil {
		log.Printf("DescribeLiveStreamPublishList failed:%s %v", url, err)
		return ""
	}
	return rsp
}
