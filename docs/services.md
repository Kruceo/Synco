# Creating a Service

## Linux (systemd)

**Path:** `/etc/systemd/system/synco.service`  
Assuming the executable is located at: **`/usr/bin/synco`**

```ini
[Unit]
Description=Synco Watch Service
After=network.target

[Service]
# Program + arguments
ExecStart=/usr/bin/synco watch --interval 60

# Automatically restart on failure
Restart=always
RestartSec=5

# User that runs the service (avoid root if possible). Put your user...
User=user

# Log output (collected by journald by default)
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

### Installation steps

1. Copy the binary to `/usr/bin/`:
   ```bash
   sudo cp synco /usr/bin/
   ```

2. Create the working directory:
   ```bash
   sudo mkdir -p /var/lib/synco
   sudo chown user:user /var/lib/synco
   ```

3. Copy the service file:
   ```bash
   sudo cp synco.service /etc/systemd/system/
   ```

4. Reload `systemd` and enable the service:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable synco
   sudo systemctl start synco
   ```

5. Check status and logs:
   ```bash
   systemctl status synco
   journalctl -u synco -f
   ```

---

## Windows (PowerShell)

Assuming the executable is located at: **`C:\synco.exe`**

```powershell
New-Service -Name "Synco" `
  -BinaryPathName '"C:\synco.exe" watch --interval 60' `
  -DisplayName "Synco Watch Service" `
  -Description "Runs Synco in watch mode with 60s interval" `
  -StartupType Automatic
```

### Installation steps

1. Copy the binary `synco.exe` to `C:\` (or a permanent location such as `C:\Program Files\Synco\`).

2. Open PowerShell as **Administrator**.

3. Run the `New-Service` command above.

4. Start the service:
   ```powershell
   Start-Service Synco
   ```

5. Check service status:
   ```powershell
   Get-Service Synco
   ```

6. To remove the service:
   ```powershell
   Stop-Service Synco
   sc.exe delete Synco
   ```

---

## Notes

- **Linux**: Logs are stored in `journalctl` by default. You can configure file logging with  
  `StandardOutput=append:/var/log/synco.log` if needed.
- **Windows**: Logs go to the *Event Viewer* by default. If you want file logging, the application itself must handle it.
- **Security**: Always run services with a dedicated user. Avoid running as `root` or `Administrator` unless absolutely required.
