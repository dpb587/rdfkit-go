package earltesting

import (
	"fmt"
	"strings"
)

// avoid trying to be a full wrapper around testing.T; only match methods for description-related logging

func (a *Assertion) Error(args ...any) {
	a.t.Helper()
	a.Log(args...)
	a.t.Fail()
}

func (a *Assertion) Errorf(format string, args ...any) {
	a.t.Helper()
	a.Logf(format, args...)
	a.t.Fail()
}

func (a *Assertion) Fatal(args ...any) {
	a.t.Helper()
	a.Log(args...)
	a.t.FailNow()
}

func (a *Assertion) Fatalf(format string, args ...any) {
	a.t.Helper()
	a.Logf(format, args...)
	a.t.FailNow()
}

func (a *Assertion) Log(args ...any) {
	a.t.Helper()

	msg := fmt.Sprintln(args...)
	a.descriptionLog.WriteString(msg)

	a.t.Log(msg)
}

func (a *Assertion) Logf(format string, args ...any) {
	a.t.Helper()

	msg := fmt.Sprintf(format, args...)

	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}

	a.descriptionLog.WriteString(msg)

	a.t.Log(msg)
}
