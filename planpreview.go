// Copyright 2021 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type PlanPreviewResult struct {
	Applications        []ApplicationResult
	FailureApplications []FailureApplication
	FailurePipeds       []FailurePiped
}

type ApplicationResult struct {
	ApplicationInfo
	SyncStrategy string // QUICK_SYNC, PIPELINE
	PlanSummary  string
	PlanDetails  string
}

type FailurePiped struct {
	PipedInfo
	Reason string
}

type FailureApplication struct {
	ApplicationInfo
	Reason      string
	PlanDetails string
}

type PipedInfo struct {
	PipedID  string
	PipedURL string
}

type ApplicationInfo struct {
	ApplicationID        string
	ApplicationName      string
	ApplicationURL       string
	EnvID                string
	EnvName              string
	EnvURL               string
	ApplicationKind      string // KUBERNETES, TERRAFORM, CLOUDRUN, LAMBDA, ECS
	ApplicationDirectory string
}

func retrievePlanPreview(
	ctx context.Context,
	remoteURL,
	baseBranch,
	headBranch,
	headCommit,
	address,
	apiKey string,
	timeout time.Duration,
) (*PlanPreviewResult, error) {

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, fmt.Errorf("failed to create a temporary directory (%w)", err)
	}
	outPath := filepath.Join(dir, "result.json")

	args := []string{
		"plan-preview",
		"--repo-remote-url", remoteURL,
		"--base-branch", baseBranch,
		"--head-branch", headBranch,
		"--head-commit", headCommit,
		"--address", address,
		"--api-key", apiKey,
		"--timeout", timeout.String(),
		"--out", outPath,
	}
	cmd := exec.CommandContext(ctx, "pipectl", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute pipectl command (%w) (%s)", err, string(out))
	}

	log.Println(string(out))

	data, err := ioutil.ReadFile(outPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read result file (%w)", err)
	}

	var r PlanPreviewResult
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("failed to parse result file (%w)", err)
	}

	return &r, nil
}

func makeCommentBody(event *githubEvent, r *PlanPreviewResult) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("@%s, ", event.SenderLogin))

	if len(r.Applications)+len(r.FailureApplications)+len(r.FailurePipeds) == 0 {
		fmt.Fprintf(&b, "This pull request does not touch any applications\n")
		return b.String()
	}

	b.WriteString(fmt.Sprintf("Here are plan-preview result for commit %s:\n", event.HeadCommit))

	if len(r.Applications) > 0 {
		if len(r.Applications) > 1 {
			fmt.Fprintf(&b, "\nHere are plan-preview for %d applications:\n", len(r.Applications))
		} else {
			fmt.Fprintf(&b, "\nHere are plan-preview for 1 application:\n")
		}
		for i, app := range r.Applications {
			fmt.Fprintf(&b, "\n%d. app: %s, env: %s, kind: %s\n", i+1, app.ApplicationName, app.EnvName, app.ApplicationKind)
			fmt.Fprintf(&b, "  sync strategy: %s\n", app.SyncStrategy)
			fmt.Fprintf(&b, "  summary: %s\n", app.PlanSummary)
			fmt.Fprintf(&b, "  details:\n\n  ---DETAILS_BEGIN---\n%s\n  ---DETAILS_END---\n", app.PlanDetails)
		}
	}

	if len(r.FailureApplications) > 0 {
		if len(r.FailureApplications) > 1 {
			fmt.Fprintf(&b, "\nNOTE: An error occurred while building plan-preview for the following %d applications:\n", len(r.FailureApplications))
		} else {
			fmt.Fprintf(&b, "\nNOTE: An error occurred while building plan-preview for the following application:\n")
		}
		for i, app := range r.FailureApplications {
			fmt.Fprintf(&b, "\n%d. app: %s, env: %s, kind: %s\n", i+1, app.ApplicationName, app.EnvName, app.ApplicationKind)
			fmt.Fprintf(&b, "  reason: %s\n", app.Reason)
			if len(app.PlanDetails) > 0 {
				fmt.Fprintf(&b, "  details:\n\n  ---DETAILS_BEGIN---\n%s\n  ---DETAILS_END---\n", app.PlanDetails)
			}
		}
	}

	if len(r.FailurePipeds) > 0 {
		if len(r.FailurePipeds) > 1 {
			fmt.Fprintf(&b, "\nNOTE: An error occurred while building plan-preview for applications of the following %d Pipeds:\n", len(r.FailurePipeds))
		} else {
			fmt.Fprintf(&b, "\nNOTE: An error occurred while building plan-preview for applications of the following Piped:\n")
		}
		for i, piped := range r.FailurePipeds {
			fmt.Fprintf(&b, "\n%d. piped: %s\n", i+1, piped.PipedID)
			fmt.Fprintf(&b, "  reason: %s\n", piped.Reason)
		}
	}

	return b.String()
}
