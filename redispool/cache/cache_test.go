package rediscache

import (
	"testing"
	"encoding/json"

	"github.com/shawnfeng/sutil/slog"
)


type testData struct {
	Key string
	Tst int64
}

func (m *testData) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *testData) Unmarshal(sdata []byte) error {
	return json.Unmarshal(sdata, m)
}

func (m *testData) Load(key string) error {
	*m = testData {
		Key: key,
		Tst: 12345,
	}

	return nil
}


func TestCache(t *testing.T) {

	cachet := NewCache([]string{"127.0.0.1:9600"}, "TTT", 60)

	err := cachet.Del("test")
	if err != nil {
		t.Errorf("%s", err)
		return
	}


	var data testData
	err = cachet.Get("test", &data)
	if err != nil {
		t.Errorf("%s", err)
		return
	}


	slog.Infof("%s", data)

	if data.Key != "test" || data.Tst != 12345 {
		t.Errorf("get err")
	}


	err = cachet.Del("test")
	if err != nil {
		t.Errorf("%s", err)
		return
	}


}
