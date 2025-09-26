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

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"synco/utils"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "setup",
	Short: "Adds a new configuration entry to synco",
	Long:  `Configure a new entry to be synchronized by synco.`,
	Run:   setup,
}

func init() {
	// Adiciona o 'serveCmd' como um subcomando de 'rootCmd'.
	rootCmd.AddCommand(serveCmd)
}

func validateSSHURL(url string) error {
	// return nil
	var sshRegEx = regexp.MustCompile(`\w+?@\w+?\.\w+?:.+?\.git`) //"git@github.com:Kruceo/Bound-Unbound.git"
	if sshRegEx.MatchString(url) {
		return nil
	}
	return fmt.Errorf("%s is not a ssh URL", url)
}
func setup(cmd *cobra.Command, args []string) {

	var sshURL string
	var selectedFiles []string
	var selectedBranch string

	if MainConfig.GetGitOrigin() == "" {
		firstConfigForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Type your git repository URL (SSH)").
					Validate(validateSSHURL).
					Value(&sshURL),
			).Title("First configuration"),
		)

		firstConfigForm.Run()

		MainConfig.SetGitOrigin(sshURL)
	}

	availableFiles := utils.GetOptionsFromDir("./")
	branchForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Type branch").
				Suggestions([]string{"main", "dev"}).
				Placeholder("synco").
				Value(&selectedBranch),
		))
	branchForm.Run()

	/**TODO create a 'in remote' file selector???*/

	filesForm := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Options(availableFiles...).
				Value(&selectedFiles),
		),
	)
	filesForm.Run()

	entryIndex, entry := MainConfig.AddEntry(selectedBranch, selectedFiles, 0)

	/**TODO error handling*/
	/**TODO verify if the selected filepaths isn't included in other entries*/
	/**TODO force upload option? (local2cloud process)*/

	cloneIfNotExists(sshURL, BlobPath)

	if exists, _ := git.BranchExistsOnline(selectedBranch); !exists {
		git.CheckoutNewBranch(selectedBranch, true)
		git.Reset("hard")

		processLocal2Cloud(entryIndex, &entry)
	} else {
		git.Checkout(selectedBranch)
		git.Reset("")

		processCloud2Local(entryIndex, &entry)
	}
}

func cloneIfNotExists(sshURL, blobPath string) error {
	gitDir := filepath.Join(blobPath, ".git")

	// Check if .git already exists
	if _, err := os.Stat(gitDir); !os.IsNotExist(err) {
		log.Debug("Repository already exists, skipping clone")
		return nil
	}

	// Run git clone
	cmd := exec.Command("git", "clone", sshURL, blobPath)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	return cmd.Run()
}
