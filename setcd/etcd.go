package setcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/client"
	"github.com/shawnfeng/sutil/slog"
	"github.com/shawnfeng/sutil/stime"
	"time"
)

type EtcdInstance struct {
	Api client.KeysAPI
}

func NewEtcdInstanceWichApi(api client.KeysAPI) *EtcdInstance {
	return &EtcdInstance{
		Api: api,
	}
}

func NewEtcdInstance(cluster []string) (*EtcdInstance, error) {
	cfg := client.Config{
		Endpoints: cluster,
		Transport: client.DefaultTransport,
	}
	c, err := client.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("create etchd client cfg error")
	}
	api := client.NewKeysAPI(c)
	if api == nil {
		return nil, fmt.Errorf("create etchd api error")
	}
	return NewEtcdInstanceWichApi(api), nil
}
func (m *EtcdInstance) Get(ctx context.Context, path string) (string, error) {
	r, err := m.Api.Get(ctx, path, &client.GetOptions{
		Recursive: false,
		Sort:      false,
	})
	if err != nil {
		return "", err
	}

	if r.Node == nil {
		return "", fmt.Errorf("etcdIns node value err location:%s", path)
	}

	return r.Node.Value, nil
}
func (m *EtcdInstance) GetNode(ctx context.Context, path string) (*client.Node, error) {
	r, err := m.Api.Get(ctx, path, &client.GetOptions{
		Recursive: true,
		Sort:      false,
	})
	if err != nil {
		return nil, err
	}

	if r.Node == nil {
		return nil, fmt.Errorf("etcdIns node value err location:%s", path)
	}

	return r.Node, nil
}
func (m *EtcdInstance) Set(ctx context.Context, path, val string) error {
	r, err := m.Api.Set(ctx, path, val, &client.SetOptions{})
	if err != nil {
		return err
	}

	if r.Node == nil {
		return fmt.Errorf("etcdIns node value err location:%s", path)
	}

	return nil
}
func (m *EtcdInstance) CreateDir(ctx context.Context, path string) error {
	_, err := m.Api.Set(ctx, path, "", &client.SetOptions{
		Dir:       true,
		PrevExist: client.PrevNoExist,
	})
	if err != nil {
		return err
	}
	return nil
}
func (m *EtcdInstance) SetTtl(ctx context.Context, path, val string, ttl time.Duration) error {
	_, err := m.Api.Set(ctx, path, val, &client.SetOptions{
		TTL: ttl,
	})
	if err != nil {
		return err
	}

	return nil
}
func (m *EtcdInstance) RefreshTtl(ctx context.Context, path string, ttl time.Duration) error {

	_, err := m.Api.Set(ctx, path, "", &client.SetOptions{
		PrevExist: client.PrevExist,
		Refresh:   true,
		TTL:       ttl,
	})
	if err != nil {
		return err
	}

	return nil
}
func (m *EtcdInstance) SetNx(ctx context.Context, path, val string) error {
	_, err := m.Api.Set(ctx, path, val, &client.SetOptions{
		PrevExist: client.PrevNoExist,
	})
	if err != nil {
		return err
	}

	return nil
}
func (m *EtcdInstance) Regist(ctx context.Context, path, val string, heatbeat time.Duration, ttl time.Duration) error {
	fun := "EtcdInstance.Regist -->"
	var isset = true
	go func() {
		for i := 0; ; i++ {
			var err error
			if isset {
				slog.Warnf(ctx, "%s create idx:%d val:%s", fun, i, val)
				_, err = m.Api.Set(ctx, path, val, &client.SetOptions{
					TTL: ttl,
				})
				if err == nil {
					isset = false
				}
			} else {
				slog.Infof(ctx, "%s refresh ttl idx:%d val:%s", fun, i, val)
				_, err = m.Api.Set(ctx, path, "", &client.SetOptions{
					PrevExist: client.PrevExist,
					TTL:       ttl,
					Refresh:   true,
				})
			}
			if err != nil {
				slog.Errorf(ctx, "%s reg idx:%d err:%s", fun, i, err)

			}

			time.Sleep(heatbeat)
		}
	}()

	return nil
}
func (m *EtcdInstance) Watch(ctx context.Context, path string, hander func(*client.Response)) {
	fun := "EtcdInstance.Watch -->"
	backoff := stime.NewBackOffCtrl(time.Millisecond*10, time.Second*5)
	var chg chan *client.Response
	go func() {
		slog.Infof(ctx, "%s start watch:%s", fun, path)
		for {
			if chg == nil {
				slog.Infof(ctx, "%s loop watch new receiver:%s", fun, path)
				chg = make(chan *client.Response)
				go m.startWatch(ctx, chg, path)
			}

			r, ok := <-chg
			if !ok {
				slog.Errorf(ctx, "%s chg info nil:%s", fun, path)
				chg = nil
				backoff.BackOff()
			} else {
				slog.Infof(ctx, "%s update path:%s", fun, r.Node.Key)
				hander(r)
				backoff.Reset()
			}
		}
	}()
}

func (m *EtcdInstance) startWatch(ctx context.Context, chg chan *client.Response, path string) {
	fun := "EtcdInstance.startWatch -->"

	for i := 0; ; i++ {
		r, err := m.Api.Get(ctx, path, &client.GetOptions{Recursive: true, Sort: false})
		if err != nil {
			slog.Warnf(ctx, "%s get path:%s err:%s", fun, path, err)
		} else {
			chg <- r
		}
		index := uint64(0)
		if r != nil {
			index = r.Index
			sresp, _ := json.Marshal(r)
			slog.Infof(ctx, "%s init get action:%s nodes:%d index:%d path:%s resp:%s", fun, r.Action, len(r.Node.Nodes), r.Index, path, sresp)
		}

		wop := &client.WatcherOptions{
			Recursive:  true,
			AfterIndex: index,
		}
		watcher := m.Api.Watcher(path, wop)
		if watcher == nil {
			slog.Errorf(ctx, "%s new watcher path:%s", fun, path)
			return
		}

		resp, err := watcher.Next(context.Background())
		// etcdIns 关闭时候会返回
		if err != nil {
			slog.Errorf(ctx, "%s watch path:%s err:%s", fun, path, err)
			close(chg)
			return
		} else {
			slog.Infof(ctx, "%s next get idx:%d action:%s nodes:%d index:%d after:%d path:%s", fun, i, resp.Action, len(resp.Node.Nodes), resp.Index, wop.AfterIndex, path)
			// 测试发现next获取到的返回，index，重新获取总有问题，触发两次，不确定，为什么？为什么？
			// 所以这里每次next前使用的afterindex都重新get了
		}
	}

}
