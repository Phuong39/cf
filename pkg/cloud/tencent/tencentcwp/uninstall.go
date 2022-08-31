package tencentcwp

import (
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util/errutil"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	cwp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cwp/v20180228"
)

func UninstallAgent(UUID string) {
	client := CWPClient("")
	request := cwp.NewDeleteMachineRequest()
	request.Uuid = common.StringPtr(UUID)
	_, err := client.DeleteMachine(request)
	errutil.HandleErr(err)
	log.Info("卸载云镜成功 (Uninstall Agent Success)")
}
