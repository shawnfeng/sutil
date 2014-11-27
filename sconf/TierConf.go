package sconf

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"io/ioutil"
	"github.com/vaughan0/go-ini"
)


type TierConf struct {
	conf map[string]map[string]string

}

func NewTierConf() *TierConf {

	return &TierConf{
		conf: make(map[string]map[string]string),

	}
}

func (m *TierConf) GetConf() map[string]map[string]string {
	return m.conf
}

func (m *TierConf) LoadFromConf(cfg map[string]map[string]string) {
	for name, section := range cfg {
		if _, ok := m.conf[name]; !ok {
			m.conf[name] = make(map[string]string)
		}

		for k, v := range section {
			m.conf[name][k] = v
		}
	}


}

func (m *TierConf) LoadFromFile(conf string) error {
	data, err := ioutil.ReadFile(conf)
	if err != nil {
		return err
	} else {
		return m.Load(data)
	}


}

func (m *TierConf) Load(cfg []byte) error {
	file, err := ini.Load(bytes.NewReader(cfg))

	if err != nil {
		return err
	}

	for name, section := range file {
		if _, ok := m.conf[name]; !ok {
			m.conf[name] = make(map[string]string)
		}

		for k, v := range section {
			m.conf[name][k] = v
		}
	}


	return nil

}


func (m *TierConf) ToSection(section string) (map[string]string, error) {
	if s, ok := m.conf[section]; ok {
		return s, nil
	} else {
		return nil, errors.New("section empty")
	}


}

func (m *TierConf) ToString(section string, property string) (string, error) {
	s, err := m.ToSection(section)

	if err != nil {
		return "", err

	} else {
		if p, ok := s[property]; ok {
			return p, nil
		} else {
			return "", errors.New("property empty")
		}


	}

}

func (m *TierConf) ToStringWithDefault(section string, property string, deft string) string {
	v, err := m.ToString(section, property)

	if err != nil {
		return deft
	} else {
		return v
	}


}

func (m *TierConf) ToInt(section string, property string) (int, error) {
	v, err := m.ToString(section, property)

	if err != nil {
		return 0, err
	} else {
		return strconv.Atoi(v)
	}

}

func (m *TierConf) ToIntWithDefault(section string, property string, deft int) int {
	v, err := m.ToInt(section, property)
	if err != nil {
		return deft

	} else {
		return v
	}
}

func (m *TierConf) ToSliceString(section string, property string, sep string) ([]string, error) {
	v, err := m.ToString(section, property)

	if err != nil {
		return nil, err
	} else {
		ss := strings.Split(v, sep)

		for i := 0; i < len(ss); i++ {
			ss[i] = strings.Trim(ss[i], " \t")
		}
		return ss, nil
	}

}

func (m *TierConf) ToSliceInt(section string, property string, sep string) ([]int, error) {
	s, err := m.ToSliceString(section, property, sep)
	if err != nil {
		return nil, err
	} else {
		ints := make([]int, 0)
		for _, v := range(s) {
			tmp, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			} else {
				ints = append(ints, tmp)
			}

		}

		return ints, nil
	}

}
