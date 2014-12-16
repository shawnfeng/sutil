package sutil

import (
	"hash/fnv"
	"io/ioutil"
	"fmt"
	"strings"
	"os"
	"encoding/json"
	"unicode/utf8"
	"crypto/md5"
	//"crypto/sha1"

	"sync"
	"code.google.com/p/go-uuid/uuid"
)


func HashV(addrs []string, key string) string {
    h := fnv.New32a()
    h.Write([]byte(key))
    hv := h.Sum32()

	return addrs[hv % uint32(len(addrs))]

}


func IsJSON(s []byte) bool {
    //var js map[string]interface{}
    var js interface{}
    return json.Unmarshal(s, &js) == nil
}

func GetUtf8Chars(s string, num int) string {
	b := []byte(s)
	rv := ""
	for i := 0; len(b)>0 && i < num; i++ {
		_, size := utf8.DecodeRune(b)
		rv += string(b[:size])
		b = b[size:]
	}

	return rv
}

var uuidMu sync.Mutex
func GetUUID() string {
	uuidMu.Lock()
	defer uuidMu.Unlock()

	uuidgen := uuid.NewUUID()
	return uuidgen.String()
}


func GetUniqueMd5() string {
	u := GetUUID()
	h := md5.Sum([]byte(u))
	return fmt.Sprintf("%x", h)
}


// 文件输出，目录不存在自动创建
func WriteFile(path string, data []byte, perm os.FileMode) error {

	idx := strings.LastIndex(path, "/")
	if idx != -1 {
		logdir := path[:idx]
		err := os.MkdirAll(logdir, 0777)
		if err != nil {
			return err
		}
	}

	return ioutil.WriteFile(path, data, perm)
}
