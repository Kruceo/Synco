# Synco

Synco is a powerful, file synchronization tool that leverages the power of Git to keep your files and configurations consistent across multiple devices. It automates the process of adding, committing, pushing, and pulling changes, making it ideal for syncing application settings or any files that need to be harmonized between different machines.

It was created to synchronize configuration files for applications that lack a native sync feature, providing a seamless and automatic way to maintain consistency.

# Goals

- Operate without any intermediate server dependencies.
- Avoid opening network ports.
- Do not require a public IP.

## Features

- **Git-Powered:** Uses a Git repository for robust and reliable file synchronization.
- **Automatic:** Watches for file changes and automatically syncs them.
- **Multi-Device:** Keep files synchronized across two or more devices.
- **Easy Setup:** An interactive setup process makes it simple to configure.
- **Flexible:** Sync different sets of files to different branches.

## How It Works

Synco works by using a central Git repository as the source of truth for your files.

1.  **Setup:** You start by configuring Synco with the SSH URL of your Git repository and selecting the local files you want to sync. This creates a "sync entry" that associates your files with a specific branch in the repository.
2.  **Staging:** Synco maintains a local clone of your repository in `~/.synco/blob`. This directory acts as a staging area for file changes.
3.  **Watching:** The `watch` command runs in the background, continuously monitoring your files for changes.
    - If a local file is modified, Synco copies it to the staging area, commits it, and pushes it to the remote repository.
    - If a change is detected in the remote repository, Synco pulls the changes and updates your local files.

This process ensures that your files are always in sync, no matter which machine you're working on.

## Prerequisites

Before you begin, make sure you have the following installed:

- [Go](https://golang.org/doc/install) (version 1.20 or later)
- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

## Installation

1.  Clone the Synco repository:
    ```bash
    git clone https://github.com/your-username/synco.git
    cd synco
    ```

2.  Build the executable:
    ```bash
    go build
    ```
    This will create a `synco` executable in the current directory. You can move this to a directory in your system's `PATH` (e.g., `/usr/local/bin`) to make it accessible from anywhere.

## Usage

Synco is controlled through two main commands: `setup` and `watch`.

### `setup`

The `setup` command is used to configure a new set of files to be synchronized. It's an interactive process that will guide you through the necessary steps.

```bash
./synco setup
```

You will be prompted for the following information:

- **Git Repository URL (SSH):** The SSH URL of the Git repository you want to use for synchronization (e.g., `git@github.com:user/repo.git`).
- **Branch:** The name of the branch you want to use for this sync entry. You can use different branches to sync different sets of files.
- **Files to Sync:** A list of local files that you want to include in this sync entry.

### `watch`

The `watch` command runs in the background and monitors your files for changes, keeping them in sync with the remote repository.

```bash
./synco watch
```

For the best experience, you can set up a systemd service on Linux or a launchd agent on macOS to run this command automatically on startup.

## Configuration

Synco stores its configuration in `~/.synco/config.json`. This file contains the Git repository URL and the details of each sync entry.

Here is an example of what the `config.json` file might look like:

```json
{
  "gitOrigin": "git@github.com:user/repo.git",
  "entries": [
    {
      "branch": "main",
      "filePaths": [
        "/path/to/your/file1.conf",
        "/path/to/your/file2.json"
      ],
      "localLastUpdate": 1678886400,
      "lastSha256": "a1b2c3d4..."
    }
  ]
}
```

- **`gitOrigin`:** The URL of the remote Git repository.
- **`entries`:** An array of sync entries.
  - **`branch`:** The branch associated with this entry.
  - **`filePaths`:** The absolute paths to the files being synced.
  - **`localLastUpdate`:** The timestamp of the last local update.
  - **`lastSha256`:** The SHA256 hash of the files at the last sync, used to detect changes.

## Contributing

Contributions are welcome! If you have a feature request, bug report, or want to contribute to the code, please open an issue or submit a pull request.

## License

This project is licensed under the Apache 2.0 License. See the [LICENSE](LICENSE) file for details.
