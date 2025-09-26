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
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"synco/config"
	"synco/utils"
	"time"

	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Inicia o servidor web",
	Long:  `Este comando inicia o servidor HTTP para a aplicação.`,
	Run:   watch,
}

/**TODO error hadnling*/
func watch(cmd *cobra.Command, args []string) {
	a, b := git.HasLogHistory()
	fmt.Println(a, b)
	for {
		entries := MainConfig.ReadAllEntries()
		for entryIndex, entry := range entries {
			if hasLog, _ := git.HasLogHistory(); hasLog {
				out, err := git.Fetch(entry.Branch)
				if err != nil {
					fmt.Println("Fetch error:", out)
					continue
				}
			}

			if curBranch, _ := git.ShowCurrentBranch(); curBranch != entry.Branch {
				out, err := git.Checkout(entry.Branch)
				if err != nil {
					fmt.Println("Checkout error:", out)
					continue
				}
			}
			processCloud2Local(entryIndex, entry)
			processLocal2Cloud(entryIndex, entry)
		}
		time.Sleep(30 * time.Second)
	}
}

func init() {
	rootCmd.AddCommand(watchCmd)
}

func processLocal2Cloud(entryIndex int, entry config.ConfigEntry) {
	/**
	concatenate all files from entry annd Calc the sha256
	compare with entry.lastSha256 and choose if its
	*/
	buffer := appendFilesToBuffer(entry.FilePaths)

	currentSha := sha256.Sum256(buffer)
	currentShaBase64 := base64.RawStdEncoding.EncodeToString(currentSha[:])

	if currentShaBase64 != entry.LastSha256 {

		fmt.Println("Difference in sha256 detected, updating cloud repository...")

		fmt.Println("Copying files from local paths to blob...")
		for _, filePath := range entry.FilePaths {
			utils.FastCopy(filePath, path.Join(BlobPath, filepath.Base(filePath)))
		}

		fmt.Println("Adding all changes...")
		if out, err := git.AddAll(); err != nil {
			fmt.Println("Add error:", err, out)
			return
		}

		fmt.Println("Committing changes...")
		if out, err := git.Commit("Update"); err != nil {
			fmt.Println("Commit error:", err, out)
			git.Reset("")
			return
		}

		fmt.Println("Pushing to remote...")
		if out, err := git.Push("HEAD"); err != nil {
			fmt.Println("Push error:", err, out)
			return
		}
		nowLocalUnix := time.Now().Unix()

		fmt.Println("Saving configuration...")
		MainConfig.SetEntry(entryIndex, entry.Branch, entry.FilePaths, uint64(nowLocalUnix), currentShaBase64)
	}
}

func processCloud2Local(entryIndex int, entry config.ConfigEntry) {
	/**
	GIT PULL MOMENT
	*/
	currentCloudLastUpdate, _, err := git.GetCloudRepoCommitTime(entry.Branch)
	if err != nil {
		fmt.Println("Show error:", err)
		return
	}

	/**TODO optimize file buffers, e.g: reuse buffers from copy and paste*/
	if entry.LocalLastUpdate < uint64(currentCloudLastUpdate.Unix()) {
		fmt.Println("Remote is more updated...", currentCloudLastUpdate.Local().String())

		fmt.Println("Executing git pull...")
		git.Pull(entry.Branch)

		fmt.Println("Copying files from blob to local paths...")
		for _, filePath := range entry.FilePaths {
			utils.FastCopy(path.Join(BlobPath, filepath.Base(filePath)), filePath)
		}
		fmt.Println("Calculating new sha256...")
		buff := appendFilesToBuffer(entry.FilePaths)

		currentSha := sha256.Sum256(buff)
		currentShaBase64 := base64.RawStdEncoding.EncodeToString(currentSha[:])

		nowLocalUnix := time.Now().Unix()

		fmt.Println("Saving configuration...")
		MainConfig.SetEntry(entryIndex, entry.Branch, entry.FilePaths, uint64(nowLocalUnix), currentShaBase64)

	}
}

// ---

// appendFilesToBuffer reformulada como um método
// Esta função não executa comandos Git, mas foi adaptada para o padrão GitWrapper.
func appendFilesToBuffer(filePaths []string) []byte {
	buffer := []byte{}

	for _, filePath := range filePaths {
		// Usando os.Open em vez de os.OpenFile (apenas leitura)
		f, err := os.Open(filePath)
		if err != nil {
			// Seu código original apenas printava o erro, mas se os arquivos são essenciais,
			// panicar ou retornar o erro é melhor. Mantenho o 'fmt.Println' para seguir
			// de perto seu original, mas idealmente seria 'return nil, err'.
			fmt.Println(err)
			continue // Pula para o próximo arquivo
		}

		data, err := io.ReadAll(f)
		if err != nil {
			fmt.Println(err)
			f.Close()
			continue
		}

		f.Close()

		// O buffer recebe o caminho do arquivo e o conteúdo
		buffer = append(buffer, []byte(filePath)...)
		buffer = append(buffer, data...)
	}

	return buffer
}
