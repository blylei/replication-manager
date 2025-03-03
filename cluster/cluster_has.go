// replication-manager - Replication Manager Monitoring and CLI for MariaDB and MySQL
// Copyright 2017 Signal 18 SARL
// Authors: Guillaume Lefranc <guillaume@signal18.io>
//          Stephane Varoqui  <svaroqui@gmail.com>
// This source code is licensed under the GNU General Public License, version 3.

package cluster

import (
	"strings"

	"github.com/signal18/replication-manager/config"
)

func (cluster *Cluster) HasServer(srv *ServerMonitor) bool {
	for _, sv := range cluster.Servers {
		//	cluster.LogPrintf(LvlInfo, "HasServer:%s %s, %s %s", sv.Id, srv.Id, sv.URL, srv.URL)
		// id can not be used for checking equality because  same srv in different clusters
		if sv.URL == srv.URL {
			return true
		}
	}
	return false
}

func (cluster *Cluster) HasSchedulerEntry(myname string) bool {
	if _, ok := cluster.Schedule[myname]; ok {
		return true
	}

	return false
}

func (cluster *Cluster) IsProvisioned() bool {
	if cluster.Conf.ProvOrchestrator == config.ConstOrchestratorOnPremise {
		return true
	}
	if cluster.Conf.Hosts == "" {
		return false
	}
	for _, db := range cluster.Servers {
		if !db.HasProvisionCookie() {
			if db.IsRunning() {
				db.SetProvisionCookie()
				cluster.LogPrintf(LvlInfo, "Can DB Connect creating cookie state:%s", db.State)
			} else {
				return false
			}
		}
	}
	for _, px := range cluster.Proxies {
		if !px.HasProvisionCookie() {
			if px.IsRunning() {
				px.SetProvisionCookie()
				cluster.LogPrintf(LvlInfo, "Can Proxy Connect creating cookie state:%s", px.State)
			} else {
				return false
			}
		}
	}
	return true
}

func (cluster *Cluster) IsInIgnoredHosts(server *ServerMonitor) bool {
	ihosts := strings.Split(cluster.Conf.IgnoreSrv, ",")
	for _, ihost := range ihosts {
		if server.URL == ihost || server.Name == ihost {
			return true
		}
	}
	return false
}

func (cluster *Cluster) IsInPreferedBackupHosts(server *ServerMonitor) bool {
	ihosts := strings.Split(cluster.Conf.BackupServers, ",")
	for _, ihost := range ihosts {
		if server.URL == ihost || server.Name == ihost {
			return true
		}
	}
	return false
}

func (cluster *Cluster) IsInIgnoredReadonly(server *ServerMonitor) bool {
	ihosts := strings.Split(cluster.Conf.IgnoreSrvRO, ",")
	for _, ihost := range ihosts {
		if server.URL == ihost || server.Name == ihost {
			return true
		}
	}
	return false
}

func (cluster *Cluster) IsInPreferedHosts(server *ServerMonitor) bool {
	ihosts := strings.Split(cluster.Conf.PrefMaster, ",")
	for _, ihost := range ihosts {
		if server.URL == ihost || server.Name == ihost {
			return true
		}
	}
	return false
}

func (cluster *Cluster) IsInCaptureMode() bool {
	if !cluster.Conf.MonitorCapture || cluster.IsNotMonitoring || len(cluster.Servers) > 0 {
		return false
	}
	for _, server := range cluster.Servers {
		if server.InCaptureMode {
			return true
		}
	}
	return false
}

func (cluster *Cluster) HasAllDbUp() bool {
	if cluster.Servers == nil {
		return false
	}
	for _, s := range cluster.Servers {
		if s.State == stateFailed || s.State == stateSuspect /*&& misc.Contains(cluster.ignoreList, s.URL) == false*/ {
			return false
		}
	}
	return true
}

func (cluster *Cluster) HasRequestDBRestart() bool {
	if cluster.Servers == nil {
		return false
	}
	for _, s := range cluster.Servers {
		if s.HasRestartCookie() {
			return true
		}
	}
	return false
}

func (cluster *Cluster) HasRequestDBRollingRestart() bool {
	ret := true
	if cluster.Servers == nil {
		return false
	}
	for _, s := range cluster.Servers {
		if !s.HasRestartCookie() {
			return false
		}
	}
	return ret
}

func (cluster *Cluster) HasRequestDBRollingReprov() bool {
	ret := true
	if cluster.Servers == nil {
		return false
	}
	for _, s := range cluster.Servers {
		if !s.HasReprovCookie() {
			return false
		}
	}
	return ret
}

func (cluster *Cluster) HasRequestDBReprov() bool {
	for _, s := range cluster.Servers {
		if s.HasReprovCookie() {
			return true
		}
	}
	return false
}

func (cluster *Cluster) HasRequestProxiesRestart() bool {
	for _, p := range cluster.Proxies {
		if p.HasRestartCookie() {
			return true
		}
	}
	return false
}

func (cluster *Cluster) HasRequestProxiesReprov() bool {
	for _, p := range cluster.Proxies {
		if p.HasReprovCookie() {
			return true
		}
	}
	return false
}

func (cluster *Cluster) IsInHostList(host string) bool {
	for _, v := range cluster.hostList {
		if v == host {
			return true
		}
	}
	return false
}

func (cluster *Cluster) IsMasterFailed() bool {
	// get real master or the virtual master
	mymaster := cluster.GetMaster()
	if mymaster == nil {
		return true
	}
	if mymaster.State == stateFailed {
		return true
	} else {
		return false
	}
}

func (cluster *Cluster) IsActive() bool {
	if cluster.Status == ConstMonitorActif {
		return true
	} else {
		return false
	}
}

func (cluster *Cluster) IsVerbose() bool {
	if cluster.Conf.Verbose {
		return true
	} else {
		return false
	}
}

func (cluster *Cluster) IsInFailover() bool {
	return cluster.sme.IsInFailover()
}

func (cluster *Cluster) IsDiscovered() bool {
	return cluster.sme.IsDiscovered()
}

func (cluster *Cluster) HaveDBTag(tag string) bool {
	for _, t := range cluster.DBTags {
		if t == tag {
			return true
		}
	}
	return false
}

func (cluster *Cluster) HaveProxyTag(tag string) bool {
	for _, t := range cluster.ProxyTags {
		if t == tag {
			return true
		}
	}
	return false
}
