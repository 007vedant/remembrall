# Remembrall Password Manager

A secure CLI password manager written in Go that helps you store and retrieve passwords for various applications and websites.

## 🔒 Security Features

- **AES-256-GCM encryption** with PBKDF2 key derivation
- **Master password authentication** for all operations
- **Hidden password input** (no shoulder surfing)
- **Copy to clipboard** for retrieved passwords
- **SQLite database** stored securely in your home directory
- **Fuzzy search** with intelligent matching

## 📦 Installation

### One-liner Installation (Recommended)

```bash
curl -sSL https://raw.githubusercontent.com/007vedant/remembrall/main/install.sh | bash
```

### Manual Installation

1. Clone the repository:
```bash
git clone https://github.com/007vedant/remembrall.git
cd remembrall
```

2. Run the installation script:
```bash
./install.sh
```

### Requirements

- Go 1.19 or later
- Unix-like system (Linux, macOS, WSL)

## 🚀 Quick Start

1. **Save your first password:**
```bash
remembrall save github
```

2. **List all stored applications:**
```bash
remembrall list
```

3. **Retrieve a password:**
```bash
remembrall get github
```

4. **Search with fuzzy matching:**
```bash
remembrall search git
```

## 💻 Usage

### Core Commands

| Command | Description | Example |
|---------|-------------|---------|
| `save <app-name>` | Save a password for an application | `remembrall save gmail` |
| `get <app-name>` | Retrieve a password (copies to clipboard) | `remembrall get gmail` |
| `update <app-name>` | Update an existing password | `remembrall update gmail` |
| `list` | List all stored applications | `remembrall list` |
| `search <query>` | Search applications with fuzzy matching | `remembrall search gmai` |

### Getting Help

```bash
remembrall --help                # Show all commands
remembrall save --help           # Help for specific command
```

## 🔍 Fuzzy Search

Remembrall includes intelligent fuzzy search that works with:

- **Partial names**: `git` finds `github`
- **Typos**: `gmai` finds `gmail`
- **Subsequences**: `gh` finds `github`
- **Word boundaries**: `work` finds `work-laptop`

## 🛡️ Security Design

### Encryption
- **Algorithm**: AES-256-GCM (authenticated encryption)
- **Key Derivation**: PBKDF2 with 100,000 iterations
- **Salt**: Random 16-byte salt per password
- **Nonce**: Random 12-byte nonce per encryption

### Master Password
- Never stored on disk
- Used to derive encryption keys
- Verified through encrypted test string
- Required for all operations

### Database
- **Location**: `~/.remembrall.db`
- **Content**: Only encrypted passwords
- **Permissions**: User-readable only

## 🗑️ Uninstallation

```bash
curl -sSL https://raw.githubusercontent.com/007vedant/remembrall/main/uninstall.sh | bash
```

Or if you have the repository:
```bash
./uninstall.sh
```

This will remove:
- The Remembrall binary
- All stored passwords
- Master password verification
- PATH configuration

## 🔧 Development

### Building from Source

```bash
git clone https://github.com/007vedant/remembrall.git
cd remembrall
go build -o remembrall cmd/remembrall/main.go
```

### Project Structure

```
remembrall/
├── cmd/remembrall/          # Main application entry point
├── internal/
│   ├── auth/               # Authentication and input handling
│   ├── crypto/             # Encryption/decryption
│   ├── db/                 # Database operations
│   ├── search/             # Fuzzy search algorithms
│   └── ui/                 # CLI commands and interface
├── pkg/models/             # Data models
├── install.sh              # Installation script
├── uninstall.sh           # Uninstallation script
└── README.md
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## ⚠️ Disclaimer

This is a personal password manager. While it uses industry-standard encryption, please:

- **Backup your passwords** elsewhere
- **Remember your master password** (it cannot be recovered)
- **Use at your own risk** for critical passwords

## 🙏 Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) CLI framework
- Uses [Go's crypto libraries](https://pkg.go.dev/crypto) for security
- Inspired by other CLI password managers