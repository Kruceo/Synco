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
