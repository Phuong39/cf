package tencentconsole

import (
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud/tencent/tencentcam"
	"github.com/teamssix/cf/pkg/util/database"
	"github.com/teamssix/cf/pkg/util/errutil"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

func DetachPolicyFromUser(userName string) {
	UserUin := GetUserUin(userName)
	request := cam.NewDetachUserPolicyRequest()
	request.PolicyId = common.Uint64Ptr(1)
	request.DetachUin = common.Uint64Ptr(UserUin)
	_, err := tencentcam.CAMClient().DetachUserPolicy(request)
	errutil.HandleErr(err)
	if err == nil {
		log.Debugf("成功移除 %s 用户的权限 (Successfully removed the privileges of the %s user)", userName, userName)
	}
}

func DeleteUser(userName string) {
	request := cam.NewDeleteUserRequest()
	request.Name = common.StringPtr(userName)
	request.Force = common.Uint64Ptr(1)
	_, err := tencentcam.CAMClient().DeleteUser(request)
	errutil.HandleErr(err)
	if err == nil {
		log.Debugf("删除 %s 用户成功 (Delete %s user successfully)", userName, userName)
	}
}

func CancelTakeoverConsole() {
	TakeoverConsoleCache := database.SelectTakeoverConsoleCache("tencent")
	if len(TakeoverConsoleCache) == 0 {
		log.Infoln("未接管过控制台，无需取消 (No takeover of the console, no need to cancel)")
	} else {
		userName := TakeoverConsoleCache[0].UserName
		DetachPolicyFromUser(userName)
		DeleteUser(userName)
		database.DeleteTakeoverConsoleCache("tencent")
		log.Infof("成功删除 %s 用户，已取消控制台接管 (Successful deletion of %s user, console takeover cancelled)", userName, userName)
	}
}
