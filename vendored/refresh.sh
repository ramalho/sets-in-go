#!/bin/bash
# Re-extract the proposed container/set, container/mapset and container/hash
# packages from Gerrit into this directory as a self-contained Go module.
#
# All CLs are still open, so patchset numbers are pinned below for
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
MAP_CL=612217 # container/hash.Map, on which hash.Set is built
MAP_PS=${MAP_PS:-30}
HASH_CL=741160 # container/hash.Set
HASH_PS=${HASH_PS:-21}
ITER_CL=745440 # iter: add Every, Some -- see internal/iter (not fetched)

# The module needs hash/maphash.Hasher, new in Go 1.27.
GO=${GO:-go1.27rc2}
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

# container/hash.Set is defined in terms of the unexported internals of
# container/hash.Map, so the two CLs must be extracted together into one
# package. The set CL also rewrites the map CL's example_test.go, so take
# example_test.go from the later change.
echo "Fetching container/hash.Map CL $MAP_CL patchset $MAP_PS"
mkdir -p hash
fetch "$MAP_CL" "$MAP_PS" src/container/hash/map.go hash/map.go
fetch "$MAP_CL" "$MAP_PS" src/container/hash/map_test.go hash/map_test.go
fetch "$MAP_CL" "$MAP_PS" src/container/hash/iter_test.go hash/iter_test.go

echo "Fetching container/hash.Set CL $HASH_CL patchset $HASH_PS"
fetch "$HASH_CL" "$HASH_PS" src/container/hash/set.go hash/set.go
fetch "$HASH_CL" "$HASH_PS" src/container/hash/set_test.go hash/set_test.go
fetch "$HASH_CL" "$HASH_PS" src/container/hash/example_test.go hash/example_test.go

# The conformance test that checks set.Set and hash.Set against the same
# abstract interfaces. Upstream it lives in src/container and is touched by
# both CLs; take the later one.
echo "Fetching container conformance test CL $HASH_CL patchset $HASH_PS"
mkdir -p container
fetch "$HASH_CL" "$HASH_PS" src/container/container_test.go container/container_test.go

# mapset.String and hash.{Map,Set}.String use internal/fmtsort, which cannot be
# imported from outside GOROOT. It only depends on cmp, reflect and slices, so
# it vendors verbatim -- but hash needs fmtsort.Compare, which the Map CL adds,
# so take the file from that CL rather than from the local GOROOT.
echo "Fetching internal/fmtsort from CL $MAP_CL patchset $MAP_PS"
mkdir -p internal/fmtsort
fetch "$MAP_CL" "$MAP_PS" src/internal/fmtsort/sort.go internal/fmtsort/sort.go
cp "$($GO env GOROOT)/LICENSE" LICENSE

# Rewrite the stdlib import paths to their vendored equivalents.
#
# mapset calls maps.Identical, which is itself an unmerged proposal (CL 760800,
# golang/go#78456) and so is absent from every released Go. The checked-in
# internal/maps shim supplies it and forwards the rest to stdlib, which lets the
# extracted function bodies stay verbatim.
#
# hash calls iter.Seq.Every and iter.Seq.Some, methods from another unmerged
# proposal (CL 745440). Those cannot be shimmed as methods -- see the comment at
# the top of internal/iter/iter.go -- so the three call sites are rewritten from
# method to function form below.
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
sed -i '' \
	-e "s|\"internal/fmtsort\"|\"$MODULE/internal/fmtsort\"|" \
	-e "s|^\([[:space:]]*\)\"iter\"\$|\1\"$MODULE/internal/iter\"|" \
	hash/map.go hash/set.go
sed -i '' \
	-e "s|\"container/hash\"|\"$MODULE/hash\"|" \
	hash/map_test.go hash/set_test.go hash/iter_test.go hash/example_test.go
sed -i '' \
	-e "s|\"container/hash\"|\"$MODULE/hash\"|" \
	-e "s|\"container/set\"|\"$MODULE/set\"|" \
	container/container_test.go

echo "Rewriting iter.Seq.Every/Some calls to function form"
sed -i '' -e "s|keys\.Every(m\.Contains)|iter.Every(keys, m.Contains)|" hash/map.go
sed -i '' \
	-e "s|seq\.Every(s\.Contains)|iter.Every(seq, s.Contains)|" \
	-e "s|other\.All()\.Some(s\.Contains)|iter.Some(other.All(), s.Contains)|" \
	hash/set.go

# Fail loudly if a new patchset introduces a call these rules do not cover.
if grep -n '\.\(Every\|Some\)(' hash/*.go container/*.go |
	grep -v 'iter\.\(Every\|Some\)('; then
	echo "error: unrewritten iter.Seq method call above; update refresh.sh" >&2
	exit 1
fi

gofmt -w set/set.go mapset/mapset.go mapset/mapset_test.go hash/*.go container/container_test.go

cat >PROVENANCE.txt <<EOF
Extracted by refresh.sh on $(date -u '+%Y-%m-%d %H:%M:%S UTC')

  container/set            $GERRIT/c/go/+/$SET_CL  patchset $SET_PS
  container/mapset         $GERRIT/c/go/+/$MAPSET_CL  patchset $MAPSET_PS
  container/hash (Map)     $GERRIT/c/go/+/$MAP_CL  patchset $MAP_PS
  container/hash (Set)     $GERRIT/c/go/+/$HASH_CL  patchset $HASH_PS
  container (conformance)  $GERRIT/c/go/+/$HASH_CL  patchset $HASH_PS
  internal/fmtsort         $GERRIT/c/go/+/$MAP_CL  patchset $MAP_PS

All CLs were unmerged proposals at the time of extraction.
Import paths were changed throughout; hash/map.go and hash/set.go also have
three iter.Seq method calls rewritten to function form (see internal/iter).
The code is otherwise verbatim.

Built and tested with $($GO version | cut -d' ' -f3-4).
EOF

$GO build ./... && $GO test ./... >/dev/null && echo "OK - builds and tests pass"
