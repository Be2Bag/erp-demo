package helpers

import (
	"strings"
	"time"

	"github.com/Be2Bag/erp-demo/models"
)

func InSet(v string, set ...string) bool { // เช็คว่า v อยู่ในชุดค่าหรือไม่
	for _, s := range set {
		if v == s {
			return true
		}
	}
	return false
}
func MaxInt(a, b int) int { // คืนค่าสูงสุดระหว่าง a, b
	if a > b {
		return a
	}
	return b
}
func DeriveTaskStatusFromSteps(steps []models.TaskWorkflowStep) string {
	if len(steps) == 0 {
		return "todo"
	}

	allDone := true
	anyInProg := false

	for _, st := range steps {
		switch st.Status {
		case "done", "skip":
			// ถือว่า step เสร็จแล้ว
		default:
			allDone = false
		}

		if st.Status == "in_progress" {
			anyInProg = true
		}
	}

	switch {
	case allDone:
		return "done"
	case anyInProg:
		return "in_progress"
	default:
		return "todo"
	}
}

func DateToISO(s string) (time.Time, error) {
	if strings.TrimSpace(s) == "" {
		return time.Time{}, nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
