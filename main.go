package main

import (
	"github.com/teamssix/cf/cmd"
	_ "github.com/teamssix/cf/cmd/alibaba"
	_ "github.com/teamssix/cf/cmd/aws"
	_ "github.com/teamssix/cf/cmd/huawei"
	_ "github.com/teamssix/cf/cmd/tencent"
)

func main() {
	cmd.Execute()
}
