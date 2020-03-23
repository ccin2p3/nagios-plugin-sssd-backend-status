package nagios

import (
	"fmt"
	"io"
	"os"

	nerrors "github.com/ccin2p3/nagios-plugin-sssd-backend-status/nagsssdbackend/errors"
)

type nagiosOutput struct {
	stdoutWriter io.Writer
	stderrWriter io.Writer
	exitFunc     func(int)
}

func (n nagiosOutput) logErrorf(format string, a ...interface{}) {
	fmt.Fprintf(n.stderrWriter, format, a...)
}

func (n nagiosOutput) logInfof(format string, a ...interface{}) {
	fmt.Fprintf(n.stdoutWriter, format, a...)
}

func (n nagiosOutput) ToNagiosOutput(errs []error) {
	if len(errs) == 0 {
		n.logInfof("sssd domains are online\n")
		n.exitFunc(nagiosOkStatus)
	}

	n.logErrorf("some sssd domains are in error\n")
	for _, err := range errs {
		derr := err.(nerrors.DomainStatusError)
		n.logErrorf("%s\n", derr)
	}

	n.exitFunc(nagiosCriticalStatus)
}

func ToNagiosOutput(errs []error) {
	nagOut := nagiosOutput{
		stdoutWriter: os.Stdout,
		stderrWriter: os.Stderr,
		exitFunc:     os.Exit,
	}

	nagOut.ToNagiosOutput(errs)
}
