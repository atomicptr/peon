package usererr

import (
	"fmt"
	"os"

	"atomicptr.dev/peon/pkg/cli/commands/hook"
	"atomicptr.dev/peon/pkg/cli/common/errormsg"
)

func ShellIntegrationNotInstalled() {
	fmt.Fprintln(
		os.Stderr,
		errormsg.Render(
			"Shell Integration Not Installed",
			fmt.Sprintf("Please add `%s` to your `%s` config file.", hook.EvalCommand(), hook.ConfigFile()),
		),
	)
}
