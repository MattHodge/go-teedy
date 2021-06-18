# go-teedy

<a href="https://github.com/MattHodge/go-teedy/releases" target="_blank">![GitHub release (latest by date)](https://img.shields.io/github/v/release/MattHodge/go-teedy?label=VERSION&style=for-the-badge)</a> <a href="https://pkg.go.dev/github.com/MattHodge/go-teedy?tab=doc" target="_blank"><img src="https://img.shields.io/badge/Go-Reference-00ADD8?style=for-the-badge&logo=go" alt="go reference" /></a> ![GitHub](https://img.shields.io/github/license/MattHodge/go-teedy?style=for-the-badge) <a href="https://app.codecov.io/gh/MattHodge/go-teedy/branch/main" target="_blank">![Codecov branch](https://img.shields.io/codecov/c/github/MattHodge/go-teedy/main?logo=codecov&style=for-the-badge)</a> <a href="https://www.twitter.com/MattHodge" target="_blank">![Twitter Follow](https://img.shields.io/twitter/follow/MattHodge?label=%40MattHodge&logo=twitter&style=for-the-badge)</a>

This repository contains:

* **go-teedy**, a Go client library for accessing the API of [Teedy](https://github.com/sismics/docs). Read the [package docs](https://pkg.go.dev/github.com/MattHodge/go-teedy) for library usage.

* **teedy-cli**, a command line tool to backup, restore and import Evernote [enex files](https://evernote.com/blog/how-evernotes-xml-export-format-works/) into [Teedy](https://github.com/sismics/docs).

## ⭐️ teedy-cli

### ⚡️ Installation

* Download `teedy-cli` for your system from the [releases page](https://github.com/MattHodge/go-teedy/releases).
* Extract the binary the location of your choosing.

### ⚙️ Commands

#### `help`

Get CLI help.

```bash
teedy-cli --help
```

#### `backup`

Backup a Teedy instance.

```bash
# Provide username and password via environment variables
export TEEDY_USERNAME=user
export TEEDY_PASSWORD=password

teedy-cli backup --url http://source.teedy.local --destinationpath ./backup
```

#### `restore`

Restore a Teedy instance from a backup.

```bash
# Provide username and password via environment variables
export TEEDY_USERNAME=user
export TEEDY_PASSWORD=password

teedy-cli restore --url http://destination.teedy.local:8080 --sourcepath ./backup
```

#### `deletedocsfortag`

Deletes all documents with a specific Tag ID.

```bash
# Provide username and password via environment variables
export TEEDY_USERNAME=user
export TEEDY_PASSWORD=password

teedy-cli deletedocsfortag --url http://teedy.local:8080 --tagid f3472d4e-ed47-414c-ad7a-be65ab54d107
```

#### `evernote`

Import an Evernote `.enex` file into Teedy.

> 🔔 [Export an Evernote notebook](https://help.evernote.com/hc/en-us/articles/209005557-Export-notes-and-notebooks) to an `.enex` file on disk first

```bash
# Provide username and password via environment variables
export TEEDY_USERNAME=user
export TEEDY_PASSWORD=password

teedy-cli evernote --url http://localhost:8080 --source-enex source.enex
```
