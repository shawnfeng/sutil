// Copyright 2014 The dbrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbrouter

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/shawnfeng/sutil/slog/slog"
)

const DefaultGroup = ""

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
	Level  int32           `json:"level"`
	Ins    []itemDbInsCfg  `json:"ins"`
}

type itemDbInsCfg struct {
	Dbcfg json.RawMessage `json:"dbcfg"`
	Group string          `json:"group"`
}

type dbInsInfo struct {
	Instance string
	DBType   string
	DBName   string
	Level    int32
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
	dbIns map[string]map[string]*dbInsInfo
}

func (m *Parser) String() string {
	return fmt.Sprintf("%v", m.dbCls.clusters)
}

func (m *Parser) GetInstance(cluster, table string) string {
	instance := m.dbCls.getInstance(cluster, table)
	return instance
}

func (m *Parser) getConfig(instance, group string) *dbInsInfo {
	if infoMap, ok := m.dbIns[group]; ok {
		if info, ok := infoMap[instance]; ok {
			return info
		}
	} else {
		if info, ok := infoMap[DefaultGroup]; ok {
			return info
		}
	}
	return &dbInsInfo{}
}

func (m *Parser) GetConfig(instance, group string) *dbInsInfo {
	return m.getConfig(instance, group)
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
		dbIns: make(map[string]map[string]*dbInsInfo),
	}

	var cfg routeConfig
	err := json.Unmarshal(jscfg, &cfg)
	if err != nil {
		return nil, fmt.Errorf("dbrouter config unmarshal error: %s", err.Error())
	}

	cls := cfg.Cluster
	for c, ins := range cls {
		if er := checkVarname(c); er != nil {
			slog.Errorf(context.TODO(), "%s cluster config name err:%s", fun, err)
			continue
		}

		if len(ins) == 0 {
			slog.Warnf(context.TODO(), "%s empty instance in cluster:%s", fun, c)
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

	inss := cfg.Instances
	for ins, db := range inss {
		if er := checkVarname(ins); er != nil {
			slog.Errorf(context.TODO(), "%s instances name config err:%s", fun, er.Error())
			continue
		}

		dbtype := db.Dbtype
		dbname := db.Dbname
		cfg := db.Dbcfg
		level := db.Level
		// 外层配置为默认的group
		group := DefaultGroup
		info, err := parseDbIns(dbtype, dbname, ins, cfg, level)
		if err != nil {
			slog.Errorf(context.TODO(), "%s parse default dbInsInfo error, %s", fun, err.Error())
		} else {
			if _, ok := r.dbIns[group]; ok {
				r.dbIns[group][ins] = info
			} else {
				r.dbIns[group] = make(map[string]*dbInsInfo)
				r.dbIns[group][ins] = info
			}
		}

		for _, itemIns := range db.Ins {
			cfg = itemIns.Dbcfg
			group = itemIns.Group
			info, err := parseDbIns(dbtype, dbname, ins, cfg, level)
			if err != nil {
				slog.Errorf(context.TODO(), "%s parse %s dbInsInfo error, %s", fun, group, err.Error())
			} else {
				if _, ok := r.dbIns[group]; ok {
					r.dbIns[group][ins] = info
				} else {
					r.dbIns[group] = make(map[string]*dbInsInfo)
					r.dbIns[group][ins] = info
				}
			}
		}

		dbInsLen := len(r.dbIns[DefaultGroup])
		for group, insMap := range r.dbIns {
			if len(insMap) > dbInsLen {
				return nil, fmt.Errorf("db instance %s is lack of default group", ins)
			} else if len(insMap) < dbInsLen {
				return nil, fmt.Errorf("db instance %s is lack of group %s", ins, group)
			}
		}
	}

	return r, nil
}

func parseDbIns(dbtype, dbname, ins string, dbcfg json.RawMessage, level int32) (*dbInsInfo, error) {
	if er := checkVarname(dbtype); er != nil {
		return nil, fmt.Errorf("dbtype instance:%s err:%s", ins, er.Error())
	}

	if er := checkVarname(dbname); er != nil {
		return nil, fmt.Errorf("dbname instance:%s err:%s", ins, er.Error())
	}

	if len(dbcfg) == 0 {
		return nil, fmt.Errorf("empty dbcfg instance:%s", ins)
	}

	var info = new(dbInsInfo)
	err := json.Unmarshal(dbcfg, info)
	if err != nil {
		return nil, fmt.Errorf("unmarshal err, cfg:%s", string(dbcfg))
	}
	info.DBType = dbtype
	info.DBName = dbname
	info.Instance = ins
	info.Level = level

	return info, nil
}

func compareParsers(originParser Parser, newParser Parser) dbConfigChange {
	// 原来实例中修改的、删除的，要通知数据库连接池关闭掉实例的数据库连接
	var dbInsChangeMap = make(map[string][]string)
	var groups []string
	for group, originIns := range originParser.dbIns {
		var dbInsChanges []string
		if newIns, ok := newParser.dbIns[group]; ok {
			dbInsChanges = compareDbInstances(originIns, newIns)
			if len(dbInsChanges) == 0 {

			}
		} else {
			dbInsChanges = compareDbInstances(originIns, make(map[string]*dbInsInfo))
		}
		if len(dbInsChanges) > 0 {
			dbInsChangeMap[group] = dbInsChanges
		}
	}

	for group, _ := range newParser.dbIns {
		groups = append(groups, group)
	}

	return dbConfigChange{
		dbInstanceChange: dbInsChangeMap,
		dbGroups:         groups,
	}
}

func compareDbInstances(originDbInstances map[string]*dbInsInfo, newDbInstances map[string]*dbInsInfo) []string {
	var dbInstanceChanges []string
	for insName, originDbInsInfo := range originDbInstances {
		if newDbInsInfo, ok := newDbInstances[insName]; ok {
			if !compareDbInfo(originDbInsInfo, newDbInsInfo) {
				dbInstanceChanges = append(dbInstanceChanges, insName)
			}
		} else {
			dbInstanceChanges = append(dbInstanceChanges, insName)
		}
	}
	return dbInstanceChanges
}

func compareDbInfo(dbInsInfo1 *dbInsInfo, dbInsInfo2 *dbInsInfo) bool {
	return dbInsInfo1.DBName == dbInsInfo2.DBName && dbInsInfo1.UserName == dbInsInfo2.UserName &&
		dbInsInfo1.PassWord == dbInsInfo2.PassWord && compareStringList(dbInsInfo1.DBAddr, dbInsInfo2.DBAddr)
}

func compareStringList(stringList1 []string, stringList2 []string) bool {
	if len(stringList1) != len(stringList2) {
		return false
	}

	for index := range stringList1 {
		if stringList1[index] != stringList2[index] {
			return false
		}
	}

	return true
}
