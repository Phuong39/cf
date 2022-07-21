<p align="center">
<img width="500" src="https://cdn.jsdelivr.net/gh/teamssix/BlogImages/imgs/202207022209168.png"><br><br>
<a href="https://github.com/teamssix/cf/stargazers"><img alt="GitHub stars" src="https://img.shields.io/github/stars/teamssix/cf"/></a>
<a href="https://github.com/teamssix/cf/releases"><img alt="GitHub releases" src="https://img.shields.io/github/release/teamssix/cf"/></a>
<a href="https://github.com/teamssix/cf/blob/main/LICENSE"><img alt="License" src="https://img.shields.io/badge/License-Apache%202.0-blue.svg"/></a>
<a href="https://github.com/teamssix/cf/releases"><img alt="Downloads" src="https://img.shields.io/github/downloads/teamssix/cf/total?color=brightgreen"/></a>
<a href="https://goreportcard.com/report/github.com/teamssix/cf"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/teamssix/cf"/></a>
<a href="https://twitter.com/intent/tweet/?text=CF%2C%20an%20amazing%20cloud%20exploitation%20framework%0Ahttps%3A%2F%2Fgithub.com%2Fteamssix%2Fcf%0A%23cloud%20%23security%20%23cloudsecurity%20%23cybersecurtiy"><img alt="tweet" src="https://img.shields.io/twitter/url?url=https://github.com/teamssix/cf" /></a>
<a href="https://twitter.com/teamssix"><img alt="Twitter" src="https://img.shields.io/twitter/follow/teamssix?label=Followers&style=social" /></a>
<a href="https://github.com/teamssix"><img alt="Github" src="https://img.shields.io/github/followers/TeamsSix?style=social" /></a><br></br>
<a href="README.md">中文</a> | English
</p>




---

CF is a cloud exploitation framework, It can facilitate the work of the red team after obtaining access key.

CF download address: [github.com/teamssix/cf/releases](https://github.com/teamssix/cf/releases)

> Currently CF only supports Alibaba cloud, and will continue to update the support for other cloud providers.

Currently CF has these functions as follows: 

* Currently supported features

  - [x] List OSS
  - [x] List ECS
  - [x] Get the STS Token in the instance metadata
  - [x] Batch execution of multiple commands used to prove permission acquisition
  - [x] Get intances reverse shell
  - [x] Support alibaba cloud
  - [x] List RDS
  - [x] Takeover console
  - [x] View permissions for access key
  - [x] Support Tencent Cloud
  - [x] ......
  
* Functions to be implemented in the future
  - [ ] Attack trail removal
  
  - [ ] Automatically detect if the current running environment is an instance, and if so, scan the local instance for credential information
  - [ ] Add the resulting credentials to the CF
  - [ ] Support other cloud providers
  - [ ] ......

## Manual

For detailed manuals, please visit: [wiki.teamssix.com/cf](https://wiki.teamssix.com/cf)

> The manual currently supports Chinese only

[![img](https://cdn.jsdelivr.net/gh/teamssix/BlogImages/imgs/202207112152449.png)](https://wiki.teamssix.com/cf)

## Case

[《我用 CF 打穿了他的云上内网》](https://zone.huoxian.cn/d/1341-cf)

## Easy to start

![](https://cdn.jsdelivr.net/gh/teamssix/BlogImages/imgs/202207180028840.png)

Configuration

```bash
cf configure
```

One step lists the cloud service resources with current access key

```bash
cf alibaba ls
```

ls permissions

```bash
cf alibaba permissions
```

takeover console

```bash
cf alibaba console
```

View help information for ecs exec commands

```bash
cf alibaba ecs exec -h
```

Batch execution of multiple commands used to prove permission acquisition

```
cf alibaba ecs exec -b
```

Get the STS Token in the instance metadata

```bash
cf alibaba ecs exec -m
```

View security group policy

```bash
cf tencent vpc ls
```

If it feels good, maybe you can give me a Star ~

## Contributor

Thank you for your contributions to CF ~

<table>
    <tr>
        <td align="center"><a href="https://github.com/teamssix"><img alt="TeamsSix"
                                src="https://avatars.githubusercontent.com/u/49087564?v=4" style="width: 100px;"/><br />TeamsSix</a></td>
        <td align="center"><a href="https://github.com/Amzza0x00"><img alt="Amzza0x00"
                                src="https://avatars.githubusercontent.com/u/32904523?v=4"  style="width: 100px;" /><br />Amzza0x00</a></td>
        <td align="center"><a href="https://github.com/Esonhugh"><img alt="Esonhugh"
                                src="https://avatars.githubusercontent.com/u/32677240?v=4"  style="width: 100px;" /><br />Esonhugh</a></td>
        <td align="center"><a href="https://github.com/Dawnnnnnn"><img alt="Dawnnnnnn"
                                src="https://avatars.githubusercontent.com/u/24506421?v=4"  style="width: 100px;" /><br />Dawnnnnnn</a></td>
</table>


A note on contributions: [CONTRIBUTING](https://github.com/teamssix/cf/blob/main/CONTRIBUTING.md)

## Warning

* This tool can only be used in legal scenarios and is strictly forbidden to be used in illegal scenarios.
* The risks involved in this tool are the responsibility of the tenant and not the cloud providers.

## More

If you are interested in cloud security, you can see my other project [Awesome Cloud Security](https://github.com/teamssix/awesome-cloud-security) , many cloud security resources are included here.

If these cloud security resources are still not enough for you, check out my [cloud security knowledge base](https://wiki.teamssix.com/)), where I have a lot of notes and articles in the direction of cloud security.

Finally, the following is my personal wechat official accounts, welcome to follow ~

**If you would like to work with me on this, you can join the team by sending your resume to admin@wgpsec.org.**

<div align=center><img width="700" src="https://cdn.jsdelivr.net/gh/teamssix/BlogImages/imgs/202204152148071.png" div align=center/></div>

<div align=center><img src="https://api.star-history.com/svg?repos=teamssix/cf&type=Timeline" div align=center/></div>