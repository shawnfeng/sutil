package sutil

import (
	"hash/fnv"
	"encoding/json"
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

