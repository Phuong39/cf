package aliconsole

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud/alibaba/aliram"
	"github.com/teamssix/cf/pkg/util/database"
	"github.com/teamssix/cf/pkg/util/errutil"
	"strings"
)

func DetachPolicyFromUser(userName string) {
	request := ram.CreateDetachPolicyFromUserRequest()
	request.Scheme = "https"
	request.PolicyType = "System"
	request.PolicyName = "AdministratorAccess"
	request.UserName = userName
	_, err := aliram.RAMClient().DetachPolicyFromUser(request)
	errutil.HandleErrNoExit(err)
	if err == nil {
		log.Debugf("成功移除 %s 用户的权限 (Successfully removed the privileges of the %s user)", userName, userName)
	}
}

func DeleteUser(userName string) {
	request := ram.CreateDeleteUserRequest()
	request.Scheme = "https"
	request.UserName = userName
	_, err := aliram.RAMClient().DeleteUser(request)
	errutil.HandleErrNoExit(err)
	if err == nil {
		log.Debugf("删除 %s 用户成功 (Delete %s user successfully)", userName, userName)
	}
}

func CancelTakeoverConsole() {
	TakeoverConsoleCache := database.SelectTakeoverConsoleCache("alibaba")
	if len(TakeoverConsoleCache) == 0 {
		log.Infoln("未接管过控制台，无需取消 (No takeover of the console, no need to cancel)")
	} else {
		userName := strings.Split(TakeoverConsoleCache[0].UserName, "@")[0]
		DetachPolicyFromUser(userName)
		DeleteUser(userName)
		database.DeleteTakeoverConsoleCache("alibaba")
		log.Infof("成功删除 %s 用户，已取消控制台接管 (Successful deletion of %s user, console takeover cancelled)", userName, userName)
	}
}
