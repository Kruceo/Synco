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
	"path"
	"path/filepath"
	"synco/config"
	"synco/utils"
	"time"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update all config entries",
	Long:  `Synchronize all configured entries by pulling from and pushing to the remote Git repository as needed.`,
	Run:   update,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func update(cmd *cobra.Command, args []string) {
	runUpdate()
}

func runUpdate() {
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

		processCloud2Local(entryIndex, &entry) /**TODO create a way to block running if one of these is realy trigered*/
		processLocal2Cloud(entryIndex, &entry)
		time.Sleep(2 * time.Second)
	}
}

/*
Tests if the sum of local entry files are equal the last stored (LastSha256), if not,
update the remote repository and lastSha256.
*/
func processLocal2Cloud(entryIndex int, entry *config.ConfigEntry) {

	buffer, err := utils.AppendFilesToBuffer(entry.FilePaths)
	if err != nil {
		log.Error("Error while buffering files", err)
		return
	}
	currentSha := sha256.Sum256(buffer)
	currentShaBase64 := base64.RawStdEncoding.EncodeToString(currentSha[:])

	if currentShaBase64 != entry.LastSha256 {

		log.Info("Difference in sha256 detected, updating cloud repository...")
		log.Debug("Old sha256:", entry.LastSha256)
		log.Debug("New sha256:", currentShaBase64)

		log.Info("Copying files from local paths to blob...")
		for _, filePath := range entry.FilePaths {
			utils.FastCopy(filePath, path.Join(config.BlobPath, filepath.Base(filePath)))
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
			utils.FastCopy(path.Join(config.BlobPath, filepath.Base(filePath)), filePath)
		}
		log.Info("Calculating new sha256...")
		buff, err := utils.AppendFilesToBuffer(entry.FilePaths)
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
