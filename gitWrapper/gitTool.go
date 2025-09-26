// Copyright 2025 Kruceo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gitwrapper

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

type GitWrapper struct {
	gitRepositoryDir string
}

func NewGitWrapper(repoDir string) GitWrapper {
	log.Debug("Creating GitWrapper for repo at:", repoDir)
	return GitWrapper{
		gitRepositoryDir: repoDir,
	}
}

// RunGitCommand is a helper to run git commands and return output + error
func (git GitWrapper) RunGitCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = git.gitRepositoryDir

	out, err := cmd.CombinedOutput()
	return string(out), err
}

func (git GitWrapper) Fetch(branch string) (string, error) {
	return git.RunGitCommand("fetch", "origin", branch)
}

func (git GitWrapper) GetCloudRepoCommitTime(branch string) (time.Time, string, error) {
	out, err := git.RunGitCommand("show", "-s", "--format=%cI", "origin/"+branch)
	if err != nil {
		return time.Time{}, out, fmt.Errorf("git show failed: %w", err)
	}

	normalizedOutput := strings.TrimSpace(out)
	t, parseErr := time.Parse(time.RFC3339, normalizedOutput)
	if parseErr != nil {
		return time.Time{}, out, fmt.Errorf("failed to parse commit time: %w", parseErr)
	}

	return t, out, nil
}

func (git GitWrapper) Checkout(branch string) (string, error) {
	return git.RunGitCommand("checkout", branch)
}

func (git GitWrapper) CheckoutNewBranch(branch string, orphan bool) (string, error) {
	if orphan {
		return git.RunGitCommand("checkout", "--orphan", branch)
	}
	return git.RunGitCommand("checkout", "-b", branch)
}

func (git GitWrapper) Reset(mode string) (string, error) {
	switch mode {
	case "hard":
		return git.RunGitCommand("reset", "--hard")
	case "soft":
		return git.RunGitCommand("reset", "--soft")
	}
	return git.RunGitCommand("reset")
}

func (git GitWrapper) Clone(url string, to string) (string, error) {
	return git.RunGitCommand("clone", url, to)
}

func (git GitWrapper) Pull(branch string) (string, error) {
	return git.RunGitCommand("pull", "origin", branch)
}

// ---

// AddAll runs `git add .`
func (git GitWrapper) AddAll() (string, error) {
	return git.RunGitCommand("add", ".")
}

// Commit runs `git commit -m <msg>`
func (git GitWrapper) Commit(message string) (string, error) {
	return git.RunGitCommand("commit", "-m", message)
}

// Push runs `git push origin <branch>`
func (git GitWrapper) Push(branch string) (string, error) {
	return git.RunGitCommand("push", "origin", branch)
}

func (git GitWrapper) BranchExistsOnline(branch string) (bool, error) {
	out, err := git.RunGitCommand("ls-remote", "--heads", "origin", branch)
	if err != nil {
		return false, err
	}
	if strings.TrimSpace(out) != "" {
		return true, nil
	}
	return false, nil
}

func (git GitWrapper) HasLogHistory() (bool, error) {
	out, err := git.RunGitCommand("log")
	if err != nil {
		return false, err
	}
	if strings.TrimSpace(out) != "" {
		return true, nil
	}
	return false, nil
}

func (git GitWrapper) ShowCurrentBranch() (string, error) {
	out, err := git.RunGitCommand("branch", "--show-current")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}
