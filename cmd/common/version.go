// Package version executes and returns the version string
// for the currently running process.
package common

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

// The value of these vars are set through linker options.
//var gitCommit = "Local build"
//var buildDate = "Moments ago"
//var gitTag = "Unknown"

var gitCommit = "{STABLE_GIT_COMMIT}"
var buildDate = "{DATE}"
var gitTag = "{STABLE_GIT_TAG}"

// Version returns the version string of this build.
func Version() string {
	if buildDate == "{DATE}" {
		now := time.Now().Format(time.RFC3339)
		buildDate = now
	}
	return fmt.Sprintf("%s. Built at: %s", BuildData(), buildDate)
}

// SemanticVersion returns the Major.Minor.Patch version of this build.
func SemanticVersion() string {
	return gitTag
}

// BuildData returns the git tag and commit of the current build.
func BuildData() string {
	// if doing a local build, these values are not interpolated
	if gitCommit == "{STABLE_GIT_COMMIT}" {
		commit, err := exec.Command("git", "rev-parse", "HEAD").Output()
		if err != nil {
			log.Println(err)
		} else {
			gitCommit = strings.TrimRight(string(commit), "\r\n")
		}
	}
	return fmt.Sprintf("Carrier/%s/%s", gitTag, gitCommit)
}
