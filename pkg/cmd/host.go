package cmd

import (
	"net"
	"os"
	"strings"
)

var (
	chostname = ""
	cip       = ""
	cuser     = ""
	cgitUser  = ""

	hostIPList   = make(map[string]string)
	hostHostList = make(map[string]string)
	trustList    = make(map[string]bool)
)

func init() { getGitUser() }

// GetIP returns current IP address or specified IP address
func GetIP(host ...string) string {
	if len(host) > 0 && strings.TrimSpace(host[0]) != "" {
		if ip, ok := hostIPList[host[0]]; ok {
			return ip
		}
		ip := GetCmdStr("host %s|awk '{print $4}'", host[0])
		if strings.Count(ip, ".") == 3 {
			hostIPList[host[0]] = ip
		}
		return ip
	}

	if cip == "" {
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			cip = GetCmdStr("ifconfig| grep -v '127.0.0.1'| grep -v '192.168.'|grep 'inet '| awk '{print $2}'|head -1")
			cip = strings.TrimPrefix(cip, "addr:")
		}
		for _, address := range addrs {
			// 检查ip地址判断是否回环地址
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					cip = ipnet.IP.String()
				}
			}
		}
	}
	if strings.Count(cip, ".") != 3 {
		cip = ""
	}
	return cip
}

// GetHostname returns current hostname or specified hostname by ip/hostname
func GetHostname(host ...string) string {
	if len(host) > 0 && strings.TrimSpace(host[0]) != "" {
		if h, ok := hostHostList[host[0]]; ok {
			return h
		}
		str := GetCmdStr("host %s", host[0])
		h := ""
		if strings.Contains(str, "has address") {
			h = GetCmdStr("host %s|awk '{print $1}'", host[0])
			hostHostList[host[0]] = h
		} else if strings.Contains(str, "domain name pointer") {
			h = strings.TrimSuffix(GetCmdStr("host %s|awk '{print $5}'", host[0]), ".")
			hostHostList[host[0]] = h
		}
		return h
	}
	if chostname == "" {
		h, err := os.Hostname()
		if err != nil {
			chostname = GetCmdStr("hostname")
		} else {
			chostname = h
		}
	}
	return chostname
}

//GetUser returns current user
func GetUser() string {
	if cuser == "" {
		cuser = GetCmdStr("echo $USER")
	}
	if cuser == "" {
		cuser = GetCmdStr("whoami")
	}
	return cuser
}

// GetCurrentUserByGit returns current user by git config
func getGitUser() {
	//if cgitUser == "" {
	if c, s := RunCmd("git config --list|grep name|head -1|awk -F '=' '{print $2}'"); c {
		cgitUser = s
		//} else {
		//	log.Debug("failed to get git user")
		//	cgitUser = "gituserUnknown"
	}
	//}
	//return cgitUser
}

// GetCurrentUserByGit returns current user by git config
func GetCurrentUserByGit() string { return cgitUser }

//CheckTrust returns whether the relationship is trust
func CheckTrust(host string) bool {
	if i, ok := trustList[host]; ok {
		return i
	}
	trustList[host], _ = RunCmd("ssh %s \"hostname\" &>/dev/null", host)
	return trustList[host]
}
