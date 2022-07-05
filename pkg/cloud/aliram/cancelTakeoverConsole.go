package aliram

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util"
)

func DetachPolicyFromUser() {
	request := ram.CreateDetachPolicyFromUserRequest()
	request.Scheme = "https"
	request.PolicyType = "System"
	request.PolicyName = "AdministratorAccess"
	request.UserName = "teamssix"
	_, err := RAMClient().DetachPolicyFromUser(request)
	util.HandleErrNoExit(err)
	if err == nil {
		log.Debugln("成功移除 teamssix 用户的权限 (Successfully removed the privileges of the teamssix user)")
	}
}

func DeleteUser() {
	request := ram.CreateDeleteUserRequest()
	request.Scheme = "https"
	request.UserName = "teamssix"
	_, err := RAMClient().DeleteUser(request)
	util.HandleErrNoExit(err)
	if err == nil {
		log.Debugln("删除 teamssix 用户成功 (Delete teamssix user successfully)")
	}
}

func CancelTakeoverConsole() {
	DetachPolicyFromUser()
	DeleteUser()
	log.Infoln("成功删除 teamssix 用户，已取消控制台接管 (Successful deletion of teamssix user, console takeover cancelled)")
}
