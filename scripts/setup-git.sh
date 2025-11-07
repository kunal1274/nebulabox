#!/bin/bash

# Setup Git configuration for NebulaBox development
# This ensures Git is properly configured for releases

set -e

echo "🔧 Setting up Git configuration for NebulaBox"
echo ""

# Check current configuration
CURRENT_NAME=$(git config --get user.name 2>/dev/null || echo "")
CURRENT_EMAIL=$(git config --get user.email 2>/dev/null || echo "")

if [ -n "$CURRENT_NAME" ] && [ -n "$CURRENT_EMAIL" ]; then
    echo "✅ Git is already configured:"
    echo "   Name:  $CURRENT_NAME"
    echo "   Email: $CURRENT_EMAIL"
    echo ""
    read -p "Do you want to change it? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "✅ Keeping current configuration"
        exit 0
    fi
fi

# Get user input
echo "Enter your Git configuration:"
echo ""

# Get name
if [ -z "$CURRENT_NAME" ]; then
    read -p "Your Name: " GIT_NAME
else
    read -p "Your Name [$CURRENT_NAME]: " GIT_NAME
    GIT_NAME=${GIT_NAME:-$CURRENT_NAME}
fi

# Get email
if [ -z "$CURRENT_EMAIL" ]; then
    read -p "Your Email: " GIT_EMAIL
else
    read -p "Your Email [$CURRENT_EMAIL]: " GIT_EMAIL
    GIT_EMAIL=${GIT_EMAIL:-$CURRENT_EMAIL}
fi

# Ask for scope
echo ""
echo "Configure for:"
echo "  1) This repository only (recommended for testing)"
echo "  2) All repositories (global)"
read -p "Choice [1]: " SCOPE
SCOPE=${SCOPE:-1}

# Set configuration
if [ "$SCOPE" = "2" ]; then
    git config --global user.name "$GIT_NAME"
    git config --global user.email "$GIT_EMAIL"
    echo ""
    echo "✅ Git configured globally"
else
    git config user.name "$GIT_NAME"
    git config user.email "$GIT_EMAIL"
    echo ""
    echo "✅ Git configured for this repository only"
fi

# Verify
echo ""
echo "Current Git configuration:"
echo "   Name:  $(git config --get user.name)"
echo "   Email: $(git config --get user.email)"
echo ""
echo "✅ Git setup complete! You can now create releases."

