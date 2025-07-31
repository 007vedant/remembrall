#!/bin/bash

# Remembrall Password Manager Uninstallation Script
# This script removes Remembrall from the system

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="remembrall"
SYSTEM_INSTALL_DIR="/usr/local/bin"
USER_INSTALL_DIR="$HOME/.local/bin"
DB_FILE="$HOME/.remembrall.db"
MASTER_FILE="$HOME/.remembrall-master"

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

# Ask for confirmation
confirm_uninstall() {
    echo -e "${YELLOW}⚠️  WARNING: This will remove Remembrall and all its data!${NC}"
    echo ""
    echo "This will delete:"
    echo "  • Remembrall binary"
    echo "  • All stored passwords ($DB_FILE)"
    echo "  • Master password verification ($MASTER_FILE)"
    echo "  • PATH configuration (if added during installation)"
    echo ""
    read -p "Are you sure you want to uninstall Remembrall? (y/N): " -r
    echo ""
    
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Uninstallation cancelled"
        exit 0
    fi
}

# Remove binary
remove_binary() {
    print_info "Removing Remembrall binary..."
    
    REMOVED=false
    
    # Check system installation
    if [[ -f "$SYSTEM_INSTALL_DIR/$BINARY_NAME" ]]; then
        if [[ $EUID -eq 0 ]]; then
            rm -f "$SYSTEM_INSTALL_DIR/$BINARY_NAME"
            print_success "Removed system binary: $SYSTEM_INSTALL_DIR/$BINARY_NAME"
            REMOVED=true
        else
            print_error "System binary found but requires root privileges to remove"
            print_info "Please run: sudo rm $SYSTEM_INSTALL_DIR/$BINARY_NAME"
        fi
    fi
    
    # Check user installation
    if [[ -f "$USER_INSTALL_DIR/$BINARY_NAME" ]]; then
        rm -f "$USER_INSTALL_DIR/$BINARY_NAME"
        print_success "Removed user binary: $USER_INSTALL_DIR/$BINARY_NAME"
        REMOVED=true
    fi
    
    if [[ $REMOVED == false ]]; then
        print_warning "No Remembrall binary found in common installation directories"
    fi
}

# Remove data files
remove_data() {
    print_info "Removing Remembrall data files..."
    
    REMOVED_FILES=()
    
    # Remove database
    if [[ -f "$DB_FILE" ]]; then
        rm -f "$DB_FILE"
        REMOVED_FILES+=("Password database")
        print_success "Removed password database: $DB_FILE"
    fi
    
    # Remove master password file
    if [[ -f "$MASTER_FILE" ]]; then
        rm -f "$MASTER_FILE"
        REMOVED_FILES+=("Master password verification")
        print_success "Removed master password file: $MASTER_FILE"
    fi
    
    if [[ ${#REMOVED_FILES[@]} -eq 0 ]]; then
        print_warning "No Remembrall data files found"
    else
        print_info "Removed ${#REMOVED_FILES[@]} data file(s)"
    fi
}

# Remove PATH configuration
remove_path_config() {
    print_info "Checking for PATH configuration..."
    
    # List of common shell profile files
    PROFILE_FILES=(
        "$HOME/.bashrc"
        "$HOME/.bash_profile"
        "$HOME/.zshrc"
        "$HOME/.profile"
        "$HOME/.config/fish/config.fish"
    )
    
    MODIFIED_FILES=()
    
    for PROFILE_FILE in "${PROFILE_FILES[@]}"; do
        if [[ -f "$PROFILE_FILE" ]] && grep -q "Added by Remembrall installer" "$PROFILE_FILE" 2>/dev/null; then
            # Create backup
            cp "$PROFILE_FILE" "$PROFILE_FILE.remembrall-backup"
            
            # Remove Remembrall-related lines
            sed -i.tmp '/# Added by Remembrall installer/,+1d' "$PROFILE_FILE"
            rm -f "$PROFILE_FILE.tmp"
            
            MODIFIED_FILES+=("$PROFILE_FILE")
            print_success "Removed PATH configuration from: $PROFILE_FILE"
            print_info "Backup created: $PROFILE_FILE.remembrall-backup"
        fi
    done
    
    if [[ ${#MODIFIED_FILES[@]} -eq 0 ]]; then
        print_info "No PATH configuration found to remove"
    else
        print_warning "Please restart your terminal for PATH changes to take effect"
    fi
}

# Verify removal
verify_removal() {
    print_info "Verifying removal..."
    
    REMAINING_ITEMS=()
    
    # Check for remaining binaries
    if command -v $BINARY_NAME &> /dev/null; then
        REMAINING_ITEMS+=("Binary still accessible via PATH")
    fi
    
    if [[ -f "$SYSTEM_INSTALL_DIR/$BINARY_NAME" ]]; then
        REMAINING_ITEMS+=("System binary: $SYSTEM_INSTALL_DIR/$BINARY_NAME")
    fi
    
    if [[ -f "$USER_INSTALL_DIR/$BINARY_NAME" ]]; then
        REMAINING_ITEMS+=("User binary: $USER_INSTALL_DIR/$BINARY_NAME")
    fi
    
    # Check for remaining data
    if [[ -f "$DB_FILE" ]]; then
        REMAINING_ITEMS+=("Database: $DB_FILE")
    fi
    
    if [[ -f "$MASTER_FILE" ]]; then
        REMAINING_ITEMS+=("Master file: $MASTER_FILE")
    fi
    
    if [[ ${#REMAINING_ITEMS[@]} -eq 0 ]]; then
        print_success "Remembrall has been completely removed"
    else
        print_warning "Some items may still remain:"
        for item in "${REMAINING_ITEMS[@]}"; do
            echo "  • $item"
        done
    fi
}

# Show completion message
show_completion() {
    echo ""
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}                             REMEMBRALL UNINSTALLATION COMPLETE!                             ${NC}"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo -e "${BLUE}🔐 Remembrall Password Manager has been uninstalled${NC}"
    echo ""
    echo -e "${YELLOW}What was removed:${NC}"
    echo "  • Remembrall binary"
    echo "  • All stored passwords and master password"
    echo "  • PATH configuration (if any)"
    echo ""
    echo -e "${YELLOW}Notes:${NC}"
    echo "  • Shell profile backups were created before modification"
    echo "  • Restart your terminal to ensure PATH changes take effect"
    echo "  • Thank you for using Remembrall! 🙏"
    echo ""
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

# Main uninstallation flow
main() {
    echo -e "${RED}╔══════════════════════════════════════════════════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${RED}║                                    REMEMBRALL PASSWORD MANAGER                                               ║${NC}"
    echo -e "${RED}║                                       UNINSTALLATION SCRIPT                                                  ║${NC}"
    echo -e "${RED}╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    
    confirm_uninstall
    remove_binary
    remove_data
    remove_path_config
    verify_removal
    show_completion
}

# Handle script interruption
trap 'print_error "Uninstallation interrupted"; exit 1' INT TERM

# Run main function
main "$@"