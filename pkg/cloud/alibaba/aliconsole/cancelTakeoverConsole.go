package aliconsole

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud/alibaba/aliram"
	"github.com/teamssix/cf/pkg/util/database"
	"github.com/teamssix/cf/pkg/util/errutil"
)

func DetachPolicyFromUser() {
	request := ram.CreateDetachPolicyFromUserRequest()
	request.Scheme = "https"
	request.PolicyType = "System"
	request.PolicyName = "AdministratorAccess"
	request.UserName = "crossfire"
	_, err := aliram.RAMClient().DetachPolicyFromUser(request)
	errutil.HandleErrNoExit(err)
	if err == nil {
		log.Debugln("成功移除 crossfire 用户的权限 (Successfully removed the privileges of the crossfire user)")
	}
}

func DeleteUser() {
	request := ram.CreateDeleteUserRequest()
	request.Scheme = "https"
	request.UserName = "crossfire"
	_, err := aliram.RAMClient().DeleteUser(request)
	errutil.HandleErrNoExit(err)
	if err == nil {
		log.Debugln("删除 crossfire 用户成功 (Delete crossfire user successfully)")
	}
}

func CancelTakeoverConsole() {
	DetachPolicyFromUser()
	DeleteUser()
	database.DeleteTakeoverConsoleCache("alibaba")
	log.Infoln("成功删除 crossfire 用户，已取消控制台接管 (Successful deletion of crossfire user, console takeover cancelled)")
}
