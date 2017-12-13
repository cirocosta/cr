package lib

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/pkg/errors"
)

type ActivityType int

const (
	ActivityUnknown ActivityType = iota
	ActivityStarted
	ActivityErrored
	ActivitySuccess
	ActivityAborted
)

type Activity struct {
	Type ActivityType
	Job  *Job
}

type Ui struct {
	writer *tabwriter.Writer
}

func NewUi() (u Ui) {
	u.writer = new(tabwriter.Writer)
	u.writer.Init(os.Stdout, 0, 8, 0, '\t', 0)

	return
}

func (u Ui) WriteActivity(a *Activity) (err error) {
	switch a.Type {
	case ActivityStarted:
		fmt.Fprintf(u.writer, "job=%s\tstatus=%s\tstart=%s\n",
			a.Job.Id,
			a.Type,
			a.Job.StartTime.String())
	case ActivityErrored, ActivitySuccess, ActivityAborted:
		fmt.Fprintf(u.writer, "job=%s\tstatus=%s\tstart=%s\telapsed=%s\n",
			a.Job.Id,
			a.Type,
			a.Job.StartTime,
			a.Job.FinishTime.Sub(a.Job.StartTime).String())
	default:
		err = errors.Errorf(
			"unknown activity type %+v", a)
		return
	}

	u.writer.Flush()
	return
}
