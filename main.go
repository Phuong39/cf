package main

import (
	"github.com/teamssix/cf/cmd"
	_ "github.com/teamssix/cf/cmd/alibaba"
	_ "github.com/teamssix/cf/cmd/keymanage"
	_ "github.com/teamssix/cf/cmd/tencent"
)

func main() {
	cmd.Execute()
}
