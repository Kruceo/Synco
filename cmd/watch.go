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

	fmt.Println("current =", currentShaBase64)
	fmt.Println("last    =", entry.LastSha256)

	if currentShaBase64 != entry.LastSha256 {

		for _, filePath := range entry.FilePaths {
			fmt.Println()
			utils.FastCopy(filePath, path.Join(BlobPath, filepath.Base(filePath)))
		}

		if out, err := git.AddAll(); err != nil {
			fmt.Println("Add error:", err, out)
			return
		}

		if out, err := git.Commit("Update"); err != nil {
			fmt.Println("Commit error:", err, out)
			git.Reset("")
			return
		}

		if out, err := git.Push("HEAD"); err != nil {
			fmt.Println("Push error:", err, out)
			return
		}
		nowLocalUnix := time.Now().Unix()
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

	println("local date =", entry.LocalLastUpdate)
	println("cloud date =", uint64(currentCloudLastUpdate.Unix()))

	/**TODO optimize file buffers, e.g: reuse buffers from copy and paste*/
	if entry.LocalLastUpdate < uint64(currentCloudLastUpdate.Unix()) {
		fmt.Println("Repo is more updated...")
		git.Pull(entry.Branch)
		for _, filePath := range entry.FilePaths {
			fmt.Println()
			utils.FastCopy(path.Join(BlobPath, filepath.Base(filePath)), filePath)
		}
		buff := appendFilesToBuffer(entry.FilePaths)

		currentSha := sha256.Sum256(buff)
		currentShaBase64 := base64.RawStdEncoding.EncodeToString(currentSha[:])

		nowLocalUnix := time.Now().Unix()

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
