package config

import (
	"strings"

	"github.com/google/go-cmp/cmp"
)

// PathReporter is a simple custom reporter that only records differences
// detected during comparison.
type PathReporter struct {
	path  cmp.Path
	diffs []string
}

// PushStep appends a path to the reporter.
func (r *PathReporter) PushStep(ps cmp.PathStep) {
	r.path = append(r.path, ps)
}

// Report appends only changed paths to the reporter.
func (r *PathReporter) Report(rs cmp.Result) {
	if !rs.Equal() {
		r.diffs = append(r.diffs, r.path.String())
	}
}

// PopStep drops a path from the reporter.
func (r *PathReporter) PopStep() {
	r.path = r.path[:len(r.path)-1]
}

// String returns a string representation of all diffs,
// each one in a single new line.
func (r *PathReporter) String() string {
	return strings.Join(r.diffs, "\n")
}

// Roots returns a slice of all affected path root elements.
func (r *PathReporter) Roots() []string {
	var rr []string
	for _, d := range r.diffs {
		c := strings.Split(d, ".")
		rr = append(rr, c[0])
	}
	return rr
}
