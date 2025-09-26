// cmd/root.go

package cmd

import (
	"fmt"
	"os"
	"path"
	"synco/config"
	gitwrapper "synco/gitWrapper"

	"github.com/spf13/cobra"
)

var BlobPath string

var MainConfig config.ConfigWrapper

var git gitwrapper.GitWrapper

var rootCmd = &cobra.Command{
	Use:   "synco",
	Short: "Uma breve descrição do seu aplicativo",
	Long: `Uma descrição mais longa que se estende por várias linhas e detalha
o que o aplicativo faz. Por exemplo, este aplicativo é um
gerenciador de arquivos fictício.`,

	// A função Run é executada quando este comando é chamado.
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute adiciona todos os comandos filhos ao comando raiz e os define
// para execução.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	BlobPath = path.Join(home, ".synco", "blob")

	if _, err := os.Stat(BlobPath); os.IsNotExist(err) {
		err := os.MkdirAll(BlobPath, 0644)
		if err != nil {
			panic(err)
		}
	}

	config, err := config.NewConfigWrapper(path.Join(BlobPath, "..", "config.json")) //  .NewConfigWrapper{ConfigPath:}
	if err != nil {
		panic(err)
	}

	MainConfig = config

	git = gitwrapper.NewGitWrapper(BlobPath)

}
