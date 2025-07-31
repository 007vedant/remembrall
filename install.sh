#!/bin/bash

# Remembrall Password Manager Installation Script
# This script builds and installs Remembrall system-wide

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="remembrall"
INSTALL_DIR="/usr/local/bin"
SOURCE_DIR="cmd/remembrall"

# Print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root for system installation
check_permissions() {
    if [[ $EUID -eq 0 ]]; then
        print_info "Running as root - installing system-wide"
        INSTALL_DIR="/usr/local/bin"
    else
        print_warning "Not running as root - installing to user directory"
        INSTALL_DIR="$HOME/.local/bin"
        mkdir -p "$INSTALL_DIR"
    fi
}

# Check system requirements
check_requirements() {
    print_info "Checking system requirements..."
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go 1.19+ from https://golang.org/dl/"
        exit 1
    fi
    
    # Check Go version
    GO_VERSION=$(go version | grep -oE 'go[0-9]+\.[0-9]+' | sed 's/go//')
    MAJOR_VERSION=$(echo $GO_VERSION | cut -d. -f1)
    MINOR_VERSION=$(echo $GO_VERSION | cut -d. -f2)
    
    if [[ $MAJOR_VERSION -lt 1 ]] || [[ $MAJOR_VERSION -eq 1 && $MINOR_VERSION -lt 19 ]]; then
        print_error "Go version $GO_VERSION found. Remembrall requires Go 1.19 or later."
        exit 1
    fi
    
    print_success "Go version $GO_VERSION found"
    
    # Check if we're in the right directory
    if [[ ! -d "$SOURCE_DIR" ]]; then
        print_error "Source directory '$SOURCE_DIR' not found. Please run this script from the project root."
        exit 1
    fi
    
    if [[ ! -f "go.mod" ]]; then
        print_error "go.mod not found. Please run this script from the project root."
        exit 1
    fi
    
    print_success "Project structure verified"
}

# Build the binary
build_binary() {
    print_info "Building Remembrall binary..."
    
    # Get build information
    BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    VERSION="1.0.0"
    
    # Build with ldflags for version info
    go build -ldflags "-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT" \
        -o "$BINARY_NAME" "$SOURCE_DIR/main.go"
    
    if [[ $? -eq 0 ]]; then
        print_success "Binary built successfully"
    else
        print_error "Failed to build binary"
        exit 1
    fi
}

# Install the binary
install_binary() {
    print_info "Installing Remembrall to $INSTALL_DIR..."
    
    # Copy binary to install directory
    if [[ $EUID -eq 0 ]]; then
        cp "$BINARY_NAME" "$INSTALL_DIR/"
        chmod 755 "$INSTALL_DIR/$BINARY_NAME"
    else
        cp "$BINARY_NAME" "$INSTALL_DIR/"
        chmod 755 "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    if [[ $? -eq 0 ]]; then
        print_success "Binary installed to $INSTALL_DIR/$BINARY_NAME"
    else
        print_error "Failed to install binary"
        exit 1
    fi
    
    # Clean up local binary
    rm -f "$BINARY_NAME"
}

# Update PATH if needed
update_path() {
    print_info "Checking PATH configuration..."
    
    # Check if install directory is in PATH
    if [[ ":$PATH:" == *":$INSTALL_DIR:"* ]]; then
        print_success "$INSTALL_DIR is already in PATH"
        return
    fi
    
    if [[ $EUID -eq 0 ]]; then
        print_success "/usr/local/bin should be in system PATH by default"
        return
    fi
    
    # For user installation, add to PATH
    print_warning "$INSTALL_DIR is not in your PATH"
    print_info "Adding PATH configuration to shell profile..."
    
    # Detect shell and add to appropriate profile
    SHELL_NAME=$(basename "$SHELL")
    case "$SHELL_NAME" in
        bash)
            PROFILE_FILE="$HOME/.bashrc"
            if [[ -f "$HOME/.bash_profile" ]]; then
                PROFILE_FILE="$HOME/.bash_profile"
            fi
            ;;
        zsh)
            PROFILE_FILE="$HOME/.zshrc"
            ;;
        fish)
            PROFILE_FILE="$HOME/.config/fish/config.fish"
            mkdir -p "$(dirname "$PROFILE_FILE")"
            ;;
        *)
            PROFILE_FILE="$HOME/.profile"
            ;;
    esac
    
    # Add PATH export if not already present
    if ! grep -q "$INSTALL_DIR" "$PROFILE_FILE" 2>/dev/null; then
        echo "" >> "$PROFILE_FILE"
        echo "# Added by Remembrall installer" >> "$PROFILE_FILE"
        echo "export PATH=\"$INSTALL_DIR:\$PATH\"" >> "$PROFILE_FILE"
        print_success "Added $INSTALL_DIR to PATH in $PROFILE_FILE"
        print_warning "Please restart your terminal or run: source $PROFILE_FILE"
    else
        print_success "PATH already configured in $PROFILE_FILE"
    fi
}

# Verify installation
verify_installation() {
    print_info "Verifying installation..."
    
    # Check if binary exists and is executable
    if [[ -x "$INSTALL_DIR/$BINARY_NAME" ]]; then
        print_success "Binary is installed and executable"
    else
        print_error "Binary installation verification failed"
        exit 1
    fi
    
    # Test basic functionality
    if "$INSTALL_DIR/$BINARY_NAME" --help &> /dev/null; then
        print_success "Binary is working correctly"
    else
        print_error "Binary execution test failed"
        exit 1
    fi
    
    # Show version info
    VERSION_OUTPUT=$("$INSTALL_DIR/$BINARY_NAME" --version 2>/dev/null || echo "Version info not available")
    print_info "Installed version: $VERSION_OUTPUT"
}

# Show usage instructions
show_usage() {
    echo ""
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}                             REMEMBRALL INSTALLATION COMPLETE!                               ${NC}"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo -e "${BLUE}🔐 Remembrall Password Manager is now installed and ready to use!${NC}"
    echo ""
    echo -e "${YELLOW}Getting Started:${NC}"
    echo "  remembrall save github          # Save your first password"
    echo "  remembrall list                 # List all stored applications"
    echo "  remembrall get github           # Retrieve a password"
    echo "  remembrall search git           # Search with fuzzy matching"
    echo "  remembrall --help               # Show all available commands"
    echo ""
    echo -e "${YELLOW}Security Features:${NC}"
    echo "  • AES-256-GCM encryption with PBKDF2 key derivation"
    echo "  • Master password authentication for all operations"
    echo "  • Hidden password input (no shoulder surfing)"
    echo "  • 5-second auto-clear for retrieved passwords"
    echo "  • SQLite database stored in your home directory"
    echo ""
    echo -e "${YELLOW}Installation Location:${NC} $INSTALL_DIR/$BINARY_NAME"
    echo -e "${YELLOW}Database Location:${NC} ~/.remembrall.db (created on first use)"
    echo ""
    if [[ $EUID -ne 0 ]]; then
        echo -e "${YELLOW}Note:${NC} If 'remembrall' command is not found, restart your terminal or run:"
        echo "  source $PROFILE_FILE"
        echo ""
    fi
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

# Main installation flow
main() {
    echo -e "${BLUE}╔══════════════════════════════════════════════════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║                                    REMEMBRALL PASSWORD MANAGER                                               ║${NC}"
    echo -e "${BLUE}║                                         INSTALLATION SCRIPT                                                  ║${NC}"
    echo -e "${BLUE}╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    
    check_permissions
    check_requirements
    build_binary
    install_binary
    update_path
    verify_installation
    show_usage
}

# Handle script interruption
trap 'print_error "Installation interrupted"; exit 1' INT TERM

# Run main function
main "$@"