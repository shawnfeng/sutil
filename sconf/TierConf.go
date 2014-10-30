package sconf

import (
	"bytes"
	"github.com/vaughan0/go-ini"
)


type TierConf struct {
	res map[string]map[string]string

}

func NewTierConf() *TierConf {

	return &TierConf{
		res: make(map[string]map[string]string),

	}
}

func (m *TierConf) GetConf() map[string]map[string]string {
	return m.res
}

func (m *TierConf) Load(cfg []byte) error {
	file, err := ini.Load(bytes.NewReader(cfg))

	if err != nil {
		return err
	}

	for name, section := range file {
		if _, ok := m.res[name]; !ok {
			m.res[name] = make(map[string]string)
		}

		for k, v := range section {
			m.res[name][k] = v
		}
	}


	return nil

}

