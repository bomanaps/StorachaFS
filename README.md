# StorachaFS

**StorachaFS** is a Go-based FUSE filesystem that mounts **Storacha** spaces as POSIX-like directories, allowing seamless read/write access to files stored on the decentralized Storacha network (IPFS + Filecoin). It enables tools like `rsync`, `cp`, or `cat` to work naturally with Storacha storage.

## Features

* Mount a Storacha space as a local folder
* Read files directly from Storacha using the Go SDK (**Guppy**)
* Upload files to Storacha transparently via filesystem write operations
* Supports incremental sync with metadata tracking
* Git-like CLI interface for easy sync commands

## Installation

### Prerequisites

* Go 1.21+
* FUSE library installed (`libfuse` on Linux/macOS)
* Access to a Storacha account/space

## Usage

### Mount a Storacha Space

```bash
./storachafs mount /mnt/storacha --space myspace
```

### Read Files

```bash
cat /mnt/storacha/hello.txt
```

### Upload Files

```bash
cp ./localfile.txt /mnt/storacha/
```

### Rsync Integration

```bash
rsync -av ./local-dir/ /mnt/storacha/
```

## CLI Commands

| Command             | Description                                 |
| ------------------- | ------------------------------------------- |
| `storachafs mount`  | Mount a Storacha space at a local directory |
| `storachafs status` | Show local vs remote changes                |
| `storachafs sync`   | Sync local files to Storacha                |
| `storachafs pull`   | Fetch new or updated files from Storacha    |

## Architecture

```
User CLI / rsync
        │
        ▼
     FUSE Mount (StorachaFS)
        │
        ▼
     Guppy (Go SDK)
        │
        ▼
Storacha Network (IPFS + Filecoin)
```

* **FUSE layer**: Exposes POSIX filesystem interface
* **Guppy**: Handles authentication (UCAN), file uploads/downloads, and space management
* **Storacha network**: Stores content in decentralized storage with IPFS hot storage and Filecoin cold storage

## Development Roadmap

### Phase 1: Core FUSE Implementation
- [ ] Basic FUSE filesystem structure
- [ ] Storacha space mounting functionality
- [ ] File read operations from Storacha
- [ ] Basic file write operations to Storacha
- [ ] Directory listing and navigation

### Phase 2: Enhanced Operations
- [ ] File metadata handling (timestamps, permissions)
- [ ] Incremental sync capabilities
- [ ] Local caching for performance
- [ ] Error handling and recovery
- [ ] Basic CLI commands (`status` etc.)

### Phase 3: Advanced Features
- [ ] Git-like sync commands (`sync`, `pull`)
- [ ] Background synchronization
- [ ] Performance optimizations
- [ ] Comprehensive logging and monitoring

### Phase 4: Production Ready
- [ ] Robust error handling and edge cases
- [ ] Configuration file support
- [ ] Authentication improvements
- [ ] Documentation and examples
- [ ] Testing and benchmarking

## Current Status

**Early Development** - StorachaFS is currently in the initial development phase. Core FUSE implementation and Storacha integration are in progress.
