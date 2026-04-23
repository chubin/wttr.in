#!/usr/bin/env bash
set -euo pipefail

# ──────────────────────────────────────────────────────────────────────────────
# build.sh — Build helper script for the project
# Usage:
#   ./build.sh build
#   ./build.sh clean
#   ./build.sh assets
#   ./build.sh all          # assets + build
#   ./build.sh help
# ──────────────────────────────────────────────────────────────────────────────

# ─── Configuration ───────────────────────────────────────────────────────────

BINARY_NAME="srv"
MAIN_PACKAGE="./main.go"           # adjust if your entry point is different
EMBED_TARGET_DIR="internal/assets/embed"

# ─── Font configuration ───────────────────────────────────────────────────────
declare -A FONTS
FONTS["DejaVuSansMono.ttf"]="/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf:fonts-dejavu-core"
FONTS["wqy-zenhei.ttc"]="/usr/share/fonts/truetype/wqy/wqy-zenhei.ttc:fonts-wqy-zenhei"
FONTS["MTLc3m.ttf"]="/usr/share/fonts/truetype/motoya-l-cedar/MTLc3m.ttf:fonts-motoya-l-cedar"
FONTS["LexiGulim.ttf"]="/usr/share/fonts/truetype/lexi/LexiGulim.ttf:fonts-lexi-gulim"
FONTS["Symbola_hint.ttf"]="/usr/share/fonts/truetype/ancient-scripts/Symbola_hint.ttf:fonts-symbola"
FONTS["NotoSansDevanagari-Regular.ttf"]="/usr/share/fonts/truetype/noto/NotoSansDevanagari-Regular.ttf:fonts-noto-core"
FONTS["NotoSansBengali-Regular.ttf"]="/usr/share/fonts/truetype/noto/NotoSansBengali-Regular.ttf:fonts-noto-core"

# ─── Asset roots (non-font) ───────────────────────────────────────────────────
ASSET_ROOTS=(
    "spec"
    "share/static"
    "share/pages"
    "share/templates"
)

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# ─── Helper functions ────────────────────────────────────────────────────────

info()  { echo -e "${GREEN}→ $1${NC}"; }
warn()  { echo -e "${YELLOW}⚠ $1${NC}" >&2; }
error() { echo -e "${RED}ERROR: $1${NC}" >&2; exit 1; }

ensure_clean_embed_dir() {
    mkdir -p "$EMBED_TARGET_DIR/fonts"
}

copy_fonts() {
    info "Copying fonts to embed directory..."

    for font_file in "${!FONTS[@]}"; do
        entry="${FONTS[$font_file]}"
        src_path="${entry%%:*}"
        pkg="${entry#*:}"

        if [[ -f "$src_path" ]]; then
            cp "$src_path" "$EMBED_TARGET_DIR/fonts/$font_file"
            info "  Copied $font_file"
        else
            warn "Font not found: $font_file"
            warn "   → Install with: sudo apt install $pkg"
        fi
    done
}

copy_assets() {
    ensure_clean_embed_dir

    # Copy regular assets
    for src_root in "${ASSET_ROOTS[@]}"; do
        if [[ ! -d "$src_root" ]]; then
            warn "Asset root not found (skipping): $src_root"
            continue
        fi

        target_dir="${EMBED_TARGET_DIR}/${src_root}"
        mkdir -p "$target_dir"

        # Use rsync to copy directory structure efficiently
        # -a = archive mode (preserves permissions, times, symlinks, etc.)
        # --delete = remove files in target that no longer exist in source
        # --exclude = optional patterns (add more if needed)
        rsync -a --delete \
            --exclude='*.tmp' \
            --exclude='*.bak' \
            --exclude='.DS_Store' \
            --exclude='Thumbs.db' \
            "${src_root}/" \
            "${target_dir}/" \
        || error "rsync failed for $src_root → $target_dir"

        info "Synced assets: $src_root/ → $target_dir/"
    done

    # Copy fonts
    copy_fonts

    info "All assets prepared for embedding"
}

# ─── Actions ─────────────────────────────────────────────────────────────────

cmd_assets() {
    info "Preparing embedded assets..."
    copy_assets
    info "Assets ready"
}

cmd_build() {
    info "Building binary: $BINARY_NAME"

    # Optional: force assets preparation before build
    cmd_assets

    local ldflags=""
    # Example: embed version / build time
    # ldflags="-X main.buildVersion=$(git describe --tags --dirty --always) -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)"

    CGO_CFLAGS="-Wno-return-local-addr" \
    go build \
        -ldflags "$ldflags" \
        -o "$BINARY_NAME" \
        "$MAIN_PACKAGE" \
    || error "Build failed"

    if [[ -f "$BINARY_NAME" ]]; then
        info "Build successful: ./$BINARY_NAME"
    else
        error "Binary not created"
    fi
}

cmd_update_css() {
    curl -L -o share/static/terminal.css \
      https://raw.githubusercontent.com/buildkite/terminal-to-html/refs/heads/main/internal/assets/terminal.css
    sed -i 's@^[.]term-container {@.term-container-disabled {@' share/static/terminal.css
}

cmd_gen() {
    info "Generating..."
    ./"$BINARY_NAME" gen
    go fmt internal/options/options.go
    info "Generating done"
}

cmd_clean() {
    info "Cleaning up..."
    rm -f "$BINARY_NAME"
    rm -rf "$EMBED_TARGET_DIR"
    info "Clean done"
}

cmd_all() {
    cmd_assets
    cmd_build
}

cmd_help() {
    cat <<EOF
Usage: ./build.sh <command>

Commands:
  assets      Prepare/copy embeddable assets
  build       Build the binary (implies assets)
  all         assets + build
  clean       Remove build artifacts and copied assets
  help        Show this help

Environment variables that affect build:
  CGO_CFLAGS  (already set to -Wno-return-local-addr by default)
EOF
}

# ─── Main dispatcher ─────────────────────────────────────────────────────────

main() {
    if [[ $# -eq 0 ]]; then
        cmd_help
        exit 0
    fi

    local cmd="$1"
    shift

    case "$cmd" in
        assets)     cmd_assets "$@" ;;
        build)      cmd_build "$@" ;;
        gen)        cmd_gen "$@" ;;
        all)        cmd_all "$@" ;;
        clean)      cmd_clean "$@" ;;
        update-css) cmd_update_css "$@" ;;
        help|--help|-h)
                    cmd_help ;;
        *)
            error "Unknown command: $cmd"
            ;;
    esac
}

main "$@"
