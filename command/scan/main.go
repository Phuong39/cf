package main

import (
	"github.com/teamssix/cf/command/scan/cmd"
	_ "github.com/teamssix/cf/command/scan/cmd/alibaba"
	_ "github.com/teamssix/cf/command/scan/cmd/tencent"
)

func main() {
	cmd.Execute()
}
