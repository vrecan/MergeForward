package merge

import (
	"errors"
	"fmt"
	DIFF "github.com/sergi/go-diff/diffmatchpatch"
	"time"
)

var MaxTime = 60 * time.Second
var ERRFailedPatch = errors.New("Patch failed to apply")

//do a merge using a diff patch
func DiffMerge(src string, dst string) (result string, err error) {

	diff := DIFF.New()

	diff.DiffTimeout = MaxTime

	// differences := diff.DiffMain(src, dst, true)
	var oks []bool
	diffs := diff.DiffMain(src, dst, false)
	patches := diff.PatchMake(dst, diffs)
	fmt.Println(patches)
	// fmt.Println(patches)
	result, oks = diff.PatchApply(patches, dst)
	for _, ok := range oks {
		if !ok {
			return result, ERRFailedPatch
		}
	}

	return
}
