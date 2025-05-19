# Google Workspace - Create and Configure Service Account

- [Google Workspace - Create and Configure Service Account](#google-workspace---create-and-configure-service-account)
  - [Overview](#overview)
  - [High Level](#high-level)
  - [Step-by-Step Guide](#step-by-step-guide)
    - [Create Google Cloud Project](#create-google-cloud-project)
    - [Create Service Account](#create-service-account)
    - [Create and Download the API Key for Service Account](#create-and-download-the-api-key-for-service-account)
    - [Create New Admin Role and Assign to Service Account](#create-new-admin-role-and-assign-to-service-account)
    - [Enable Domain-Wide Delegation for Service Account](#enable-domain-wide-delegation-for-service-account)
    - [Extract Customer ID of Google Workspace](#extract-customer-id-of-google-workspace)


## Overview

This guide explains how to configure a Service Account to fetch Google Groups members inside your Google Workspace.

This Service Account is essential for providing group-based authorization for OPKSSH.

The resulting Service Account will have the following permissions:
- `https://www.googleapis.com/auth/admin.directory.user.readonly`
- `https://www.googleapis.com/auth/admin.directory.group.readonly`

It will also use [domain-wide delegation](https://support.google.com/a/answer/162106?hl=en&src=supportwidget0&authuser=0).

The outcome of this configuration is:
- **Step 6** - Obtain the **Service Account Email**
- **Step 6** - Obtain the **Service Account Unique ID**
- **Step 8** - Obtain the **API Key** for the **Service Account** as a **JSON file**
- **Step 18** - Obtain the **Customer ID**

## High Level

- **Prerequisite:** [Create Google Cloud project "opkssh"](../google-cloud-project/README.md)
- Create a **Service Account** named "opkssh" under the project "opkssh"
- **Create** and **download** the **API Key** for the **Service Account** "opkssh"
- Create a new admin role "opkssh" (with "Users => Read" and "Groups => Read" permissions) and assign it to the Service Account
- Enable domain-wide delegation for the Service Account
- Extract the Customer ID of your Google Workspace

## Step-by-Step Guide

### Create Google Cloud Project

0. Ensure you have the correct project name "opkssh" in the top left corner.

![](./00-prerequisite-google-cloud-project.png)

- If you see a different name, **click on it** and select **opkssh**
- If you **do not see** the name **opkssh** in the project list, you need to complete the guide [Create Google Cloud project](../google-cloud-project/README.md)

### Create Service Account

1. Navigate to [Admin SDK API](https://console.cloud.google.com/apis/library/admin.googleapis.com). You will see the following:

  ![](./01-admin-skd-api.png)

2. Click the "Enable" button and you will see the following:

  ![](./02-admin-skd-api.png)

3. Navigate to [API & Services](https://console.cloud.google.com/apis/credentials) and click **Create Credentials**, then choose **Service Account**.

  ![](./03-create-service-account.png)

4. Enter the name of the Service Account as "opkssh" and click "Done".

  ![](./04-create-service-account.png)

### Create and Download the API Key for Service Account

5. Navigate to [IAM & Admin => Service Account](https://console.cloud.google.com/iam-admin/serviceaccounts) and click on the created Service Account.

  ![](./05-select-service-account.png)

6. On the Service Account screen:
  - Copy the **Service Account Email**; you will need it later.
  - Copy the **Unique ID**; you will need it later.
  - Click the **Keys** tab.

  ![](./06-service-account-details.png)

7. Choose "Add Key" => "Create New Key"

  ![](./07-add-key.png)

8.  Choose **JSON** and click **Create**. Your browser will **download a JSON file with the key** for your Service Account.

  ![](./08-create-api-key.png)

### Create New Admin Role and Assign to Service Account

9. Navigate to [Google Admin => Account => Admin Roles](https://admin.google.com/ac/roles) and click **Create new role**.

  ![](./09-create-new-role.png)

10. Enter the name "opkssh" and click "Continue".

  ![](./10-create-new-role.png)

11. On the **Select Privileges** screen, **scroll down** to the **Admin API privileges** section and select **Groups** => **Read** and **Users** => **Read**, then click **Continue**.

12. On the **Review Privileges** screen, double-check the name, description, and selected privileges, then click "Create Role".

  ![](./12-review-role-priviledges.png)

13. On the screen with the new admin role, click **Assign service accounts**.

  ![](./13-assign-service-account.png)

14. On the **Add service account** screen, enter the Service Account email (**step 6**) and click **Add**.

  ![](./14-add-service-account.png)

15. Click **Assign Role**.

  ![](./15-assign-role.png)

### Enable Domain-Wide Delegation for Service Account

16. Navigate to [Google Admin => Security => API Controls => Domain-wide delegation](https://admin.google.com/ac/owl/domainwidedelegation?hl=en_US) and click **Add New**.

  ![](./16-enable-domain-wide-delegation.png)

17. Configure the scope of the delegation:
  - **Client ID**: **Unique ID** from **step 6**
  - **OAuth scopes** (comma-delimited):
    - `https://www.googleapis.com/auth/admin.directory.user.readonly`
    - `https://www.googleapis.com/auth/admin.directory.group.readonly`

  ![](./17-domain-wide-delegation.png)

### Extract Customer ID of Google Workspace

18. Navigate to [Google Admin => Account Settings](https://admin.google.com/ac/accountsettings/profile?hl=en_US) to find your "Customer ID".

  ![](./18-customer-id.png)
