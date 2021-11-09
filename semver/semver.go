package semver

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// SemVer holds a major, minor and patch (teeny) version number.
type SemVer struct {
	Major int
	Minor int
	Patch int
}

const (
	Patch = iota
	Minor
	Major
)

// New returns a SemVer from a string.
// If the string is not a valid SemVer, it effectively returns "v0.0.0".
func New(s string) *SemVer {
	ver := &SemVer{}
	s = strings.TrimPrefix(s, "v")
	a := strings.Split(s, ".")
	if len(a) != 3 {
		return ver
	}

	ver.Major, _ = strconv.Atoi(a[0])
	ver.Minor, _ = strconv.Atoi(a[1])
	ver.Patch, _ = strconv.Atoi(a[2])
	return ver
}

// Bump the versiopn by one, depending on the part argument.
// Major: increase the first segment.
// Minor: increase the second segment.
// Patch: increase the third (last) segment.
func (svl *SemVer) Bump(e int) {
	switch e {
	case Patch:
		svl.Patch++
	case Minor:
		svl.Minor++
		svl.Patch = 0
	case Major:
		svl.Major++
		svl.Minor = 0
		svl.Patch = 0
	}
}

// String returns the SemVer as a string in the format "vX.Y.Z".
func (ver *SemVer) String() string {
	return fmt.Sprintf("v%d.%d.%d", ver.Major, ver.Minor, ver.Patch)
}

// SemVerList is a list of semantic versions.
type SemVerList []*SemVer

// Len returns the length of the list. Used for sorting.
func (list SemVerList) Len() int {
	return len(list)
}

// Less returns true if the first version is less than the second.
func (svl SemVerList) Less(i, j int) bool {
	if svl[i].Major < svl[j].Major {
		return true
	}

	if svl[i].Major > svl[j].Major {
		return false
	}

	if svl[i].Minor < svl[j].Minor {
		return true
	}

	if svl[i].Minor > svl[j].Minor {
		return false
	}

	return svl[i].Patch < svl[j].Patch
}

// Swap two versions in the list. Used for sorting.
func (svl SemVerList) Swap(i, j int) {
	svl[i], svl[j] = svl[j], svl[i]
}

// Sort the list of versions by major, minor then patch version.
func (svl SemVerList) Sort() {
	sort.Sort(svl)
}

// Last returns the most recent version in the list.
func (svl SemVerList) Last() *SemVer {
	if len(svl) == 0 {
		return nil
	}

	svl.Sort()
	fmt.Printf("Sorted: %+v\n", svl)
	return svl[len(svl)-1]
}
