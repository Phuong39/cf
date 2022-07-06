<p align="center">
<img width="500" src="https://cdn.jsdelivr.net/gh/teamssix/BlogImages/imgs/202207022209168.png"><br><br>
<a href="https://github.com/teamssix/cf/stargazers"><img alt="GitHub stars" src="https://img.shields.io/github/stars/teamssix/cf"/></a>
<a href="https://github.com/teamssix/cf/releases"><img alt="GitHub issues" src="https://img.shields.io/github/release/teamssix/cf"/></a>
<a href="https://github.com/teamssix/cf/blob/main/LICENSE"><img alt="License" src="https://img.shields.io/badge/License-Apache%202.0-blue.svg"/></a>
<a href="https://github.com/teamssix/cf/issues"><img alt="GitHub issues" src="https://img.shields.io/github/issues/teamssix/cf"/></a>
<a href="https://goreportcard.com/report/github.com/teamssix/cf"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/teamssix/cf"/></a>
<a href="https://twitter.com/intent/tweet/?text=CF%2C%20an%20amazing%20cloud%20exploitation%20framework%0Ahttps%3A%2F%2Fgithub.com%2Fteamssix%2Fcf%0A%23cloud%20%23security%20%23cloudsecurity%20%23cybersecurtiy"><img alt="tweet" src="https://img.shields.io/twitter/url?url=https://github.com/teamssix/cf" /></a>
<a href="https://twitter.com/teamssix"><img alt="Twitter" src="https://img.shields.io/twitter/url/https/twitter.com/teamssix.svg?style=social&label=Follow" /></a>
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
  - [x] ......
* Features intended to be supported in the short term
  - [ ] Attack trail removal
  - [ ] View permissions for access key
  - [ ] ......
* Features intended to be supported in the long term

  - [ ] Automatically detect if the current running environment is an instance, and if so, scan the local instance for credential information
  - [ ] Add the resulting credentials to the CF
  - [ ] Support other cloud providers such as Tencent Cloud
  - [ ] ......

## Manual

For detailed manuals, please visit: [wiki.teamssix.com/cf](https://wiki.teamssix.com/cf)

> The manual currently supports Chinese only

[![img](https://cdn.jsdelivr.net/gh/teamssix/BlogImages/imgs/202207052324971.png)](https://wiki.teamssix.com/cf)

## Easy to start

![img](https://cdn.jsdelivr.net/gh/teamssix/BlogImages/imgs/202207052307021.png)

Configuration

```bash
cf configure
```

![img](https://cdn.jsdelivr.net/gh/teamssix/BlogImages/imgs/202207022241064.png)

One step lists the cloud service resources with current access key

```bash
cf ls
```

![img](https://cdn.jsdelivr.net/gh/teamssix/BlogImages/imgs/202207040107386.png)

takeover console

```bash
cf console
```

![img](https://cdn.jsdelivr.net/gh/teamssix/BlogImages/imgs/202207052308228.png)

View help information for ecs exec commands

```bash
cf ecs exec -h
```

![img](https://cdn.jsdelivr.net/gh/teamssix/BlogImages/imgs/202207022215293.png)

Batch execution of multiple commands used to prove permission acquisition

```
cf ecs exec -b
```

![img](https://cdn.jsdelivr.net/gh/teamssix/BlogImages/imgs/202207022241381.png)

Get the STS Token in the instance metadata

```bash
cf ecs exec -m
```

![img](https://cdn.jsdelivr.net/gh/teamssix/BlogImages/imgs/202207022241672.png)

If it feels good, maybe you can give me a Star ~

## Warning

* This tool can only be used in legal scenarios and is strictly forbidden to be used in illegal scenarios.
* The risks involved in this tool are the responsibility of the tenant and not the cloud providers.

## More

If you are interested in cloud security, you can see my other project [Awesome Cloud Security](https://github.com/teamssix/awesome-cloud-security) , many cloud security resources are included here.

If these cloud security resources are still not enough for you, check out my [cloud security knowledge base](https://wiki.teamssix.com/)), where I have a lot of notes and articles in the direction of cloud security.

Finally, the following is my personal wechat official accounts, welcome to follow ~

<div align=center><img width="700" src="https://cdn.jsdelivr.net/gh/teamssix/BlogImages/imgs/202204152148071.png" div align=center/></div>

<div align=center><img src="https://api.star-history.com/svg?repos=teamssix/cf&type=Timeline" div align=center/></div>









