#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2026 The jma-openapi contributors
#
# SPDX-License-Identifier: MPL-2.0

set -euo pipefail

repo_root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$repo_root"

copyright_holder="The jma-openapi contributors"
copyright_year="$(date +%Y)"

annotate_file() {
  uvx reuse annotate \
    --copyright "$copyright_holder" \
    --license MPL-2.0 \
    --year "$copyright_year" \
    --fallback-dot-license \
    --skip-existing \
    "$1"
}

annotate_dir() {
  uvx reuse annotate \
    --copyright "$copyright_holder" \
    --license MPL-2.0 \
    --year "$copyright_year" \
    --fallback-dot-license \
    --skip-existing \
    --recursive \
    "$1"
}

uvx reuse download MPL-2.0 MIT

annotate_file .gitignore
annotate_file .golangci.yml
annotate_file Dockerfile
annotate_file README.md
annotate_file compose.yaml
annotate_file go.mod
annotate_file go.sum
annotate_file mise.toml

annotate_dir cmd
annotate_dir deploy
annotate_dir docs
annotate_dir internal
annotate_dir openapi
annotate_dir scripts
annotate_dir tests

uvx reuse annotate \
  --copyright "$copyright_holder" \
  --license MPL-2.0 \
  --year "$copyright_year" \
  --force-dot-license \
  --skip-existing \
  docs/adr/adr-001-language-and-router-selection.md

uvx reuse annotate \
  --copyright "Scalar contributors" \
  --license MIT \
  --force-dot-license \
  --skip-existing \
  web/scalar-api-reference.js
