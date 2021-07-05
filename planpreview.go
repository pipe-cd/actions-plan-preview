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

	log.Println("Plan-preview result:")
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
