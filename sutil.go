package sutil

import (
	"hash/fnv"
)


func HashV(addrs []string, key string) string {
    h := fnv.New32a()
    h.Write([]byte(key))
    hv := h.Sum32()

	return addrs[hv % uint32(len(addrs))]

}




