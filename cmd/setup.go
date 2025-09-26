// cmd/serve.go

package cmd

import (
	"fmt"
	"regexp"
	"synco/utils"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "setup",
	Short: "Inicia o servidor web",
	Long:  `Este comando inicia o servidor HTTP para a aplicação.`,
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
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Type branch").
				Suggestions([]string{"main", "dev"}).
				Placeholder("synco").
				Value(&selectedBranch),
			huh.NewMultiSelect[string]().
				Options(availableFiles...).
				Value(&selectedFiles),
		))
	form.Run()

	entryIndex, entry := MainConfig.AddEntry(selectedBranch, selectedFiles, 0)

	/**TODO error handling*/
	/**TODO verify if the selected filepaths isn't included in other entries*/
	git.Clone(sshURL, BlobPath)

	if exists, _ := git.BranchExistsOnline(selectedBranch); !exists {
		git.CheckoutNewBranch(selectedBranch, true)
		git.Reset("hard")

		processLocal2Cloud(entryIndex, entry)
	} else {
		git.Checkout(selectedBranch)
		git.Reset("")

		processCloud2Local(entryIndex, entry)
	}
}
