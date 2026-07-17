# tfcred — Terraform Credentials Helper with Secure Context Management

`tfcred` is a secure Terraform credentials helper designed for Windows DevOps engineers. It enables seamless management of multiple Terraform Cloud, Terraform Enterprise, and Terraform Cloud Europe credentials using named, directory-scoped contexts.

Rather than relying on a single global token or storing credentials in plaintext configuration files, `tfcred` implements Terraform's native [Credentials Helper Protocol](https://developer.hashicorp.com/terraform/internals/credentials-helpers). Terraform invokes `tfcred` automatically whenever it requires credentials, allowing the appropriate token to be resolved on demand based on the current working directory.

This enables multiple organizations, teams, and personal accounts to coexist on the same machine without manually editing configuration files or swapping environment variables.

---

# 🛡️ Security & Storage Architecture

## No Plaintext Tokens on Disk

`tfcred` stores only configuration metadata on disk, including:

- Context name
- Organization identifier
- Token type
- Target Terraform domain

No API tokens are ever written to disk.

Configuration is stored at:

```text
%LOCALAPPDATA%\tfcred\contexts.json
```

---

## Windows User Profile Isolation

All configuration is stored within the current user's profile and protected by standard Windows NTFS permissions.

---

## Windows Credential Manager

Terraform API tokens are stored exclusively inside **Windows Credential Manager**.

Secrets are protected using the Windows **Data Protection API (DPAPI)**, making them:

- Encrypted at rest
- Bound to the current Windows user account
- Unavailable to other users
- Never stored in plaintext
- Not automatically portable to another machine

---

# 🚀 Installation

## Recommended (WinGet)

Install `tfcred` without administrator privileges.

```powershell
winget install amiasea.tfcred
```

---

## Build From Source

## 🛠️ Building From Source

If you want to compile and install the project locally on a fresh development machine:

1. Compile the unified single-binary structure using **GoReleaser**:
   ```powershell
   goreleaser build --snapshot --clean
   ```
   *This handles compiling the Windows executable and packages it inside the local `.\dist\` directory.*

2. Run the repository installation script to wire up your local environment tracking:
   ```powershell
   .\scripts\install.ps1
   ```
   *This surgically populates the local `%APPDATA%\terraform.d\plugins` folder, registers your global `tfcred` terminal execution alias, and initializes your `terraform.tfrc` control pointers so you can begin testing immediately.*

---

# What the Installer Configures

The installer automatically:

- Installs the `terraform-credentials-tfcred.exe` executable.
- Registers a Windows **App Paths** execution alias so `tfcred` can be executed without modifying `%PATH%`.
- Creates a dedicated, highly visible Terraform CLI configuration file located at:
  ```text
  %USERPROFILE%\terraform.tfrc
  ```

- Configures that file to use the `tfcred` credentials helper.
- Sets the user-level `TF_CLI_CONFIG_FILE` environment variable to point to `terraform.tfrc`.

Because Terraform always honors `TF_CLI_CONFIG_FILE` when it is set, Terraform ignores the default configuration locations (such as `terraform.rc`) as well as the legacy `credentials.tfrc.json` file. All Terraform CLI configuration is therefore centralized into the managed `terraform.tfrc` file.

---

# 💻 CLI Reference

`tfcred` is implemented as a single executable.

When Terraform invokes it through the Credentials Helper Protocol (`get`), it emits the required JSON response.

When executed directly by a user, it provides an interactive command-line interface.

## ⚙️ Configuration Commands

<table>
<thead>
<tr>
<th style="white-space: nowrap;">Command</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td style="white-space: nowrap;"><b><code>tfcred init</code></b><br><sub><code>--domain &lt;domain&gt;</code></sub></td>
<td>Initialize local configuration storage and configure the default Terraform domain.</td>
</tr>
<tr>
<td style="white-space: nowrap;"><b><code>tfcred config</code></b></td>
<td>Display the current configuration.</td>
</tr>
<tr>
<td style="white-space: nowrap;"><b><code>tfcred config</code></b><br><sub><code>--default-domain &lt;domain&gt;</code></sub></td>
<td>Change the default Terraform domain.</td>
</tr>
<tr>
<td style="white-space: nowrap;"><b><code>tfcred add</code></b><br><sub><code>--context &lt;name&gt;</code> <code>--org &lt;org&gt;</code> <code>--token-type &lt;type&gt;</code> <code>--token &lt;token&gt;</code><br><code>[--domain &lt;domain&gt;]</code> <code>[--switch]</code> <code>[--force]</code></sub></td>
<td>Create a context and securely store its token in Windows Credential Manager.</td>
</tr>
<tr>
<td style="white-space: nowrap;"><b><code>tfcred list</code></b></td>
<td>List all configured contexts and their associated Terraform metadata.</td>
</tr>
<tr>
<td style="white-space: nowrap;"><b><code>tfcred switch &lt;context&gt;</code></b></td>
<td>Bind a context to the current working directory.</td>
</tr>
<tr>
<td style="white-space: nowrap;"><b><code>tfcred remove &lt;context&gt;</code></b></td>
<td>Remove a context, its directory bindings, and its stored credential.</td>
</tr>
<tr>
<td style="white-space: nowrap;"><b><code>tfcred purge</code></b><br><sub><code>--force</code></sub></td>
<td>Remove all contexts, directory bindings, and stored credentials.</td>
</tr>
</tbody>
</table>

---

## 🔍 Inspection & Diagnostics

| Command | Description |
|:---------|:------------|
| **`tfcred current`** | Display the context currently bound to the working directory. |
| **`tfcred status`** | Display the active context metadata resolved from the current directory. |
| **`tfcred context`** <br><sub><span style="color:#4EA3FF">`--json` `--all`</span></sub> | Display detailed context diagnostics, including resolved vault key information. |

---

## 🚀 Initialize

```powershell
tfcred init --domain app.terraform.io (--domain supplies a default)
```

---

# 🗂️ Context Management

## Create Contexts Across Different Directories

`tfcred add` creates a named credential context and stores its token securely in Windows Credential Manager.

A context represents a Terraform credential identity, including:

- Context name
- Organization
- Token type
- Terraform domain

The context metadata is stored in:

```text
%LOCALAPPDATA%\tfcred\contexts.json
```

Only metadata is stored in this file. The actual Terraform token is stored separately in Windows Credential Manager.

The directory where tfcred add --switch is executed is important because the switch operation updates the directory-to-context mapping used later during Terraform credential resolution.

```powershell
cd "C:\Work\infrastructure-live\production"

tfcred add `
  --context platform `
  --org "acme-corp" `
  --token-type team `
  --token "glt-XXX..." `
  --switch
```

This creates the context entry:

```text
platform
 ├── Organization: acme-corp
 ├── Token Type: team
 └── Domain: app.terraform.io
 ```

The --switch flag sets the current directory's active context:

```text
C:\Work\infrastructure-live\production -> platform
```

## Same Directory — Switch Production to Another Context

Multiple context entries can exist globally. Switching changes which existing context is active for the current directory.

```powershell
cd "C:\Work\infrastructure-live\production"

tfcred add `
  --context personal `
  --org "test-sandbox" `
  --token-type user `
  --token "glt-YYY..." `
  --switch
```

### The new context is created, and the existing directory mapping is updated:

```text
C:\Work\infrastructure-live\production -> personal
```

The previous platform context still exists, but it is no longer the active context for this directory.

The directory mapping acts as a pointer to the selected context:

```text
C:\Work\infrastructure-live\production -> personal
```

## Different Directory — Independent Context Assignment

A different working directory has its own independent context mapping.

```powershell
cd "C:\Work\infrastructure-live\staging"

tfcred add `
  --context staging `
  --org "acme-corp" `
  --token-type org `
  --token "glt-ZZZ..." `
  --switch
```

This creates another context entry and maps the staging directory:

```text
C:\Work\infrastructure-live\staging -> staging
```

The final directory mappings are independent:

```text
C:\Work\infrastructure-live\production -> personal
C:\Work\infrastructure-live\staging    -> staging
```

---

How Terraform Credential Resolution Works

When Terraform requires credentials, it invokes the Terraform Credentials Helper protocol:

```powershell
tfcred get <domain>
```

### 🔄 Credential Resolution Pipeline

`tfcred` resolves credentials dynamically using your current terminal location through the following sequential steps:

1. **Terraform CLI Interception:** Terraform natively executes `tfcred get <domain>` in the background during initialization or planning.
2. **Directory Detection:** `tfcred` dynamically reads the user's active current working directory (`cwd`) context.
3. **Workspace Context Matching:** The folder path is parsed and cross-referenced against your stored directory mappings inside `%LOCALAPPDATA%\tfcred\contexts.json`.
4. **Context Key Retrieval:** The specific named active context key (e.g., `platform` or `personal`) bound to that directory is safely retrieved.
5. **Metadata Verification:** The context's configuration metadata registry is loaded to verify organization parameters and target domain scopes.
6. **OS Vault Extraction:** The matching, hardware-encrypted token is securely extracted from the Windows Credential Manager API via native Win32 runtime calls.
7. **Spec-Compliant Delivery:** The token payload is formatted into a pure JSON stream and output to `stdout` for Terraform to cleanly digest.

Example:

```text
Current working directory:

C:\Work\infrastructure-live\production

        |
        v

Directory mapping:

C:\Work\infrastructure-live\production -> personal

        |
        v

Context metadata:

personal
 ├── Organization: test-sandbox
 ├── Token Type: user
 └── Domain: app.terraform.io

        |
        v

Windows Credential Manager

        |
        v

Terraform credential response
```

## Directory Resolution Behavior

Directory matching is exact.

Example mappings:

```text
C:\Work\infrastructure-live\production -> personal
C:\Work\infrastructure-live\staging    -> staging
```

Terraform executed from:

```text
C:\Work\infrastructure-live\production
```

resolves:

```text
personal
```

However, Terraform executed from:

```text
C:\Work\infrastructure-live\production\module-a
```

does <b>not</b> automatically inherit the parent directory context.

The directory must have its own mapping if it should resolve a context.

## Important Context Switching Behavior

Because Terraform credential resolution is based on the current working directory, always verify your location before running:

```powershell
tfcred switch <context>
```

tfcred switch changes the active context pointer for the directory you are currently in. It does not globally change the Terraform credential.

Example:

```powershell
cd "C:\Work\infrastructure-live\production"

tfcred switch platform
```

Updates:

```text
C:\Work\infrastructure-live\production -> platform
```

It does not affect:

```text
C:\Work\infrastructure-live
```

or:

```text
C:\Work\infrastructure-live\production\module-a
```

unless those directories have their own mappings.

## Context Switching Summary

The effective Terraform credential is determined by:

```text
Current Working Directory
          |
          v
Directory Context Mapping
          |
          v
Context Metadata
          |
          v
Windows Credential Manager Token
          |
          v
Terraform Credential Helper Response
```

The directory context mapping is the critical part of the lookup chain.

Always create and switch contexts from the directory where Terraform operations will actually execute.

---

## Use Terraform Normally

```powershell
terraform init
terraform plan
```

Terraform automatically invokes `tfcred` whenever credentials are required.

Changing directories and switching contexts causes Terraform to retrieve the appropriate credential by current directory lookup and requested domain automatically, without modifying Terraform configuration files or environment variables.

---

# 🧼 Uninstallation

Remove `tfcred` using WinGet:

```powershell
winget uninstall amiasea.tfcred
```

---

# What the Uninstaller Removes

The uninstaller:

- Removes the `tfcred` executable.
- Deletes all `tfcred` secrets from Windows Credential Manager.
- Removes the Windows **App Paths** registration.
- Removes the user-level `TF_CLI_CONFIG_FILE` environment variable.

Terraform simply resumes using its standard configuration discovery behavior once `TF_CLI_CONFIG_FILE` has been removed.