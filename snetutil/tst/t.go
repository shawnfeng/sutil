package main

import (
	"net"
	"fmt"
	"bytes"
)


func IpBetween(from, to, test net.IP) (bool, error) {
    if from == nil || to == nil || test == nil {
        return false, fmt.Errorf("An ip input is nil")
    }

    from16 := from.To16()
    to16 := to.To16()
    test16 := test.To16()
    if from16 == nil || to16 == nil || test16 == nil {
        return false, fmt.Errorf("An ip did not convert to a 16 byte")
    }

    if bytes.Compare(test16, from16) >= 0 && bytes.Compare(test16, to16) <= 0 {
        return true, nil
    }
    return false, nil
}


func IpBetweenStr(from, to, test string) (bool, error) {
	return IpBetween(net.ParseIP(from), net.ParseIP(to), net.ParseIP(test))
}

func ck(host string) {

	ip := net.ParseIP(host)

	if ip == nil {
		fmt.Errorf("ParseIP error:%s", host)
	}

	fmt.Println("ADDR TYPE", ip,
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



}


func GetLocalIP() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return ""
    }
    for _, address := range addrs {
        // check the address type and if it is not a loopback the display it
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
				fmt.Println(address, ipnet.IP.String())
                //return ipnet.IP.String()
            }
        }
    }
    return ""
}


func HandleIpBetween(from string, to string, test string, assert bool) {

    res, err := IpBetweenStr(from, to, test)
	if err != nil {
        fmt.Printf("check err:%s from:%s to:%s test:%s\n", err, from, to, test)
		return
	}
    if res != assert {
        fmt.Println("Assertion (have: %s should be: %s) failed on range %s-%s with test %s", res, assert, from, to, test)
		return
    }

	fmt.Println("ck OK", from, to, test, assert)
}

func main() {
	HandleIpBetween("0.0.0.0", "255.255.255.255", "0.0.0.0", true)
	HandleIpBetween("0.0.0.0", "255.255.255.255", "255.255.255.255", true)
	HandleIpBetween("0.0.0.0", "255.255.255.255", "128.128.128.128", true)
    HandleIpBetween("0.0.0.0", "128.128.128.128", "255.255.255.255", false)
    HandleIpBetween("74.50.153.0", "74.50.153.4", "74.50.153.0", true)
    HandleIpBetween("74.50.153.0", "74.50.153.4", "74.50.153.4", true)
    HandleIpBetween("74.50.153.0", "74.50.153.4", "74.50.153.5", false)
    HandleIpBetween("2001:0db8:85a3:0000:0000:8a2e:0370:7334", "74.50.153.4", "74.50.153.2", false)
    HandleIpBetween("2001:0db8:85a3:0000:0000:8a2e:0370:7334", "2001:0db8:85a3:0000:0000:8a2e:0370:8334", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", true)
    HandleIpBetween("2001:0db8:85a3:0000:0000:8a2e:0370:7334", "2001:0db8:85a3:0000:0000:8a2e:0370:8334", "2001:0db8:85a3:0000:0000:8a2e:0370:7350", true)
    HandleIpBetween("2001:0db8:85a3:0000:0000:8a2e:0370:7334", "2001:0db8:85a3:0000:0000:8a2e:0370:8334", "2001:0db8:85a3:0000:0000:8a2e:0370:8334", true)
    HandleIpBetween("2001:0db8:85a3:0000:0000:8a2e:0370:7334", "2001:0db8:85a3:0000:0000:8a2e:0370:8334", "2001:0db8:85a3:0000:0000:8a2e:0370:8335", false)
    HandleIpBetween("::ffff:192.0.2.128", "::ffff:192.0.2.250", "::ffff:192.0.2.127", false)
    HandleIpBetween("::ffff:192.0.2.128", "::ffff:192.0.2.250", "::ffff:192.0.2.128", true)
    HandleIpBetween("::ffff:192.0.2.128", "::ffff:192.0.2.250", "::ffff:192.0.2.129", true)
    HandleIpBetween("::ffff:192.0.2.128", "::ffff:192.0.2.250", "::ffff:192.0.2.250", true)
    HandleIpBetween("::ffff:192.0.2.128", "::ffff:192.0.2.250", "::ffff:192.0.2.251", false)
    HandleIpBetween("::ffff:192.0.2.128", "::ffff:192.0.2.250", "192.0.2.130", true)
    HandleIpBetween("192.0.2.128", "192.0.2.250", "::ffff:192.0.2.130", true)
    HandleIpBetween("idonotparse", "192.0.2.250", "::ffff:192.0.2.130", false)

	//10.0.0.0/8：10.0.0.0～10.255.255.255
	//172.16.0.0/12：172.16.0.0～172.31.255.255
	//192.168.0.0/16：192.168.0.0～192.168.255.255
    HandleIpBetween("10.0.0.0", "10.255.255.255", "10.0.0.0", true)
    HandleIpBetween("10.0.0.0", "10.255.255.255", "10.255.255.255", true)
    HandleIpBetween("10.0.0.0", "10.255.255.255", "10.1.2.3", true)
    HandleIpBetween("10.0.0.0", "10.255.255.255", "11.1.2.3", false)

    HandleIpBetween("172.16.0.0", "172.31.255.255", "10.0.0.0", false)

    HandleIpBetween("172.16.0.0", "172.31.255.255", "172.56.15.175", false)


	//192.168.0.0/16：192.168.0.0～192.168.255.255

    HandleIpBetween("192.168.0.0", "192.168.255.255", "192.169.0.0", false)


	return

	fmt.Println(GetLocalIP())
	return

	ck("127.0.0.1")
	ck("127.0.0.2")
	ck("127.0.0.3")

	ck("172.56.15.175")


	ck("10.169.17.40")

	ck("192.168.1.198")


}
