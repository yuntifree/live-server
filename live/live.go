package live

import (
	"fmt"

	"github.com/yuntifree/components/strutil"
)

//GenAuthKey generate auth key
func GenAuthKey(uri, key string, timestamp, rand, uid int64) string {
	str := fmt.Sprintf("%s-%d-%d-%d-%s", uri, timestamp, rand,
		uid, key)
	return strutil.MD5(str)
}
