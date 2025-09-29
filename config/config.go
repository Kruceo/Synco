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

package config

import (
	"encoding/json"
	"io"
	"os"

	"github.com/charmbracelet/log"
)

type ConfigWrapper struct {
	ConfigPath string
}

func NewConfigWrapper(path string) (ConfigWrapper, error) {
	log.Debug("Creating ConfigWrapper for config at:", path)
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create empty config object
		defaultCfg := ConfigObj{
			Entries:   []ConfigEntry{},
			GitOrigin: "",
		}

		// Serialize to JSON
		data, err := json.MarshalIndent(defaultCfg, "", "  ")
		if err != nil {
			return ConfigWrapper{}, err
		}

		// Write file
		if writeErr := os.WriteFile(path, data, 0644); writeErr != nil {
			return ConfigWrapper{}, writeErr
		}
	}

	return ConfigWrapper{ConfigPath: path}, nil
}

type ConfigObj struct {
	Entries   []ConfigEntry `json:"entries"`
	GitOrigin string        `json:"gitOrigin"`
}

type ConfigEntry struct {
	Branch          string   `json:"branch"`
	FilePaths       []string `json:"filePaths"`
	LocalLastUpdate uint64   `json:"localLastUpdate"`
	LastSha256      string   `json:"lastSha256"`
}

func (c ConfigWrapper) readConfig() ConfigObj {
	configFile, err := os.OpenFile(c.ConfigPath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		panic(err)
	}
	/**TODO error handling*/
	data, _ := io.ReadAll(configFile)
	var obj ConfigObj
	_ = json.Unmarshal(data, &obj)
	return obj
}

func (c ConfigWrapper) GetGitOrigin() string {
	return c.readConfig().GitOrigin
}

func (c ConfigWrapper) SetGitOrigin(sshUrl string) error {
	cc := c.readConfig()
	cc.GitOrigin = sshUrl
	c.writeConfig(cc)
	return nil
}

func (c ConfigWrapper) ReadAllEntries() []ConfigEntry {
	config := c.readConfig()
	return config.Entries
}

func (c ConfigWrapper) writeConfig(config ConfigObj) {
	configFile, err := os.OpenFile(c.ConfigPath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		panic(err)
	}
	configFile.Truncate(0)
	/**TODO error handling*/
	data, _ := json.Marshal(config)
	_, _ = configFile.Write(data)
}

func (c ConfigWrapper) AddEntry(branch string, filesToWatch []string, lastUpdate uint64) (int, ConfigEntry) {
	currentConfig := c.readConfig()
	currEntr := currentConfig.Entries
	toAdd := ConfigEntry{Branch: branch, FilePaths: filesToWatch, LocalLastUpdate: lastUpdate}
	currEntr = append(currEntr, toAdd)
	currentConfig.Entries = currEntr
	c.writeConfig(currentConfig)
	return len(currEntr) - 1, toAdd
}

func (c ConfigWrapper) SetEntry(index int, branch string, filesToWatch []string, lastUpdate uint64, lastSha256 string) {
	currentConfig := c.readConfig()
	currentConfig.Entries[index] = ConfigEntry{Branch: branch, FilePaths: filesToWatch, LocalLastUpdate: lastUpdate, LastSha256: lastSha256}
	c.writeConfig(currentConfig)
}

func (c ConfigWrapper) RemoveEntry(index int) error {
	currentConfig := c.readConfig()
	var newEntries []ConfigEntry
	for i, v := range currentConfig.Entries {
		if i != index {
			newEntries = append(newEntries, v)
		}
	}
	currentConfig.Entries = newEntries
	c.writeConfig(currentConfig)
	return nil
}
