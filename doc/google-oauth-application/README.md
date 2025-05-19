# Google Workspace - Create and Configure OAuth Application

- [Google Workspace - Create and Configure OAuth Application](#google-workspace---create-and-configure-oauth-application)
  - [Overview](#overview)
  - [High Level](#high-level)
  - [Step-by-Step Guide](#step-by-step-guide)
    - [Create Google Cloud Project](#create-google-cloud-project)
    - [Configure Branding, Consent Screen, and Audience](#configure-branding-consent-screen-and-audience)
    - [Create OAuth Application and Obtain "Client ID" and "Client Secret"](#create-oauth-application-and-obtain-client-id-and-client-secret)

## Overview

This guide explains how to configure an OAuth application to authenticate employees of your Google Workspace (and optionally, contractors) for SSH access to your servers using [github.com/openpubkey/opkssh](https://github.com/openpubkey/opkssh).

You have two options for configuring the OAuth application:

- Option A - **Internal** - **Most secure**: Only actual users of your organization will be able to authenticate. **External contractors** will **not** be able to authenticate.
- Option B - **External** - **Less secure**: In addition to users of your organization, **external contractors** from a specific Google Group will **also** be able to authenticate.

The outcome of this configuration will be:
- **Client ID** - **see step 11**
- **Client Secret** - **see step 11**

## High Level

- **Prerequisite:** [Create Google Cloud project "opkssh"](../google-cloud-project/README.md)
- Configure **branding**, **consent screen**, and **audience**
- Create the **OAuth application** and obtain the **Client ID** and **Client Secret**

## Step-by-Step Guide

### Create Google Cloud Project

0. Ensure you have selected the correct project name, "opkssh", in the top left corner.

![](./00-prerequisite-google-cloud-project.png)

- If you see a different name, **click on it** and select **opkssh**.
- If you **do not see** the name **opkssh** in the project list, you need to complete the guide [Create Google Cloud project](../google-cloud-project/README.md).

### Configure Branding, Consent Screen, and Audience

**Without these steps,** you **cannot** create any OAuth application.

1. Navigate to [Google Auth Platform](https://console.cloud.google.com/auth/overview). Review the highlighted areas and click **Get Started**.

  ![](./01-branding.png)

2. Fill out the "App Information" section - **App name** and **User support email**. Click **Next**.

  ![](./02-branding-app-information.png)

3. Decide on the **Audience** for the application - it can be **Internal** or **External**. Choose the appropriate option and click **Next**.

  ![](./03-branding-app-audience.png)

  - Option A - **Internal** - **Most secure**: Only actual users of your organization will be able to authenticate. **External contractors** will **not** be able to authenticate.

    ![](./03-option-A-branding-app-audience-internal.png)

  - Option B - **External** - **Less secure**: In addition to users of your organization, **external contractors** from a specific Google Group will **also** be able to authenticate.

    ![](./03-option-B-branding-app-audience-external.png)

4. Enter your **Contact Information** email - Google will use this to notify you about any changes to your project.

  ![](./04-branding-app-contact-information.png)

5. You need to accept the "Google API Services: User Data Policy" agreement.

  ![](./05-branding-app-consent.png)

### Create OAuth Application and Obtain "Client ID" and "Client Secret"

1. Navigate to [API & Services => Credentials](https://console.cloud.google.com/apis/credentials) and click **Create Credentials**.

  ![](./06-api-and-services.png)

2. Choose **Application Type** - select **Web Application**.

  ![](./07-create-oauth-app-choose-type.png)

3. Enter the **Name** of the application.

  ![](./08-create-oauth-app-name.png)

4. Under the **Authorized redirect URIs** section, click **Add URI**.

  ![](./09-create-oauth-app-redirects.png)

5. Enter the following URLs (these are the default URLs for the OPKSSH client application) and click **Create**:

  ```text
  http://localhost:3000/login-callback
  http://localhost:10001/login-callback
  http://localhost:11110/login-callback
  ```

  ![](./10-create-oauth-app-redirects.png)

6. On the final screen, you will find your **Client ID** and **Client Secret**. You can click the highlighted button to **Copy** the field value or **Download** the JSON file containing this information.
    
  ![](./11-create-oauth-app-info.png)
