#!/bin/bash
# Re-extract Alan Donovan's proposed container/set and container/mapset
# packages from Gerrit into this directory as a self-contained Go module.
#
# Both CLs are still open, so patchset numbers are pinned below for
# reproducibility. Bump them (or set them in the environment) to pick up
# newer revisions:
#
#     SET_PS=14 MAPSET_PS=28 ./refresh.sh
#
# Pass "current" to follow whatever the latest patchset is:
#
#     SET_PS=current MAPSET_PS=current ./refresh.sh

set -euo pipefail

SET_CL=745441
SET_PS=${SET_PS:-13}
MAPSET_CL=724420
MAPSET_PS=${MAPSET_PS:-27}

MODULE=github.com/ramalho/sets-in-go/vendored
GERRIT=https://go-review.googlesource.com
cd "$(dirname "$0")"

# The /content endpoint returns the file as plain base64 (unlike Gerrit's JSON
# endpoints, which are prefixed with )]}' to defeat JSON hijacking).
fetch() { # fetch <change> <patchset> <path-in-go-repo> <dest>
	local url="$GERRIT/changes/$1/revisions/$2/files/${3//\//%2F}/content"
	echo "  $3 -> $4"
	curl -sfS "$url" | base64 -d >"$4"
}

echo "Fetching container/set CL $SET_CL patchset $SET_PS"
mkdir -p set
fetch "$SET_CL" "$SET_PS" src/container/set/set.go set/set.go

echo "Fetching container/mapset CL $MAPSET_CL patchset $MAPSET_PS"
mkdir -p mapset
fetch "$MAPSET_CL" "$MAPSET_PS" src/container/mapset/mapset.go mapset/mapset.go
fetch "$MAPSET_CL" "$MAPSET_PS" src/container/mapset/mapset_test.go mapset/mapset_test.go

# mapset.String uses internal/fmtsort, which cannot be imported from outside
# GOROOT. It only depends on cmp, reflect and slices, so it vendors verbatim.
echo "Copying internal/fmtsort from $(go env GOROOT)"
mkdir -p internal/fmtsort
cp "$(go env GOROOT)/src/internal/fmtsort/sort.go" internal/fmtsort/sort.go
cp "$(go env GOROOT)/LICENSE" LICENSE

# Rewrite the three stdlib import paths to their vendored equivalents.
# mapset also calls maps.Identical, which is itself an unmerged proposal
# (CL 760800, golang/go#78456) and so is absent from every released Go. The
# checked-in internal/maps shim supplies it and forwards the rest to stdlib,
# which lets the extracted function bodies stay verbatim.
echo "Rewriting import paths"
export LC_ALL=C
sed -i '' \
	-e "s|\"internal/fmtsort\"|\"$MODULE/internal/fmtsort\"|" \
	-e "s|^\([[:space:]]*\)\"maps\"\$|\1\"$MODULE/internal/maps\"|" \
	mapset/mapset.go
sed -i '' \
	-e "s|^\([[:space:]]*\)\"maps\"\$|\1\"$MODULE/internal/maps\"|" \
	-e "s|\"container/mapset\"|\"$MODULE/mapset\"|" \
	mapset/mapset_test.go
sed -i '' \
	-e "s|\"container/mapset\"|\"$MODULE/mapset\"|" \
	set/set.go

gofmt -w set/set.go mapset/mapset.go mapset/mapset_test.go

cat >PROVENANCE.txt <<EOF
Extracted by refresh.sh on $(date -u '+%Y-%m-%d %H:%M:%S UTC')

  container/set     $GERRIT/c/go/+/$SET_CL     patchset $SET_PS
  container/mapset  $GERRIT/c/go/+/$MAPSET_CL  patchset $MAPSET_PS
  internal/fmtsort  $(go env GOROOT)/src/internal/fmtsort/sort.go ($(go version | cut -d' ' -f3))

Both CLs were unmerged proposals at the time of extraction.
Only the import paths were changed; the code is otherwise verbatim.
EOF

go build ./... && go test ./... >/dev/null && echo "OK - builds and tests pass"
