# Configure OpenSSH

- [Configure OpenSSH](#configure-openssh)
  - [Overview](#overview)
  - [Guides](#guides)
    - [Install and configure "opkssh" (client and server)](#install-and-configure-opkssh-client-and-server)
    - [Configure "opkssh" on the client side](#configure-opkssh-on-the-client-side)
    - [Configure "opkssh" on the server side](#configure-opkssh-on-the-server-side)
    - [Install and configure "opkssh-plugin-google-workspace"](#install-and-configure-opkssh-plugin-google-workspace)

## Overview

A Google Administrator must complete the following guides:
- [Create Google Cloud project "opkssh"](./../google-cloud-project/README.md)
- [Create OAuth application "opkssh"](./../google-oauth-application/README.md)
- [Create Service Account "opkssh"](./google-servive-account/README.md)

On both the client and server sides, you must complete the guide:
- [Install and configure "opkssh" (client and server)](#install-and-configure-opkssh-client-and-server)

Client side only:
- [Configure "opkssh" on the client side](#configure-opkssh-on-the-client-side)

Server side only:
- If you do not need group-based authorization:
  - [Configure "opkssh" on the server side](#configure-opkssh-on-the-server-side)
- If you need group-based authorization:
  - [Install and configure "opkssh-plugin-google-workspace"](#install-and-configure-opkssh-plugin-google-workspace)
  
## Guides

### Install and configure "opkssh" (client and server)

A Google Administrator must complete the following guides:
- [Create Google Cloud project "opkssh"](./../google-cloud-project/README.md)
- [Create OAuth application "opkssh"](./../google-oauth-application/README.md)

Install "opkssh":
```bash
wget -qO- "https://raw.githubusercontent.com/openpubkey/opkssh/main/scripts/install-linux.sh" | sudo bash
```

Configure "opkssh" providers:
```bash
echo https://accounts.google.com <CLIENT ID from the guide "Create OAuth application 'opkssh'"> 24h | sudo tee /etc/opk/providers
```

### Configure "opkssh" on the client side

A Google Administrator must complete the following guides:
- [Create Google Cloud project "opkssh"](./../google-cloud-project/README.md)
- [Create OAuth application "opkssh"](./../google-oauth-application/README.md)

You must complete the guide:
- [Install and configure "opkssh" (client and server)](#install-and-configure-opkssh-client-and-server)

Install "openssh-client":
```bash
sudo apt install openssh-client
```

Configure "opkssh" in your home directory:
```bash
mkdir -p ~/.opk
cat >~/.opk/config.yml <<EOF
---
default_provider: google-workspace
providers:
  - alias: google-workspace
    issuer: https://accounts.google.com
    client_id:     <Client ID from the guide "Create OAuth application 'opkssh'"> 
    client_secret: <Client Secret from the guide "Create OAuth application 'opkssh'"> 
    scopes: openid email profile
    access_type: offline
    prompt: consent
    redirect_uris:
      - http://localhost:3000/login-callback
      - http://localhost:10001/login-callback
      - http://localhost:11110/login-callback

EOF
```

### Configure "opkssh" on the server side

A Google Administrator must complete the following guides:
- [Create Google Cloud project "opkssh"](./../google-cloud-project/README.md)
- [Create Service Account "opkssh"](./google-servive-account/README.md)

You must complete the guide:
- [Install and configure "opkssh" (client and server)](#install-and-configure-opkssh-client-and-server)

Install the OpenSSH server:
```bash
sudo apt install openssh-server
```

Configure the OpenSSH server to use "opkssh":
```bash
cat >/etc/ssh/sshd_config.d/opkssh.config <<EOF
AuthorizedKeysCommand /usr/bin/opkssh verify %u %k %t
AuthorizedKeysCommandUser opksshuser
EOF
```

### Install and configure "opkssh-plugin-google-workspace"

```bash
# Install the binary
# TODO: acquire binary
sudo mv opkssh-plugin-google-workspace /usr/local/bin/opkssh-plugin-google-workspace
sudo chown root:opksshuser /usr/local/bin/opkssh-plugin-google-workspace
sudo chmod 750 /usr/local/bin/opkssh-plugin-google-workspace

# Connect the plugin to opkssh

sudo cat >/etc/opk/policy.d/google-workspace.yml <<EOF
name: Google Workspace
command: /usr/local/bin/opkssh-plugin-google-workspace -v
EOF

sudo chown root:opksshuser /etc/opk/policy.d/google-workspace.yml

# Configure opkssh-plugin-google-workspace
sudo cat >/etc/opkssh-plugin-google-workspace/config.yaml <<EOF
google:
  oauth:
    client_id: <Client ID from the guide "Create OAuth application 'opkssh'"> 
  service_account:
    email:    <Service Account Email from the guide "Create Service Account 'opkssh'">
    key_file: <Path to API Key file from the guide "Create Service Account 'opkssh'">
  workspace:  
    customer_id: <Customer ID from the guide "Create Service Account 'opkssh'">
policy:
  <linux user name>:
    users:
      - <user's email>
      - <user's email>
    groups:
      - <group's email>
      - <group's email>
EOF

sudo chown root:opksshuser /etc/opkssh-plugin-google-workspace/config.yaml
```
