#!/usr/bin/env bash
# Refactor all occurrences of “testo” → “testo” (module path, imports, package names, other mentions).
# macOS-compatible: uses BSD sed (-i '') and skips vendor directory entirely.
set -euox pipefail

OLD_MOD="github.com/metafates/testo"
NEW_MOD="github.com/metafates/testo"
SHORT_OLD="testo"
SHORT_NEW="testo"

# 1) Update go.mod’s module declaration exactly.
sed -i '' "s|^module[[:space:]]\+$OLD_MOD\$|module $NEW_MOD|" go.mod

# 2) Delete go.sum to avoid stale checksums.
rm -f go.sum

# 3) Walk every file (excluding vendor), replacing:
#      - Full module path OLD_MOD → NEW_MOD
#      - Every “testo” substring → “testo” (covers imports, package, comments, strings, etc.)
find . \
	-type d -name vendor -prune -false -o \
	-type f -print |
	xargs -I{} sed -i '' -e "s|$OLD_MOD|$NEW_MOD|g" -e "s|$SHORT_OLD|$SHORT_NEW|g" {}

# 4) Regenerate go.sum (and tidy).
go mod tidy

echo "✔ Complete: all ‘$SHORT_OLD’ references → ‘$SHORT_NEW’, module path updated."
