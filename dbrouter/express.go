// Copyright 2014 The dbrouter Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbrouter

import (
	"fmt"
	"regexp"
)

type clsEntry struct {
	full  map[string]*dbExpress
	regex map[string]*dbExpress
}

type dbCluster struct {
	clusters map[string]*clsEntry
}

type dbExpress struct {
	lookup *dbLookupCfg
	reg    *regexp.Regexp
}

func (m *dbExpress) String() string {
	return fmt.Sprintf("look:%s reg:%s", m.lookup, m.reg)
}

func (m *dbCluster) addInstance(cluster string, lcfg *dbLookupCfg) error {
	if _, ok := m.clusters[cluster]; !ok {
		m.clusters[cluster] = &clsEntry{
			full:  make(map[string]*dbExpress),
			regex: make(map[string]*dbExpress),
		}
	}

	match := lcfg.Match
	if match == "full" {
		if m.clusters[cluster].full[lcfg.Express] != nil {
			return fmt.Errorf("dup match full in cluster:%s express:%s", cluster, lcfg.Express)
		}

		m.clusters[cluster].full[lcfg.Express] = &dbExpress{lookup: lcfg}

	} else if match == "regex" {
		if m.clusters[cluster].regex[lcfg.Express] != nil {
			return fmt.Errorf("dup match regex in cluster:%s express:%s", cluster, lcfg.Express)
		}

		reg, err := regexp.CompilePOSIX(lcfg.Express)
		if err != nil {
			return err
		}

		m.clusters[cluster].regex[lcfg.Express] = &dbExpress{lookup: lcfg, reg: reg}

	} else {
		return fmt.Errorf("match type:%s not support", match)
	}

	return nil
}

func (m *dbCluster) getLookup(cluster string, table string) *dbLookupCfg {

	exp := m.clusters[cluster]
	if exp == nil {
		return nil
	}

	// 先全匹配查找
	en := exp.full[table]
	if en != nil {
		return en.lookup
	}

	// 正则
	for _, e := range exp.regex {
		// 必须全部匹配上
		f := e.reg.FindString(table)
		if table == f {
			return e.lookup
		}
	}

	return nil
}

func (m *dbCluster) getInstance(cluster string, table string) string {
	if lk := m.getLookup(cluster, table); lk != nil {
		return lk.Instance
	} else {
		return ""
	}

}
