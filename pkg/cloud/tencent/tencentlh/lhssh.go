package tencentlh

import (
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/cloud"
	"github.com/teamssix/cf/pkg/util/cmdutil"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var (
	LHSSHCacheFilePath     = cmdutil.ReturnCacheFile("tencent", "LHSSH")
	lhSSHKeyBackupFilePath = cmdutil.ReturnCacheFile("tencent", "sshKeyBackup")
	SSHCacheHeader         = []string{"绑定的实例ID(AssociatedInstanceIds)", "创建时间 (CreatedTime)", "密钥ID (KeyId)", "密钥名称(KeyName)", "区域(Region)"}
	selectRegionList       = []string{"全部地区 (all regions)", "ap-beijing", "ap-chengdu", "ap-guangzhou", "ap-hongkong", "ap-shanghai", "ap-singapore", "eu-frankfurt", "na-siliconvalley", "na-toronto", "ap-mumbai", "eu-moscow", "ap-tokyo", "ap-nanjing"}
)

type sshKeyInstanceBackups struct {
	Original []sshKeyInstanceBackup `json:"original"`
	New      []sshKeyInstanceBackup `json:"new"`
}
type sshKeyInstanceBackup struct {
	AssociatedInstanceIds string `json:"AssociatedInstanceIds"`
	KeyID                 string `json:"KeyId"`
}
type SSHKeyInstances struct {
	Response struct {
		KeyPairSet []struct {
			AssociatedInstanceIds []string    `json:"AssociatedInstanceIds"`
			CreatedTime           time.Time   `json:"CreatedTime"`
			KeyID                 string      `json:"KeyId"`
			KeyName               string      `json:"KeyName"`
			PrivateKey            interface{} `json:"PrivateKey"`
			PublicKey             string      `json:"PublicKey"`
		} `json:"KeyPairSet"`
		RequestID  string `json:"RequestId"`
		TotalCount int    `json:"TotalCount"`
	} `json:"Response"`
}
type createSSHKeyResp struct {
	Response struct {
		KeyPair struct {
			AssociatedInstanceIds []interface{} `json:"AssociatedInstanceIds"`
			CreatedTime           interface{}   `json:"CreatedTime"`
			KeyID                 string        `json:"KeyId"`
			KeyName               string        `json:"KeyName"`
			PrivateKey            string        `json:"PrivateKey"`
			PublicKey             string        `json:"PublicKey"`
		} `json:"KeyPair"`
		RequestID string `json:"RequestId"`
	} `json:"Response"`
}

func GetSSHKeysListInfo(lhFlushCache bool) {

	if lhFlushCache == false {
		cmdutil.PrintSSHCacheFile(LHSSHCacheFilePath, SSHCacheHeader, "tencent", "LHSSH")
	} else {
		GetSSHKeysListInfoRealTime()
	}
}

func GetSSHKeysListInfoRealTime() {
	data := make([][]string, 0)
	var region string
	regionList := make([]string, 0)
	prompt := &survey.Select{
		Message: "选择一个地区 (Choose a region): ",
		Options: selectRegionList,
	}
	survey.AskOne(prompt, &region)
	if region == "全部地区 (all regions)" {
		for index, i := range selectRegionList {
			if index == 0 {
				continue
			}
			regionList = append(regionList, i)
		}
	} else {
		regionList = append(regionList, region)
	}

	//fmt.Println(data)
	for _, i := range regionList {
		client := LHClient(i)
		request := lighthouse.NewDescribeKeyPairsRequest()
		response, err := client.DescribeKeyPairs(request)
		if _, ok := err.(*errors.TencentCloudSDKError); ok {
			fmt.Printf("An API error has returned: %s", err)
			return
		}
		if err != nil {
			panic(err)
		}
		// 输出json格式的字符串回包
		tmp := new(SSHKeyInstances)
		json.Unmarshal([]byte(response.ToJsonString()), tmp)
		for _, item := range tmp.Response.KeyPairSet {
			data = append(data, []string{sliceToString(item.AssociatedInstanceIds), item.CreatedTime.String(), item.KeyID, item.KeyName, i})
			if item.AssociatedInstanceIds != nil {
				for _, item2 := range item.AssociatedInstanceIds {
					writeSSHkeyBackupFile(item2, item.KeyID, false)
				}
			}
		}
		if tmp.Response.KeyPairSet == nil {
			data = append(data, []string{"", "", "", "", i})
		}
	}
	td := cloud.TableData{Header: SSHCacheHeader, Body: data}
	cloud.PrintTable(td, "SSHKEYINFO")
	if region == "全部地区 (all regions)" {
		cmdutil.WriteCacheFile(td, LHSSHCacheFilePath, "all", "all")
	}
}

func GenerateSSHKeyOnInstance(region string, keyName string, specifiedInstanceID string, lhFlushCache bool) {

	var InstancesList []Instances
	if lhFlushCache == false {
		data := cmdutil.ReadCacheFile(LHCacheFilePath, "tencent", "LH")
		for _, i := range data {
			if specifiedInstanceID != "all" {
				if specifiedInstanceID == i[1] {
					obj := Instances{
						InstanceId:       i[1],
						InstanceName:     i[2],
						OSName:           i[3],
						OSType:           i[4],
						Status:           i[5],
						PrivateIpAddress: i[6],
						PublicIpAddress:  i[7],
						RegionId:         i[8],
					}
					InstancesList = append(InstancesList, obj)
				}
			} else {
				obj := Instances{
					InstanceId:       i[1],
					InstanceName:     i[2],
					OSName:           i[3],
					OSType:           i[4],
					Status:           i[5],
					PrivateIpAddress: i[6],
					PublicIpAddress:  i[7],
					RegionId:         i[8],
				}
				InstancesList = append(InstancesList, obj)
			}
		}
	} else {
		InstancesList = ReturnInstancesList(region, false, specifiedInstanceID)
	}
	if len(InstancesList) == 0 {
		if specifiedInstanceID == "all" {
			log.Warnf("未发现实例，可以使用 --flushCache 刷新缓存后再试 (No instances found, You can use the --flushCache command to flush the cache and try again)")
		} else {
			log.Warnf("未找到 %s 实例的相关信息 (No information found about the %s instance)", specifiedInstanceID, specifiedInstanceID)
		}
	} else {
		if specifiedInstanceID == "all" {
			var (
				selectInstanceIDList []string
				selectInstanceID     string
			)
			selectInstanceIDList = append(selectInstanceIDList, "全部实例 (all instances)")
			for _, i := range InstancesList {
				selectInstanceIDList = append(selectInstanceIDList, fmt.Sprintf("%s (%s)", i.InstanceId, i.OSName))
			}
			prompt := &survey.Select{
				Message: "选择一个实例 (Choose a instance): ",
				Options: selectInstanceIDList,
			}
			survey.AskOne(prompt, &selectInstanceID)
			for _, j := range InstancesList {
				if selectInstanceID != "all" {
					if selectInstanceID == fmt.Sprintf("%s (%s)", j.InstanceId, j.OSName) {
						InstancesList = nil
						InstancesList = append(InstancesList, j)
					}
				}
			}
		}

		for _, i := range InstancesList {
			regions := strings.Split(i.RegionId, "-")
			region = regions[0] + "-" + regions[1]
			client := LHClient(region)
			request := lighthouse.NewCreateKeyPairRequest()
			request.KeyName = common.StringPtr(keyName)
			response, err := client.CreateKeyPair(request)
			if _, ok := err.(*errors.TencentCloudSDKError); ok {
				fmt.Printf("An API error has returned: %s", err)
				return
			}
			if err != nil {
				panic(err)
			}
			// 输出json格式的字符串回包
			respSt := new(createSSHKeyResp)
			json.Unmarshal([]byte(response.ToJsonString()), respSt)
			err = bindSSHKey(respSt.Response.KeyPair.KeyID, i.InstanceId, client)
			if err == nil {
				writeSSHkeyBackupFile(i.InstanceId, respSt.Response.KeyPair.KeyID, true)
				ioutil.WriteFile("rsa_pri_"+i.InstanceId+"_"+respSt.Response.KeyPair.KeyID+".txt", []byte(respSt.Response.KeyPair.PrivateKey), 0755)
				fmt.Println("绑定成功，私钥文件已保存在当前目录下")
			}
		}

	}
}

func DeleteSSHKeyOnInstance(region string, specifiedInstanceID string, lhFlushCache bool) {

	if !isFileExists(lhSSHKeyBackupFilePath) {
		fmt.Println("密钥缓存文件不存在，可先执行generate方法生成并绑定密钥对")
		return
	}

	var InstancesList []Instances
	if lhFlushCache == false {
		data := cmdutil.ReadCacheFile(LHCacheFilePath, "tencent", "LH")
		for _, i := range data {
			if specifiedInstanceID != "all" {
				if specifiedInstanceID == i[1] {
					obj := Instances{
						InstanceId:       i[1],
						InstanceName:     i[2],
						OSName:           i[3],
						OSType:           i[4],
						Status:           i[5],
						PrivateIpAddress: i[6],
						PublicIpAddress:  i[7],
						RegionId:         i[8],
					}
					InstancesList = append(InstancesList, obj)
				}
			} else {
				obj := Instances{
					InstanceId:       i[1],
					InstanceName:     i[2],
					OSName:           i[3],
					OSType:           i[4],
					Status:           i[5],
					PrivateIpAddress: i[6],
					PublicIpAddress:  i[7],
					RegionId:         i[8],
				}
				InstancesList = append(InstancesList, obj)
			}
		}
	} else {
		InstancesList = ReturnInstancesList(region, false, specifiedInstanceID)
	}
	if len(InstancesList) == 0 {
		if specifiedInstanceID == "all" {
			log.Warnf("未发现实例，可以使用 --flushCache 刷新缓存后再试 (No instances found, You can use the --flushCache command to flush the cache and try again)")
		} else {
			log.Warnf("未找到 %s 实例的相关信息 (No information found about the %s instance)", specifiedInstanceID, specifiedInstanceID)
		}
	} else {
		if specifiedInstanceID == "all" {
			var (
				selectInstanceIDList []string
				selectInstanceID     string
			)
			selectInstanceIDList = append(selectInstanceIDList, "全部实例 (all instances)")
			for _, i := range InstancesList {
				selectInstanceIDList = append(selectInstanceIDList, fmt.Sprintf("%s (%s)", i.InstanceId, i.OSName))
			}
			prompt := &survey.Select{
				Message: "选择一个实例 (Choose a instance): ",
				Options: selectInstanceIDList,
			}
			survey.AskOne(prompt, &selectInstanceID)
			for _, j := range InstancesList {
				if selectInstanceID != "all" {
					if selectInstanceID == fmt.Sprintf("%s (%s)", j.InstanceId, j.OSName) {
						InstancesList = nil
						InstancesList = append(InstancesList, j)
					}
				}
			}
		}

		for _, i := range InstancesList {
			regions := strings.Split(i.RegionId, "-")
			region = regions[0] + "-" + regions[1]
			client := LHClient(region)

			jsonByte, _ := ioutil.ReadFile(lhSSHKeyBackupFilePath)
			LoadedByte := new(sshKeyInstanceBackups)
			err := json.Unmarshal(jsonByte, LoadedByte)
			if err != nil {
				fmt.Println(err.Error())
				panic(err)
			}
			tmpNewKeyInfos := make([]sshKeyInstanceBackup, 0)

			if LoadedByte.New == nil {
				fmt.Printf("该实例不存在需要删除的密钥")
				return
			}

			for _, instanceGenerateKeyInfo := range LoadedByte.New {
				if instanceGenerateKeyInfo.AssociatedInstanceIds == i.InstanceId {
					unBindFlag := 0
					for true {
						request := lighthouse.NewDisassociateInstancesKeyPairsRequest()
						request.KeyIds = common.StringPtrs([]string{instanceGenerateKeyInfo.KeyID})
						request.InstanceIds = common.StringPtrs([]string{instanceGenerateKeyInfo.AssociatedInstanceIds})
						_, err := client.DisassociateInstancesKeyPairs(request)
						if _, ok := err.(*errors.TencentCloudSDKError); ok {
							if strings.Contains(err.Error(), "Code=UnsupportedOperation.KeyPairNotBoundToInstance") {
								unBindFlag = 1
							}
						}
						if err != nil && !strings.Contains(err.Error(), "Code=UnsupportedOperation.KeyPairNotBoundToInstance") && !strings.Contains(err.Error(), "DisassociateInstancesKeyPairs") {
							fmt.Printf("An API error has returned: %s", err)
							if strings.Contains(err.Error(), "Code=ResourceNotFound.KeyIdNotFound") {
								break
							}
							panic(err)
						}
						if unBindFlag == 1 {
							fmt.Println("解绑完成 执行删除操作")
							break
						}
						fmt.Println("等待解绑完成")
						time.Sleep(5 * time.Second)
					}

					request2 := lighthouse.NewDeleteKeyPairsRequest()
					request2.KeyIds = common.StringPtrs([]string{instanceGenerateKeyInfo.KeyID})
					_, err = client.DeleteKeyPairs(request2)
					//TODO删除时密钥对不存在
					if _, ok := err.(*errors.TencentCloudSDKError); ok {
						if strings.Contains(err.Error(), "Code=ResourceNotFound.KeyIdNotFound") {
							fmt.Println("正在删除中")
						}
					}
					if err != nil && !strings.Contains(err.Error(), "Code=ResourceNotFound.KeyIdNotFound") {
						fmt.Printf("An API error has returned: %s", err)
						panic(err)
					}

				}
				tmpNewKeyInfo := new(sshKeyInstanceBackup)
				tmpNewKeyInfo.KeyID = instanceGenerateKeyInfo.KeyID
				tmpNewKeyInfo.AssociatedInstanceIds = i.InstanceId
				LoadedByte.New = tmpNewKeyInfos
				newByteNeedToWrite, _ := json.Marshal(LoadedByte)
				ioutil.WriteFile(lhSSHKeyBackupFilePath, newByteNeedToWrite, 0777)
			}
			fmt.Println("完成解绑并删除（存在延迟）")
		}

	}
}

func bindSSHKey(keyId string, instanceId string, client *lighthouse.Client) error {

	request := lighthouse.NewAssociateInstancesKeyPairsRequest()
	request.KeyIds = common.StringPtrs([]string{keyId})
	request.InstanceIds = common.StringPtrs([]string{instanceId})
	_, err := client.AssociateInstancesKeyPairs(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return err
	}
	if err != nil {
		panic(err)
	}
	return nil
}

func sliceToString(slice []string) string {
	result := ""
	for index, item := range slice {
		result += item
		if index+1 != len(slice) {
			result += "\n"
		}
	}
	return result
}

func writeSSHkeyBackupFile(instanceId string, keyId string, isNew bool) {

	jsonByte, _ := ioutil.ReadFile(lhSSHKeyBackupFilePath)
	LoadedByte := new(sshKeyInstanceBackups)
	err := json.Unmarshal(jsonByte, LoadedByte)
	if err != nil && err.Error() == "unexpected end of JSON input" {
		b, _ := json.Marshal(LoadedByte)
		ioutil.WriteFile(lhSSHKeyBackupFilePath, b, 0755)
	} else if err != nil {
		fmt.Printf(err.Error())
	}

	if isNew {
		tmp := new(sshKeyInstanceBackup)
		tmp.AssociatedInstanceIds = instanceId
		tmp.KeyID = keyId
		LoadedByte.New = append(LoadedByte.New, *tmp)
		bytes, _ := json.Marshal(LoadedByte)
		ioutil.WriteFile(lhSSHKeyBackupFilePath, bytes, 0755)
		return
	}

	keyFlag := 0
	for _, item := range LoadedByte.New {
		if item.KeyID == keyId {
			keyFlag = 1
		}
	}
	for _, item := range LoadedByte.Original {
		if item.KeyID == keyId {
			keyFlag = 1
		}
	}
	if keyFlag == 0 {
		tmp := new(sshKeyInstanceBackup)
		tmp.AssociatedInstanceIds = instanceId
		tmp.KeyID = keyId
		LoadedByte.Original = append(LoadedByte.Original, *tmp)
		bytes, _ := json.Marshal(LoadedByte)
		ioutil.WriteFile(lhSSHKeyBackupFilePath, bytes, 0755)
	}

}

func isFileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {

		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
