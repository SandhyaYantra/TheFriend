# TheFriend 🌐🔌

TheFriend is a lightweight, purely CLI-based tool that mounts remote devices as local file systems over the internet. No cables, no same-network restrictions, just seamless data browsing and transferring.

Imagine plugging in a USB flash drive, but the drive is actually a device located miles away.

## ✨ Features (Fitur Utama)
* **True Remote Mounting:** Mount remote directories as local drives using FUSE.
* **Network Agnostic:** Works across different networks (NAT traversal built-in). No need to be on the same WiFi.
* **Lazy Loading:** Browse files instantly. Data is only transferred when you read or write to a specific file, saving bandwidth.
* **End-to-End Encrypted:** All data streams are encrypted using AES-256-GCM. Your files remain yours.
* **Pure CLI:** Zero GUI bloat. Designed for power users and terminal lovers.

## 🏗️ Architecture (Arsitektur)
TheFriend uses a signaling server for peer discovery and establishes an encrypted P2P tunnel between devices. If a direct connection fails due to strict NATs, it safely falls back to a secure relay. The local file system is simulated using FUSE (Filesystem in Userspace).

## 🚀 Installation (Cara Pemasangan)
*(Note: Instructions will be updated once the binary releases are ready)*

Dependencies required:
* `fuse3` (Linux) or `macFUSE` (macOS)
* Go 1.20+ (for building from source)

```bash
git clone [https://github.com/SandhyaYantra/thefriend.git](https://github.com/SandhyaYantra/thefriend.git)
cd thefriend
make build