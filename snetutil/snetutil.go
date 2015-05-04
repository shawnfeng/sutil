package snetutil

import (
	"net"
	"time"
	"fmt"
	"bytes"
	"io/ioutil"
	"encoding/binary"
	"strings"
	"errors"
	"net/http"
)




func GetInterIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}


	for _, addr := range addrs {
		//fmt.Printf("Inter %v\n", addr)
		ip := addr.String()
		if "10." == ip[:3] {
			return strings.Split(ip, "/")[0], nil
		} else if "172." == ip[:4] {
			return strings.Split(ip, "/")[0], nil
		} else if "196." == ip[:4] {
			return strings.Split(ip, "/")[0], nil
		} else if "192." == ip[:4] {
			return strings.Split(ip, "/")[0], nil
		}



	}

	return "", errors.New("no inter ip")
}

// 获取首个外网ip v4
func GetExterIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}


	for _, addr := range addrs {
		//fmt.Printf("Inter %v\n", addr)
		ips := addr.String()
		idx := strings.LastIndex(ips, "/")
		if idx == -1 {
                        continue
		}
		ipv := net.ParseIP(ips[:idx])
		if ipv == nil {
                        continue
		}

		ipv4 := ipv.To4()
		if ipv4 == nil {
                        // ipv6
                        continue
		}
		ip := ipv4.String()

		if "10." != ip[:3] && "172." != ip[:4] && "196." != ip[:4] && "127." != ip[:4] {
			return ip, nil
		}

	}

	return "", errors.New("no exter ip")
}

func GetServAddr(a net.Addr) (string, error) {
	addr := a.String()
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}

	ip := net.ParseIP(host)

	if ip == nil {
		return "", fmt.Errorf("ParseIP error:%s", host)
	}
	/*
	slog.Infoln("ADDR TYPE", ip,
		"IsGlobalUnicast",
		ip.IsGlobalUnicast(),
		"IsInterfaceLocalMulticast",
		ip.IsInterfaceLocalMulticast(),
		"IsLinkLocalMulticast",
		ip.IsLinkLocalMulticast(),
		"IsLinkLocalUnicast",
		ip.IsLinkLocalUnicast(),
		"IsLoopback",
		ip.IsLoopback(),
		"IsMulticast",
		ip.IsMulticast(),
		"IsUnspecified",
		ip.IsUnspecified(),
	)
    */

	raddr := addr
	if ip.IsUnspecified() {
		// 没有指定ip的情况下，使用内网地址
		inerip, err := GetInterIp()
		if err != nil {
			return "", err
		}

		raddr = net.JoinHostPort(inerip, port)
	}

	//slog.Tracef("ServAddr --> addr:[%s] ip:[%s] host:[%s] port:[%s] raddr[%s]", addr, ip, host, port, raddr)

	return raddr, nil
}



// Request.RemoteAddress contains port, which we want to remove i.e.:
// "[::1]:58292" => "[::1]"
func IpAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

func IpAddrPort(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return ""
	}
	return s[idx+1:]
}


// 获取http请求的client的地址
func IpAddressHttpClient(r *http.Request) string {
	hdr := r.Header
	hdrRealIp := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")

	if hdrRealIp == "" && hdrForwardedFor == "" {
		return IpAddrFromRemoteAddr(r.RemoteAddr)
	}

	if hdrForwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		// TODO: should return first non-local address
		for _, ip := range(parts) {
			if len(ip) > 5 && "10." != ip[:3] && "172." != ip[:4] && "196." != ip[:4] && "127." != ip[:4] {
				return ip
			}
		}

	}

	return hdrRealIp
}



func PackdataPad(data []byte, pad byte) []byte {
	sendbuff := make([]byte, 0)
	// no pad
	var pacLen uint64 = uint64(len(data))
	buff := make([]byte, 20)
	rv := binary.PutUvarint(buff, pacLen)

	sendbuff = append(sendbuff, buff[:rv]...) // len
	sendbuff = append(sendbuff, data...) //data
	sendbuff = append(sendbuff, pad) //pad

	return sendbuff

}

func Packdata(data []byte) []byte {
	return PackdataPad(data, 0)
}


// 最小的消息长度、最大消息长度，数据流，包回调
// 正常返回解析剩余的数据，nil
// 否则返回错误
func UnPackdata(lenmin uint, lenmax uint, packBuff []byte, readCall func([]byte)) ([]byte, error) {
	for {

		// n == 0: buf too small
		// n  < 0: value larger than 64 bits (overflow)
        //     and -n is the number of bytes read
		pacLen, sz := binary.Uvarint(packBuff)
		if sz < 0 {
			return packBuff, errors.New("package head error")
		} else if sz == 0 {
			return packBuff, nil
		}

		// sz > 0

		// must < lenmax
		if pacLen > uint64(lenmax) {
			return packBuff, errors.New("package too long")
		} else if pacLen < uint64(lenmin) {
			return packBuff, errors.New("package too short")
		}

		apacLen := uint64(sz)+pacLen+1
		if uint64(len(packBuff)) >= apacLen {
			pad := packBuff[apacLen-1]
			if pad != 0 {
				return packBuff, errors.New("package pad error")
			}

			readCall(packBuff[sz:apacLen-1])
			packBuff = packBuff[apacLen:]
		} else {
			return packBuff, nil
		}


	}


	return nil, errors.New("unknown err")


}

func HttpReqGetOk(url string, timeout time.Duration) ([]byte, error) {

	client := &http.Client{Timeout: timeout}
	response, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("statuscode:%d body:%s", response.StatusCode, body)

	} else {
		return body, nil
	}


}

func HttpReqPostOk(url string, data []byte, timeout time.Duration) ([]byte, error) {
	body, status, err := HttpReqPost(url, data, timeout)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, errors.New(fmt.Sprintf("status:%d err:%s", status, body))

	} else {
		return body, nil
	}

}

func HttpReqPost(url string, data []byte, timeout time.Duration) ([]byte, int, error) {
	client := &http.Client{Timeout: timeout}

	reqest, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, 0, err
	}
	reqest.Header.Set("Connection","Keep-Alive")

	response, err := client.Do(reqest)
	if err != nil {
		return nil, 0, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, 0, err
	}

	return body, response.StatusCode, nil

}




func PackageSplit(conn net.Conn, readtimeout time.Duration, readCall func([]byte)) (bool, []byte, error) {
	buffer := make([]byte, 2048)
	packBuff := make([]byte, 0)

	for {
		conn.SetReadDeadline(time.Now().Add(readtimeout))
		bytesRead, err := conn.Read(buffer)
		if err != nil {
			return true, nil, err
		}



		packBuff = append(packBuff, buffer[:bytesRead]...)

		packBuff, err = UnPackdata(1, 1024*5, packBuff, readCall)

		if err != nil {
			return false, packBuff, err
		}


	}

	return false, nil, errors.New("fuck err")

}
