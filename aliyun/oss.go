package aliyun

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"log"
	"strings"
	"time"

	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	simplejson "github.com/bitly/go-simplejson"
)

const (
	innerEndpoint  = "oss-cn-shenzhen-internal.aliyuncs.com"
	endpoint       = "oss-cn-shenzhen.aliyuncs.com"
	yuntiBucket    = "yuntinews"
	bucketURL      = "http://yuntinews.oss-cn-shenzhen.aliyuncs.com"
	newsCdnURL     = "http://news.yunxingzh.com"
	expireInterval = 15 * 60
	imgOuterHost   = "http://yuntiimgs.oss-cn-shenzhen.aliyuncs.com"
	maxImageSize   = 4 * 1024 * 1024
	imageBucket    = "yuntiimgs"
	fileBucket     = "yuntifile"
	ossCbURL       = "http://video.yunxingzh.com/image/callback"
	ossCbBody      = "filename=${object}&size=${size}&mimeType=${mimeType}&height=${imageInfo.height}&width=${imageInfo.width}"
	ossCbBodyType  = "application/x-www-form-urlencoded"
	fileCdnURL     = "http://file.yunxingzh.com"
	imgCdnURL      = "http://img.yunxingzh.com"
)

//UploadOssFile upload content to aliyun oss
func UploadOssFile(filename, content string) bool {
	return uploadOssBucket(filename, content, yuntiBucket)
}

//UploadYuntiFile upload content to aliyun oss
func UploadYuntiFile(filename, content string) bool {
	return uploadOssBucket(filename, content, fileBucket)
}

//UploadOssImg upload img to aliyun oss
func UploadOssImg(filename, content string) bool {
	return uploadOssBucket(filename, content, imageBucket)
}

//UploadOssImgFromFile upload img to aliyun oss read from file
func UploadOssImgFromFile(filename, filepath string) bool {
	return uploadOssBucketFromFile(filename, filepath, imageBucket)
}

func uploadOssBucketFromFile(filename, filepath, ossbucket string) bool {
	client, err := oss.New(innerEndpoint, accessKeyID, accessKeySecret)
	if err != nil {
		log.Printf("oss init failed:%v", err)
		return false
	}

	bucket, err := client.Bucket(ossbucket)
	if err != nil {
		log.Printf("bucket init failed:%v", err)
		return false
	}

	err = bucket.UploadFile(filename, filepath, 100*1024)
	if err != nil {
		log.Printf("upload failed %s: %v", filename, err)
		return false
	}

	return true
}

//uploadOssBucket upload file to alioss bucket
func uploadOssBucket(filename, content, ossbucket string) bool {
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		log.Printf("oss init failed:%v", err)
		return false
	}

	bucket, err := client.Bucket(ossbucket)
	if err != nil {
		log.Printf("bucket init failed:%v", err)
		return false
	}

	err = bucket.PutObject(filename, strings.NewReader(content))
	if err != nil {
		log.Printf("PutObject failed %s: %v", filename, err)
		return false
	}

	return true
}

//GenOssNewsURL generate oss news download url
func GenOssNewsURL(filename string) string {
	return newsCdnURL + "/" + filename
}

//GenOssFileURL generate oss file download url
func GenOssFileURL(filename string) string {
	return fileCdnURL + "/" + filename
}

//GenOssImgURL generate oss img download url
func GenOssImgURL(filename string) string {
	return imgCdnURL + "/" + filename
}

func getISO8601Time(ts time.Time) string {
	return ts.Format("2006-01-02T15:04:05Z")
}

func genPolicy(expire time.Time) string {
	json, _ := simplejson.NewJson([]byte(`{}`))
	expireStr := getISO8601Time(expire)
	var c1 = [3]interface{}{"content-length-range", 0, maxImageSize}
	var c2 = [3]interface{}{"starts-with", "$key", ""}
	var conditions = [2]interface{}{c1, c2}
	json.Set("expiration", expireStr)
	json.Set("conditions", conditions)
	data, _ := json.Encode()
	return base64.StdEncoding.EncodeToString(data)
}

func genHmacSign(content, key string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(content))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func genCallback() string {
	json, _ := simplejson.NewJson([]byte(`{}`))
	json.Set("callbackUrl", ossCbURL)
	json.Set("callbackBody", ossCbBody)
	json.Set("callbackBodyType", ossCbBodyType)
	data, _ := json.Encode()
	return base64.StdEncoding.EncodeToString(data)
}

//FillPolicyResp generate upload policy response
func FillPolicyResp(json *simplejson.Json) {
	expire := time.Now().Add(expireInterval * time.Second)
	json.Set("accessid", accessKeyID)
	json.Set("host", imgOuterHost)
	policy := genPolicy(expire)
	json.Set("policy", policy)
	sign := genHmacSign(policy, accessKeySecret)
	json.Set("signature", sign)
	json.Set("dir", "")
	json.Set("expire", expire.Unix())
	callback := genCallback()
	json.Set("callback", callback)
	return
}

//FillCallbackInfo for apply_image_upload fill callback info
func FillCallbackInfo(js *simplejson.Json) {
	js.Set("bucket", imageBucket)
	js.Set("callbackurl", ossCbURL)
	js.Set("callbackbody", ossCbBody)
}
