package cmd

func Kill(pids string) string {
	return "kill -9 " + pids
}

func Killall(pname string) string {
	return "killall -9 -p " + pname
}

func GetPstreeStatus(pname string) string {
	return "pstree $USER|grep " + pname
}

func GetPsStatus(pname string) string {
	return "ps aux|grep " + pname + "|grep -v 'supervise'|grep -v 'grep'|grep -v '_control'|grep -v noah|grep -v 'log'|awk '{print $8}'"
}

func GetPID(mname string) string {
	return "ps aux|grep " + mname + "|grep -v 'supervise'|grep -v 'grep'|grep -v '_control'|grep -v noah|awk '{print $2}'"
}
func GetSPID(mname string) string {
	return "ps aux|grep " + mname + "|grep 'supervise'|grep -v 'grep'|grep -v '_control'|grep -v noah|awk '{print $2}'"
}
