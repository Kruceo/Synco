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
	"io"
	"os"
	"path"
	"path/filepath"
	"synco/config"
	"synco/utils"
	"time"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch config entries and manage updates continuously",
	Run:   watch,
}

/**TODO error hadnling*/
func watch(cmd *cobra.Command, args []string) {
	for {
		entries := MainConfig.ReadAllEntries()
		for entryIndex, entry := range entries {

			if curBranch, _ := git.ShowCurrentBranch(); curBranch != entry.Branch {

				log.Debugf("Switching branch: %s to %s", curBranch, entry.Branch)

				out, err := git.Checkout(entry.Branch)
				if err != nil {
					log.Error("Initial checkout error: "+curBranch+"=>"+entry.Branch, out)
					continue
				}
			}

			/*If not has log, probably is a new orphan branch*/
			if hasLog, _ := git.HasLogHistory(); hasLog {
				log.Debug("Repository has log history")

				log.Debugf("Fetching branch: %s", entry.Branch)
				out, err := git.Fetch(entry.Branch)

				if err != nil {
					log.Error("Fetch error:", out)
					continue
				}
			}
			log.Debugf("Processing entry: %d %s", entryIndex, entry.Branch)

			processCloud2Local(entryIndex, &entry)
			processLocal2Cloud(entryIndex, &entry)
			time.Sleep(5 * time.Second)
		}
		time.Sleep(30 * time.Second)
	}
}

func init() {
	rootCmd.AddCommand(watchCmd)
}

func processLocal2Cloud(entryIndex int, entry *config.ConfigEntry) {
	/**
	concatenate all files from entry annd Calc the sha256
	compare with entry.lastSha256 and choose if its
	*/
	buffer, err := appendFilesToBuffer(entry.FilePaths)
	if err != nil {
		log.Error("Error while buffering files", err)
		return
	}
	currentSha := sha256.Sum256(buffer)
	currentShaBase64 := base64.RawStdEncoding.EncodeToString(currentSha[:])

	if currentShaBase64 != entry.LastSha256 {

		log.Info("Difference in sha256 detected, updating cloud repository...", currentShaBase64, entry.LastSha256)

		log.Info("Copying files from local paths to blob...")
		for _, filePath := range entry.FilePaths {
			utils.FastCopy(filePath, path.Join(BlobPath, filepath.Base(filePath)))
		}

		log.Info("Adding all changes...")
		if out, err := git.AddAll(); err != nil {
			log.Error("Add error:", err, out)
			return
		}

		log.Info("Committing changes...")
		if out, err := git.Commit("Update"); err != nil {
			log.Error("Commit error:", err, out)
			git.Reset("")
			return
		}

		log.Info("Pushing to remote...")
		if out, err := git.Push("HEAD"); err != nil {
			log.Error("Push error:", err, out)
			return
		}
		nowLocalUnix := time.Now().Unix()

		log.Info("Saving configuration...")
		MainConfig.SetEntry(entryIndex, entry.Branch, entry.FilePaths, uint64(nowLocalUnix), currentShaBase64)

		/*Important: update current entry, the next process in queue will get the old sha256 or date if it not changes*/
		entry.LastSha256 = currentShaBase64
		entry.LocalLastUpdate = uint64(nowLocalUnix)
	}
}

/*
Tests if the remote repository last commit is more recent from last push date (processLocal2Cloud run).
If yes, this will download from remote and restore in local.
*/
func processCloud2Local(entryIndex int, entry *config.ConfigEntry) {
	currentCloudLastUpdate, out, err := git.GetCloudRepoCommitTime(entry.Branch)
	if err != nil {
		log.Error("Show error:\n"+out, err)
		return
	}

	/**TODO optimize file buffers, e.g: reuse buffers from copy and paste*/
	if entry.LocalLastUpdate < uint64(currentCloudLastUpdate.Unix()) {
		log.Info("Remote is more updated...", currentCloudLastUpdate.Local().String())

		log.Info("Executing git pull...")
		git.Pull(entry.Branch)

		log.Info("Copying files from blob to local paths...")
		for _, filePath := range entry.FilePaths {
			utils.FastCopy(path.Join(BlobPath, filepath.Base(filePath)), filePath)
		}
		log.Info("Calculating new sha256...")
		buff, err := appendFilesToBuffer(entry.FilePaths)
		if err != nil {
			log.Error("Error while buffering files", err)
			return
		}
		currentSha := sha256.Sum256(buff)
		currentShaBase64 := base64.RawStdEncoding.EncodeToString(currentSha[:])

		nowLocalUnix := time.Now().Unix()

		log.Info("Saving configuration...")
		MainConfig.SetEntry(entryIndex, entry.Branch, entry.FilePaths, uint64(nowLocalUnix), currentShaBase64)

		/*Important: update current entry, the next process in queue will get the old sha256 or date if it not changes*/
		entry.LastSha256 = currentShaBase64
		entry.LocalLastUpdate = uint64(nowLocalUnix)
	}
}

func appendFilesToBuffer(filePaths []string) ([]byte, error) {
	buffer := []byte{}

	for _, filePath := range filePaths {
		f, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}

		data, err := io.ReadAll(f)
		if err != nil {
			f.Close()
			return nil, err
		}

		f.Close()

		// O buffer recebe o caminho do arquivo e o conteúdo
		buffer = append(buffer, []byte(filePath)...)
		buffer = append(buffer, data...)
	}

	return buffer, nil
}
