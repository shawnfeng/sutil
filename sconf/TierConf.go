package sconf

import (
	"bytes"
	"fmt"
	"strconv"
	"regexp"
	"strings"
	"io/ioutil"
	"sort"
	"github.com/vaughan0/go-ini"
)


type TierConf struct {
	reg *regexp.Regexp
	conf map[string]map[string]string

}

func NewTierConf() *TierConf {

	return &TierConf{
		conf: make(map[string]map[string]string),
		reg: regexp.MustCompile("\\$\\{.*?\\}"),

	}
}

func (m *TierConf) StringCheck() (string, error) {
	keys := make([]string, 0)
	for s, _ := range m.conf {
		keys = append(keys, s)
	}
	sort.Strings(keys)
	rv := ""
	for _, s := range keys {
		rv += fmt.Sprintf("[%s]\n", s)

		ps := make([]string, 0)
		for s, _ := range m.conf[s] {
			ps = append(ps, s)
		}
		sort.Strings(ps)

		for _, p := range ps {
			v, err := m.ToString(s, p)
			if err != nil {
				return "", err
			}
			rv += fmt.Sprintf("%s=%s\n", p, v)
		}
		rv += fmt.Sprintf("\n")
	}

	return rv, nil
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
	configs := strings.Split(conf, ",")

	for _, c := range configs {
		if err := m.LoadFromOneFile(c); err != nil {
			return err
		}
	}

	return nil

}


func (m *TierConf) LoadFromOneFile(conf string) error {
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

func (m *TierConf) toString(history []string, section string, property string) (string, error) {
	s, err := m.ToSection(section)

	if err != nil {
		return "", err

	} else {
		if p, ok := s[property]; ok {
			v, perr := m.parseVar(history, p)
			if perr != nil {
				return "", perr
			} else {
				return v, nil
			}
		} else {
			return "", fmt.Errorf("property empty:%s.%s", section, property)
		}


	}

}



func (m *TierConf) parseVar(history []string, value string) (string, error) {

	ids := m.reg.FindAllStringIndex(value, -1)

	var rv string = ""

	lastpos := 0
	for _, index := range(ids) {
		rv += value[lastpos:index[0]]
		pv := value[index[0]:index[1]]
		v := strings.Trim(pv, " \t${}")

		tmp := strings.Index(v, ".")

		if tmp == -1 {
			rv += pv
		} else {
			trims := strings.Trim(v[:tmp], " \t")
			trimp := strings.Trim(v[tmp+1:], " \t")
			// 检查循环引用
			newhis := fmt.Sprintf("%s.%s", trims, trimp)
			for _, his := range history {
				if newhis == his {
					return "", fmt.Errorf("cyclic reference:${%s}", his)
				}
			}
			history = append(history, newhis)
			newval, err := m.toString(history, trims, trimp)
			history = history[:len(history)-1]
			if err != nil {
				if strings.Index(err.Error(), "cyclic reference") != -1 {
					return "", err
				} else {
					newval = pv
				}
			}

			rv += newval

			//fmt.Println(v[:tmp], v[tmp:], pv, ids, history)
		}

		lastpos = index[1]
	}

	rv += value[lastpos:]

	return rv, nil

}


func (m *TierConf) ToSection(section string) (map[string]string, error) {
	if s, ok := m.conf[section]; ok {
		return s, nil
	} else {
		return nil, fmt.Errorf("section empty:%s", section)
	}


}

func (m *TierConf) ToString(section string, property string) (string, error) {
	return m.toString(nil, section, property)

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

func (m *TierConf) ToInt32(section string, property string) (int32, error) {
	v, err := m.ToString(section, property)

	if err != nil {
		return 0, err
	} else {
		i, err := strconv.ParseInt(v, 10, 32)
		return int32(i), err
	}

}

func (m *TierConf) ToInt64(section string, property string) (int64, error) {
	v, err := m.ToString(section, property)

	if err != nil {
		return 0, err
	} else {
		i, err := strconv.ParseInt(v, 10, 64)
		return i, err
	}

}

func (m *TierConf) ToUint64(section string, property string) (uint64, error) {
	v, err := m.ToString(section, property)

	if err != nil {
		return 0, err
	} else {
		i, err := strconv.ParseUint(v, 10, 64)
		return i, err
	}

}

func (m *TierConf) ToUint32(section string, property string) (uint32, error) {
	v, err := m.ToString(section, property)

	if err != nil {
		return 0, err
	} else {
		i, err := strconv.ParseUint(v, 10, 32)
		return uint32(i), err
	}

}

func (m *TierConf) ToFloat64(section string, property string) (float64, error) {
	v, err := m.ToString(section, property)

	if err != nil {
		return 0, err
	} else {
		i, err := strconv.ParseFloat(v, 64)
		return i, err
	}

}

func (m *TierConf) ToFloat32(section string, property string) (float32, error) {
	v, err := m.ToString(section, property)

	if err != nil {
		return 0, err
	} else {
		i, err := strconv.ParseFloat(v, 32)
		return float32(i), err
	}

}


func (m *TierConf) ToBool(section string, property string) (bool, error) {
	v, err := m.ToString(section, property)

	if err != nil {
		return false, err
	} else {
		return strconv.ParseBool(v)
	}

}


func (m *TierConf) ToBoolWithDefault(section string, property string, deft bool) bool {
	v, err := m.ToBool(section, property)

	if err != nil {
		return deft
	} else {
		return v
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
