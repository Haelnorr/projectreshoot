#!/bin/sh
HOOKS_DIR=".githooks"
GIT_HOOKS_DIR=".git/hooks"

mkdir -p "$GIT_HOOKS_DIR"

for hook in "$HOOKS_DIR"/*; do
    hook_name=$(basename "$hook")
    cp "$hook" "$GIT_HOOKS_DIR/$hook_name"
    chmod +x "$GIT_HOOKS_DIR/$hook_name"
done

echo "Git hooks installed!"

