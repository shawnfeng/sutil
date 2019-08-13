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

type dbInsCfgEtcd struct {
	Dbtype string          `json:"dbtype"`
	Dbname string          `json:"dbname"`
	Dbcfg  json.RawMessage `json:"dbcfg"`
	Type   int32           `json:"type"`
	// TODO 可能要新添加压测配置
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

type routeConfigEtcd struct {
	Cluster   map[string][]*dbLookupCfg  `json:"cluster"`
	Instances map[string]*dbInsCfgEtcd   `json:"instances"`
}

type Parser struct {
	dbCls *dbCluster
	dbIns map[string]*dbInsInfo
	shadowDbIns map[string]*dbInsInfo  // 影子实例配置，用于全链路压测
}

func (m *Parser) String() string {
	return fmt.Sprintf("%v", m.dbCls.clusters)
}

func (m *Parser) GetInstance(cluster, table string) string {
	instance := m.dbCls.getInstance(cluster, table)
	return instance
	/*info := m.getConfig(instance)
	return info.Instance*/
}

func (m *Parser) getConfig(instance string) *dbInsInfo {
	if info, ok := m.dbIns[instance]; ok {
		return info
	}
	return &dbInsInfo{}
}

func (m *Parser) GetShadowConfig(instance string) *dbInsInfo {
	if info, ok := m.shadowDbIns[instance]; ok {
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

func NewParserEtcd(jscfg []byte) (*Parser, error) {
	fun := "NewParserEtcd -->"

	r := &Parser{
		dbCls: &dbCluster{
			clusters: make(map[string]*clsEntry),
		},
		dbIns: make(map[string]*dbInsInfo),
	}

	var cfg routeConfigEtcd
	err := json.Unmarshal(jscfg, &cfg)
	if err != nil {
		slog.Errorf(context.TODO(), "%s dbrouter config unmarshal:%s", fun, err.Error())
		return r, nil
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

	inss := cfg.Instances
	for ins, db := range inss {
		if er := checkVarname(ins); er != nil {
			slog.Errorf(context.TODO(), "%s instances name config err:%s", fun, er.Error())
			continue
		}

		// for _, db := range dbs {
			dbtype := db.Dbtype
			dbname := db.Dbname
			cfg := db.Dbcfg

			if er := checkVarname(dbtype); er != nil {
				slog.Errorf(context.TODO(), "%s dbtype instance:%s err:%s", fun, ins, er.Error())
				continue
			}

			if er := checkVarname(dbname); er != nil {
				slog.Errorf(context.TODO(), "%s dbname instance:%s err:%s", fun, ins, er.Error())
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

			// TODO 标识字段待定
			if db.Type == 0 {
				if _, ok := r.dbIns[ins]; ok {
					slog.Errorf(context.TODO(), "%s dbname dup, ins:%s, cfg:%v", fun, ins, string(cfg))
					continue
				}

				r.dbIns[ins] = &info
			} else if db.Type == 1 {
				if _, ok := r.shadowDbIns[ins]; ok {
					slog.Errorf(context.TODO(), "%s shadow dbname dup, ins:%s, cfg:%v", fun, ins, string(cfg))
					continue
				}

				r.shadowDbIns[ins] = &info
			}
		}
	// }

	return r, nil
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
			slog.Errorf(context.TODO(), "%s instances name config err:%s", fun, er.Error())
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
			slog.Errorf(context.TODO(), "%s dbname instance:%s err:%s", fun, ins, er.Error())
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

func compareParsers(originParser Parser, newParser Parser) dbInstanceChange {
	// 原来实例中修改的、删除的，要通知数据库连接池关闭掉实例的数据库连接
	dbInsChanges := compareDbInstances(originParser.dbIns, newParser.dbIns)
	shadowDbInsChanges := compareDbInstances(originParser.shadowDbIns, newParser.shadowDbIns)
	return dbInstanceChange{
		dbInsChanges: dbInsChanges,
		shadowDbInsChanges: shadowDbInsChanges,
	}
}

func compareDbInstances(originDbInstances map[string]*dbInsInfo, newDbInstances map[string]*dbInsInfo) []string {
	var dbInstanceChanges []string
	for insName, originDbInsInfo := range originDbInstances {
		if newDbInsInfo, ok := newDbInstances[insName]; ok {
			if !compareDbInfo(*originDbInsInfo, *newDbInsInfo) {
				dbInstanceChanges = append(dbInstanceChanges, insName)
			}
		} else {
			dbInstanceChanges = append(dbInstanceChanges, insName)
		}
	}
	return dbInstanceChanges
}

func compareDbInfo(dbInsInfo1 dbInsInfo, dbInsInfo2 dbInsInfo) bool {
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
