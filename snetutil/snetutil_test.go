package snetutil

import (
	"testing"
	"fmt"
	"net"
)

func cmpbyte(a []byte, b []byte) bool {

	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}


	return true

}


func TestGetInterIp(t *testing.T) {
	//t.Errorf()

	ip, err := GetInterIp()
	fmt.Printf("GetInterIp:%s %v\n", ip, err)

}


func TestGetExterIp(t *testing.T) {
	//t.Errorf()
	// SplitHostPort splits a network address of the form "host:port", "[host]:port" or "[ipv6-host%zone]:port" into host or ipv6-host%zone and port. A literal address or host name for IPv6 must be enclosed in square brackets, as in "[::1]:80", "[ipv6-host]:http" or "[ipv6-host%zone]:80".
	host, port, err := net.SplitHostPort("127.0.0.1:333")
	fmt.Printf("SplitHostPort: %s-%s-%v\n", host, port, err)
	host, port, err = net.SplitHostPort("[::1]:80")
	fmt.Printf("SplitHostPort: %s-%s-%v\n", host, port, err)

	addrs, err := net.InterfaceAddrs()
	for _, addr := range addrs {
		fmt.Println(addr.Network(), addr.String())
	}


	ip, err := GetExterIp()
	fmt.Printf("GetExterIp:%s %v\n", ip, err)

}


func TestPackdataPad(t *testing.T) {
	// PackdataPad(data []byte, pad byte) []byte

	data := []byte("ABCD")
	fmt.Println(data)

	pdata := PackdataPad(data, 0)

	fmt.Println(pdata)

	if pdata[0] != 4 {
		t.Errorf("error len")
	}

	if pdata[1] != 65 {
		t.Errorf("error data")
	}

	if pdata[2] != 66 {
		t.Errorf("error data")
	}

	if pdata[3] != 67 {
		t.Errorf("error data")
	}

	if pdata[4] != 68 {
		t.Errorf("error data")
	}


	if pdata[5] != 0 {
		t.Errorf("error pad")
	}




	data = []byte{1, 2, 3, 200, 255}
	fmt.Println(data)

	pdata = PackdataPad(data, 10)

	fmt.Println(pdata)

	if pdata[0] != 5 {
		t.Errorf("error len")
	}

	if pdata[1] != 1 {
		t.Errorf("error data")
	}

	if pdata[2] != 2 {
		t.Errorf("error data")
	}

	if pdata[3] != 3 {
		t.Errorf("error data")
	}

	if pdata[4] != 200 {
		t.Errorf("error data")
	}


	if pdata[5] != 255 {
		t.Errorf("error data")
	}


	if pdata[6] != 10 {
		t.Errorf("error pad")
	}

}

func TestUnPackdataempty(t *testing.T) {

	// value larger than 64 bits (overflow)
	buff := []byte{}

	fmt.Println(buff)

	surplus, err := UnPackdata(
		1,
		1000,
		buff,
		func (pb []byte) {
			t.Errorf("error here")
		},
	)

	fmt.Printf("surplus:%v err:%v\n", surplus, err)
	if err != nil {
		t.Errorf("empty error %s", err)
	}

	if len(surplus) != 0 {
		t.Errorf("error packa")
	}



}



func TestUnPackdata0(t *testing.T) {

	// value larger than 64 bits (overflow)
	buff := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0}

	fmt.Println(buff)

	surplus, err := UnPackdata(
		1,
		1000,
		buff,
		func (pb []byte) {
			t.Errorf("error here")
		},
	)

	fmt.Printf("surplus:%v err:%v\n", surplus, err)
	if err == nil ||  "package head error" != err.Error() {
		t.Errorf("error package head anys")
	}

}


func TestUnPackdata1(t *testing.T) {

	// 超长测试
	buff := []byte{255, 1}

	fmt.Println(buff)

	surplus, err := UnPackdata(
		1,
		100,
		buff,
		func (pb []byte) {
			t.Errorf("error here")
		},
	)

	fmt.Printf("surplus:%v err:%v\n", surplus, err)
	if err == nil || "package too long" != err.Error() {
		t.Errorf("error package too long")
	}

}


func TestUnPackdata2(t *testing.T) {

	// 超长测试
	buff := []byte{4, 0}

	fmt.Println(buff)

	surplus, err := UnPackdata(
		2,
		3,
		buff,
		func (pb []byte) {
			t.Errorf("error here")
		},
	)

	fmt.Printf("surplus:%v err:%v\n", surplus, err)
	if err == nil || "package too long" != err.Error() {
		t.Errorf("error package too long")
	}

}


func TestUnPackdata3(t *testing.T) {

	buff := []byte{8, 0}

	fmt.Println(buff)

	surplus, err := UnPackdata(
		10,
		100,
		buff,
		func (pb []byte) {
			t.Errorf("error here")
		},
	)

	fmt.Printf("surplus:%v err:%v\n", surplus, err)
	if err == nil || "package too short" != err.Error() {
		t.Errorf("error package too short")
	}

}

func TestUnPackdata3_1(t *testing.T) {

	buff := []byte{5, 1, 2, 3, 4, 5, 1}

	fmt.Println(buff)

	surplus, err := UnPackdata(
		0,
		100,
		buff,
		func (pb []byte) {
			if string(pb) != string([]byte{1, 2, 3, 4, 5}) {
				t.Errorf("error packa")
			}
			fmt.Println("callback", pb)

		},
	)


	fmt.Printf("surplus:%v err:%v\n", surplus, err)
	if err == nil || "package pad error" != err.Error() {
		t.Errorf("error package pad error")
	}


}




func TestUnPackdata4(t *testing.T) {

	buff := []byte{5, 1, 2, 3, 4, 5, 0}

	fmt.Println(buff)

	surplus, err := UnPackdata(
		0,
		100,
		buff,
		func (pb []byte) {
			if string(pb) != string([]byte{1, 2, 3, 4, 5}) {
				t.Errorf("error packa")
			}
			fmt.Println("callback", pb)

		},
	)


	if err != nil{
		t.Errorf("error packa err:%s", err)
	}

	fmt.Println(surplus)
	if len(surplus) != 0 {
		t.Errorf("error packa")
	}


}


func TestUnPackdata5(t *testing.T) {

	buff := []byte{5, 1, 2, 3, 4, 5, 0, 1}

	fmt.Println(buff)

	surplus, err := UnPackdata(
		0,
		100,
		buff,
		func (pb []byte) {
			if string(pb) != string([]byte{1, 2, 3, 4, 5}) {
				t.Errorf("error packa")
			}
			fmt.Println("callback", pb)
		},
	)


	if err != nil{
		t.Errorf("error packa err:%s", err)
	}

	fmt.Println(surplus)
	if len(surplus) != 1 || surplus[0] != 1  {
		t.Errorf("error packa")
	}


}




func TestUnPackdata6(t *testing.T) {

	buff := []byte{5, 1, 2, 3, 4, 5, 0, 1}

	fmt.Println(buff)

	surplus, err := UnPackdata(
		0,
		100,
		buff,
		func (pb []byte) {
			if string(pb) != string([]byte{1, 2, 3, 4, 5}) {
				t.Errorf("error packa")
			}
			fmt.Println("callback", pb)
		},
	)


	if err != nil{
		t.Errorf("error packa err:%s", err)
	}

	fmt.Println(surplus)
	if len(surplus) != 1 || surplus[0] != 1  {
		t.Errorf("error packa")
	}


}

// 只有半个数据包
func TestUnPackdata6_1(t *testing.T) {

	buff := []byte{5, 1, 2, }

	fmt.Println(buff)

	surplus, err := UnPackdata(
		0,
		100,
		buff,
		func (pb []byte) {
			t.Errorf("error packa")
		},
	)


	if err != nil{
		t.Errorf("error packa err:%s", err)
	}

	fmt.Println(surplus)
	if !cmpbyte(surplus, buff)  {
		t.Errorf("error packa")
	}


}



// 半个数据包，只有len，且len为0
func TestUnPackdata7(t *testing.T) {

	buff := []byte{5, 1, 2, 3, 4, 5, 0, 0}

	fmt.Println(buff)

	surplus, err := UnPackdata(
		0,
		100,
		buff,
		func (pb []byte) {
			if string(pb) != string([]byte{1, 2, 3, 4, 5}) {
				t.Errorf("error packa")
			}
			fmt.Println("callback", pb)
		},
	)


	if err != nil{
		t.Errorf("error packa err:%s", err)
	}

	fmt.Println(surplus)
	if len(surplus) != 1 || surplus[0] != 0  {
		t.Errorf("error packa")
	}


}


// 半个数据包，且len够，但是数据不够
func TestUnPackdata8(t *testing.T) {

	buff := []byte{5, 1, 2, 3, 4, 5, 0, 5, 1, 2}

	fmt.Println(buff)

	surplus, err := UnPackdata(
		0,
		100,
		buff,
		func (pb []byte) {
			if string(pb) != string([]byte{1, 2, 3, 4, 5}) {
				t.Errorf("error packa")
			}
			fmt.Println("callback", pb)
		},
	)


	if err != nil{
		t.Errorf("error packa err:%s", err)
	}

	fmt.Println(surplus)
	if !cmpbyte(surplus, []byte{5, 1, 2}) {
		t.Errorf("error packa")
	}


}


// 半个数据包，且len够，数据够，pad不够
func TestUnPackdata9(t *testing.T) {
	buff := []byte{3, 1, 2, 3, 0, 2, 1, 2}

	fmt.Println(buff)

	surplus, err := UnPackdata(
		0,
		100,
		buff,
		func (pb []byte) {
			if !cmpbyte(pb, []byte{1, 2, 3}) {
				t.Errorf("error packa")
			}
			fmt.Println("callback", pb)
		},
	)


	if err != nil{
		t.Errorf("error packa err:%s", err)
	}

	fmt.Println(surplus)
	if !cmpbyte(surplus, []byte{2, 1, 2}) {
		t.Errorf("error packa")
	}


}


// 正好两个个数据包
func TestUnPackdata10(t *testing.T) {
	buff := []byte{3, 1, 2, 3, 0, 2, 1, 2, 0}

	fmt.Println(buff)

	// 解析的包计数
	pcn := 0
	surplus, err := UnPackdata(
		0,
		100,
		buff,
		func (pb []byte) {
			if pcn == 0 {
				if !cmpbyte(pb, []byte{1, 2, 3}) {
					t.Errorf("error packa")
				}
			} else if pcn == 1 {
				if !cmpbyte(pb, []byte{1, 2}) {
					t.Errorf("error packa")
				}
			}

			pcn++
			fmt.Println("callback", pb)
		},
	)


	if err != nil{
		t.Errorf("error packa err:%s", err)
	}

	fmt.Println(surplus)
	if !cmpbyte(surplus, []byte{}) {
		t.Errorf("error packa")
	}


	if pcn != 2 {
		t.Errorf("error more")
	}


}


