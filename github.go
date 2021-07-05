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
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/go-github/v36/github"
	"golang.org/x/oauth2"
)

type githubEvent struct {
	Owner       string
	Repo        string
	RepoRemote  string
	PRNumber    int
	HeadBranch  string
	HeadCommit  string
	BaseBranch  string
	IsComment   bool
	SenderLogin string
}

// parsePullRequestEvent uses the given environment variables
// to parse and build githubEvent struct.
// Currently, we support 2 kinds of event as below:
// - PullRequestEvent
//   https://pkg.go.dev/github.com/google/go-github/v36/github#PullRequestEvent
// - IssueCommentEvent
//   https://pkg.go.dev/github.com/google/go-github/v36/github#IssueCommentEvent
func parseGitHubEvent() (*githubEvent, error) {
	const (
		pullRequestEventName = "pull_request"
		commentEventName     = "issue_comment"
	)

	eventName := os.Getenv("GITHUB_EVENT_NAME")
	if eventName != pullRequestEventName && eventName != commentEventName {
		return nil, fmt.Errorf("unexpected event %s, only %q and %q event are supported", eventName, pullRequestEventName, commentEventName)
	}

	eventPath := os.Getenv("GITHUB_EVENT_PATH")
	payload, err := ioutil.ReadFile(eventPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read event payload: %v", err)
	}

	event, err := github.ParseWebHook(eventName, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to parse event payload: %v", err)
	}

	switch e := event.(type) {
	case *github.PullRequestEvent:
		return &githubEvent{
			Owner:      e.Repo.Owner.GetName(),
			Repo:       e.Repo.GetName(),
			RepoRemote: e.Repo.GetSSHURL(),
			PRNumber:   e.GetNumber(),
			HeadBranch: e.PullRequest.Head.GetRef(),
			HeadCommit: e.PullRequest.Head.GetSHA(),
			BaseBranch: e.PullRequest.Base.GetRef(),
		}, nil

	case *github.IssueCommentEvent:
		return &githubEvent{
			Owner:      e.Repo.Owner.GetName(),
			Repo:       e.Repo.GetName(),
			RepoRemote: e.Repo.GetSSHURL(),
			PRNumber:   e.Issue.GetNumber(),
			//HeadBranch:  e.Issue.Head.GetRef(),
			//HeadCommit:  e.PullRequest.Head.GetSHA(),
			//BaseBranch:  e.PullRequest.Base.GetRef(),
			IsComment:   true,
			SenderLogin: e.Sender.GetName(),
		}, nil

	default:
		return nil, fmt.Errorf("got an unexpected event type, got: %t", e)
	}
}

func sendComment(ctx context.Context, token string, pr int, body string) (*github.IssueComment, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	// https://pkg.go.dev/github.com/google/go-github/v36/github#IssueComment
	c, _, err := client.Issues.CreateComment(ctx, "owner", "repo", pr, &github.IssueComment{
		Body: &body,
	})
	return c, err
}
