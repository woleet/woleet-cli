# Install woleet-cli

## Get latest binaries

The latest binaries can be found [here](https://github.com/woleet/woleet-cli/releases)

Just download it, decompress it and execute it (add execution permissions if necessary)

### For Linux or MacOS
You can use this command to install the latest binaries in /usr/local/bin

```bash
CLI_URL=$(curl --silent https://api.github.com/repos/woleet/woleet-cli/releases/latest | grep 'browser_download_url' | grep -ioE "https://.*$(uname -s)_x86_64.tar.gz") && \
sudo curl -L "$CLI_URL" | sudo tar -xz -C /usr/local/bin woleet-cli && \
sudo chmod +x /usr/local/bin/woleet-cli
```
