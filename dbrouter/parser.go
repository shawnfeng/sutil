// Copyright 2014 The dbrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbrouter

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/shawnfeng/sutil/slog"
)

type dbLookupCfg struct {
	Instance string `json:"instance"`
	Match    string `json:"match"`
	Express  string `json:"express"`
}

func (m *dbLookupCfg) String() string {
	return fmt.Sprintf("ins:%s exp:%s match:%s", m.Instance, m.Express, m.Match)
}

type dbInsCfg struct {
	Dbtype string          `json:"dbtype"`
	Dbname string          `json:"dbname"`
	Dbcfg  json.RawMessage `json:"dbcfg"`
}

type dbInsInfo struct {
	Instance string
	DBType   string
	DBName   string
	DBAddr   []string `json:"addrs"`
	UserName string   `json:"user"`
	PassWord string   `json:"passwd"`
}

type routeConfig struct {
	Cluster   map[string][]*dbLookupCfg `json:"cluster"`
	Instances map[string]*dbInsCfg      `json:"instances"`
}

type Parser struct {
	dbCls *dbCluster
	dbIns map[string]*dbInsInfo
}

func (m *Parser) String() string {
	return fmt.Sprintf("%s", m.dbCls.clusters)
}

func (m *Parser) GetInstance(cluster, table string) string {
	instance := m.dbCls.getInstance(cluster, table)
	info := m.getConfig(instance)
	return info.Instance
}

func (m *Parser) getConfig(instance string) *dbInsInfo {
	if info, ok := m.dbIns[instance]; ok {
		return info
	}
	return &dbInsInfo{}
}

func (m *Parser) GetConfig(instance string) *dbInsInfo {
	return m.getConfig(instance)
}

// 检查用户输入的合法性
// 1. 只能是字母或者下划线
// 2. 首字母不能为数字，或者下划线
func checkVarname(varname string) error {
	if len(varname) == 0 {
		return fmt.Errorf("is empty")
	}

	f := varname[0]
	if !((f >= 'a' && f <= 'z') || (f >= 'A' && f <= 'Z')) {
		return fmt.Errorf("first char is not alpha")
	}

	for _, c := range varname {

		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			continue
		} else if c >= '0' && c <= '9' {
			continue
		} else if c == '_' {
			continue
		} else {
			return fmt.Errorf("is contain not [a-z] or [A-Z] or [0-9] or _")
		}
	}

	return nil
}

func NewParser(jscfg []byte) (*Parser, error) {
	fun := "NewParser -->"

	r := &Parser{
		dbCls: &dbCluster{
			clusters: make(map[string]*clsEntry),
		},
		dbIns: make(map[string]*dbInsInfo),
	}

	var cfg routeConfig
	err := json.Unmarshal(jscfg, &cfg)
	if err != nil {
		slog.Errorf(context.TODO(), "%s dbrouter config unmarshal:%s", fun, err.Error())
		return r, nil
	}

	inss := cfg.Instances
	for ins, db := range inss {
		if er := checkVarname(ins); er != nil {
			slog.Errorf(context.TODO(),"%s instances name config err:%s", fun, err.Error())
			continue
		}

		dbtype := db.Dbtype
		dbname := db.Dbname
		cfg := db.Dbcfg

		if er := checkVarname(dbtype); er != nil {
			slog.Errorf(context.TODO(), "%s dbtype instance:%s err:%s", fun, ins, er.Error())
			continue
		}

		if er := checkVarname(dbname); er != nil {
			slog.Errorf(context.TODO(), "%sdbname instance:%s err:%s", fun, ins, er.Error())
			continue
		}

		if len(cfg) == 0 {
			slog.Errorf(context.TODO(), "%s empty dbcfg instance:%s", fun, ins)
			continue
		}

		var info dbInsInfo
		err := json.Unmarshal(cfg, &info)
		if err != nil {
			slog.Errorf(context.TODO(), "%s unmarshal err, cfg:%s", fun, string(cfg))
			continue
		}
		info.DBType = dbtype
		info.DBName = dbname
		info.Instance = ins

		if _, ok := r.dbIns[ins]; ok {
			slog.Errorf(context.TODO(), "%s dbname dup, ins:%s, cfg:%v", fun, ins, string(cfg))
			continue
		}

		r.dbIns[ins] = &info
	}

	cls := cfg.Cluster
	for c, ins := range cls {
		if er := checkVarname(c); er != nil {
			slog.Errorf(context.TODO(), "%s cluster config name err:%s", fun, err)
			continue
		}

		if len(ins) == 0 {
			slog.Errorf(context.TODO(), "%s empty instance in cluster:%s", fun, c)
			continue
		}

		for _, v := range ins {
			if len(v.Express) == 0 {
				slog.Errorf(context.TODO(), "%s empty express in cluster:%s instance:%s", fun, c, v.Instance)
				continue
			}

			if er := checkVarname(v.Match); er != nil {
				slog.Errorf(context.TODO(), "%s match in cluster:%s instance:%s err:%s", fun, c, v.Instance, err)
				continue
			}

			if er := checkVarname(v.Instance); er != nil {
				slog.Errorf(context.TODO(), "%s instance name in cluster:%s instance:%s err:%s", fun, c, v.Instance, err)
				continue
			}

			if err := r.dbCls.addInstance(c, v); err != nil {
				return nil, fmt.Errorf("load instance lookup rule err:%s", err.Error())
			}
		}
	}

	return r, nil
}
