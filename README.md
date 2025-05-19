# Google Workspace Plugin for OPKSSH

- [Google Workspace Plugin for OPKSSH](#google-workspace-plugin-for-opkssh)
  - [Overview](#overview)
    - [Why do I need a custom OAuth application instead of using the default one?](#why-do-i-need-a-custom-oauth-application-instead-of-using-the-default-one)
    - [Why do I need a custom plugin for OPKSSH?](#why-do-i-need-a-custom-plugin-for-opkssh)
  - [Configuration](#configuration)
    - [Google Workspace (Google Admin)](#google-workspace-google-admin)
    - [OpenSSH Client and Server (End-User and DevOps)](#openssh-client-and-server-end-user-and-devops)
  - [How it works](#how-it-works)

## Overview

Read the article from Cloudflare: [Open-sourcing OpenPubkey SSH (OPKSSH): integrating single sign-on with SSH](https://blog.cloudflare.com/open-sourcing-openpubkey-ssh-opkssh-integrating-single-sign-on-with-ssh/)

This plugin enables [opkssh](https://github.com/openpubkey/opkssh) to authenticate users from Google Workspace and authorize them based on Google Groups.

### Why do I need a custom OAuth application instead of using the default one?

| Aspect                        | Default OPKSSH configuration         | Your custom OAuth App                       |
| ----------------------------- | ------------------------------------ | ------------------------------------------- |
| Control over who can log in   | Open to all Google accounts          | You control domain or allow-list via policy |
| OAuth consent screen branding | Generic (Cloudflare/OpenPubkey)      | Shows your company name                     |
| Visibility into usage/logs    | None                                 | Full via Google Cloud Console               |
| User trust during login       | "Unknown app" warning                | Branded with your organization—more trustworthy |
| Best practice                 | Not suitable for enterprise/prod use | Recommended for production use              |

### Why do I need a custom plugin for OPKSSH?

| Requirement                                    | opkssh | opkssh-plugin-google-workspace (internal audience of OAuth App) | opkssh-plugin-google-workspace (external audience of OAuth App) |
| ---------------------------------------------- | ------ | --------------------------------------------------------------- | --------------------------------------------------------------- |
| Group-based authorization (organization users) | :x:    | :white_check_mark:                                              | :white_check_mark:                                              |
| Group-based authorization (external members)   | :x:    | :x:                                                             | :white_check_mark:                                              |

## Configuration

### Google Workspace (Google Admin)

- [Create a Google Cloud project "opkssh"](./doc/google-cloud-project/README.md) to hold the necessary OAuth application and Service Account within your Google Workspace. Complete this guide to obtain:
    - **Google Cloud Project** to hold the **OAuth application** and **Service Account**

- [Create an OAuth application "opkssh"](./doc/google-oauth-application/README.md) to enable authentication of your employees and contractors via the `opkssh` client. Complete this guide to obtain:
    - **Client ID** of the OAuth application
    - **Client Secret** of the OAuth application

- [Create a Service Account "opkssh"](./doc/google-servive-account/README.md) to allow `opkssh-plugin-google-workspace` to perform group-based authorization of your employees on your servers. Complete this guide to obtain:
    - **Email** of the **Service Account**
    - **JSON file** with the **API Key** of the **Service Account**
    - **Customer ID** of **Google Workspace**

:warning: You need a Service Account only if you want to use the `opkssh-plugin-google-workspace` plugin.

### OpenSSH Client and Server (End-User and DevOps)

- [Configure SSH](./doc/configure-ssh/README.md) on both client and server.

## How it works

The latest [github.com/openpubkey/opkssh](https://github.com/openpubkey/opkssh) provides "policy plugins" functionality.

Please see the [official documentation from opkssh about this](https://github.com/openpubkey/opkssh/blob/main/docs/policyplugins.md).

The "opkssh-plugin-google-workspace" is a "policy plugin". It receives the following environment variables from "opkssh":
- `OPKSSH_PLUGIN_U` - principal, your system user name to authorize
- `OPKSSH_PLUGIN_EMAIL` - user's email from the incoming token
- `OPKSSH_PLUGIN_EMAIL_VERIFIED` - whether the user's email in the incoming token is verified
- `OPKSSH_PLUGIN_AUD` - audience of the user's token

You can configure your policies in the configuration file `/etc/opkssh-plugin-google-workspace/config.yaml`. Example:
```yaml
policy:
  foo:
    users:
      - work@company.name
      - personal@gmail.com
  bar:
    groups:
      - employee-group@company.name
      - contractor-group@company.name
```

To authorize an incoming user, the plugin needs to access the Google Admin API to fetch group members:
```yaml
google:
  oauth:
    client_id: <Client ID from the "Create OAuth application 'opkssh'" guide> 
  service_account:
    email:    <Service Account Email from the "Create Service Account 'opkssh'" guide>
    key_file: <Path to API Key file from the "Create Service Account 'opkssh'" guide>
  workspace:  
    customer_id: <Customer ID from the "Create Service Account 'opkssh'" guide>
```

The plugin saves the content of the group cache to `/var/cache/opkssh-plugin-google-workspace/cache.json`.
Default cache settings:
```yaml
cache:
  path: /var/cache/opkssh-plugin-google-workspace/cache.json
  duration: 15min
```

The plugin writes logs to `/var/log/opkssh-plugin-google-workspace.log`.

A full example config with all settings:
```yaml
cache:
  path: /var/cache/opkssh-plugin-google-workspace/cache.json
  duration: 15min
google:
  oauth:
    client_id: <Client ID from the "Create OAuth application 'opkssh'" guide> 
  service_account:
    email:    <Service Account Email from the "Create Service Account 'opkssh'" guide>
    key_file: <Path to API Key file from the "Create Service Account 'opkssh'" guide>
  workspace:  
    customer_id: <Customer ID from the "Create Service Account 'opkssh'" guide>
policy:
  foo:
    users:
      - work@company.name
      - personal@gmail.com
  bar:
    groups:
      - employee-group@company.name
      - contractor-group@company.name
```
