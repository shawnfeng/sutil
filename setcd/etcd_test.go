package setcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/client"
	"math/rand"
	"testing"
	"time"
)

var etcdIns *EtcdInstance

func init() {
	instance, e := NewEtcdInstance([]string{"http://infra0.etcd.ibanyu.com:20002", "http://infra1.etcd.ibanyu.com:20002"})
	//instance, e := NewEtcdInstance([]string{"http://localhost:20002"})
	if e != nil {
		panic("create etcdIns panic")
	}
	etcdIns = instance
}

func testpath(path string) string {
	return fmt.Sprintf("/test%s", path)
}

func TestEtcdInstance_Get(t *testing.T) {
	s, e := etcdIns.Get(context.Background(), "/roc/db/route")
	if e != nil {
		t.Error(e)
		return
	}
	t.Log(s)
	t.Log("=====")

}

func Test(t *testing.T) {
	t.Log("=====")

}

func TestEtcdInstance_CreateDir(t *testing.T) {
	s, e := etcdIns.Get(context.Background(), testpath("/d1"))
	if e != nil {
		t.Error(e)
		return
	}
	t.Log(s)
	e = etcdIns.CreateDir(context.Background(), testpath("/d1"))
	if e != nil {
		t.Error(e)
		return
	}
	t.Log("===")

	e = etcdIns.CreateDir(context.Background(), testpath("/d1"))
	if e != nil {
		t.Error("exist", e)
		return
	}
	t.Log("===")
	return
}
func TestEtcdInstance_GetNode(t *testing.T) {
	node, e := etcdIns.GetNode(context.Background(), testpath("/d1"))
	if e != nil {
		t.Error(e)
		return
	}
	bytes, i := json.Marshal(node)
	t.Log("===", string(bytes), i)
	return
}
func TestEtcdInstance_RefreshTtl(t *testing.T) {
	path := testpath("/d1/r1")
	e := etcdIns.SetTtl(context.Background(), path, "r111", time.Second*5)
	node, e := etcdIns.GetNode(context.Background(), path)
	t.Log("5s", node.TTL)
	time.Sleep(1 * time.Second)
	node, e = etcdIns.GetNode(context.Background(), path)
	t.Log("4s", node.TTL)
	e = etcdIns.RefreshTtl(context.Background(), path, time.Second*10)
	if e != nil {
		t.Error(e)
		return
	}
	node, e = etcdIns.GetNode(context.Background(), path)
	t.Log("10s", node.TTL)
	e = etcdIns.SetTtl(context.Background(), path, "r111", time.Second*5)
	node, e = etcdIns.GetNode(context.Background(), path)
	t.Log("5s", node.TTL)
	t.Log("===")
	return
}

func TestEtcdInstance_Set(t *testing.T) {
	path := testpath("/d1/c1")
	e := etcdIns.Set(context.Background(), path, "c11111")
	if e != nil {
		t.Error(e)
		return
	}

	s, e := etcdIns.Get(context.Background(), path)
	if e != nil {
		t.Error(e)
		return
	}
	t.Log("===", s, e)

	path = testpath("/d1/c1/f1")
	e = etcdIns.Set(context.Background(), path, "f11111")
	if e != nil {
		t.Error(e)
		return
	}

	s, e = etcdIns.Get(context.Background(), path)
	if e != nil {
		t.Error(e)
		return
	}
	t.Log("===", s, e)
	return
}
func TestEtcdInstance_SetNx(t *testing.T) {
	path := testpath("/d1/e1")
	e := etcdIns.SetNx(context.Background(), path, "e111111")
	if e != nil {
		t.Error(e)
		return
	}

	s, e := etcdIns.Get(context.Background(), path)
	if e != nil {
		t.Error(e)
		return
	}
	t.Log("===", s, e)
	return
}
func TestEtcdInstance_SetTtl(t *testing.T) {
	path := testpath("/d1/e3")
	e := etcdIns.SetTtl(context.Background(), path, "e333", time.Second*3)
	if e != nil {
		t.Error(e)
		return
	}
	s, i := etcdIns.Get(context.Background(), path)
	t.Log("===", s, i)
	time.Sleep(time.Second * 4)
	s, i = etcdIns.Get(context.Background(), path)
	t.Log("xxxx", s, i)
	return
}
func TestEtcdInstance_Regist(t *testing.T) {
	path := testpath("/d1/r2")
	e := etcdIns.Set(context.Background(), path, "rrr")
	e = etcdIns.Regist(context.Background(), path, "rrr1", time.Second*3, time.Second*10)
	if e != nil {
		t.Error(e)
		return
	}

	go func() {
		for {
			s, _ := etcdIns.GetNode(context.Background(), path)
			fmt.Println("get=====", s.Value, s.TTL)
			time.Sleep(1 * time.Second)
		}
	}()
	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(5 * time.Second)
			e := etcdIns.Set(context.Background(), path, fmt.Sprintf("ccc%d", i))
			fmt.Println("set===", e)
		}
	}()

	//go func() {
	//	for {
	//		v := fmt.Sprintf("%d", rand.Int31n(1000))
	//		e := etcdIns.Set(context.Background(), path,v)
	//		t.Log("set=====",e,v)
	//		time.Sleep(1*time.Second)
	//	}
	//}()

	t.Log("===")
	time.Sleep(10 * time.Hour)
	return
}
func TestEtcdInstance_Watch(t *testing.T) {
	path := testpath("/d1/w2")
	e := etcdIns.CreateDir(context.Background(), path)
	if e != nil {
		t.Error(e)
		return
	}
	etcdIns.Watch(context.Background(), path, func(response *client.Response) {
		bytes, _ := json.Marshal(response.Node)
		fmt.Println("watch change ", response.Action, string(bytes))
	})
	go func() {
		for {
			s, _ := etcdIns.GetNode(context.Background(), path)
			fmt.Println("get=====", s.Value, s.TTL)
			time.Sleep(1 * time.Second)
		}
	}()
	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(time.Duration(rand.Int31n(5)+1) * time.Second)
			tmppath := path
			if i%1 == 0 {
				tmppath = fmt.Sprintf("%s/q1", path)
			}
			e := etcdIns.Set(context.Background(), tmppath, fmt.Sprintf("ccc%d", i))
			fmt.Println("set===", tmppath, e)
		}
	}()
	t.Log("===")
	time.Sleep(10 * time.Hour)
	return
}
