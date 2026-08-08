#!/usr/bin/env python3
"""Replace tabs with spaces in .go files, recursively.

By default it walks the directory this script lives in and rewrites every
`.go` file below it in place, expanding each tab into 4 spaces.

    ./untab.py                  # rewrite snippets/**/*.go
    ./untab.py --dry-run        # report what would change, touch nothing
    ./untab.py -w 2 some/dir    # 2 spaces per tab, different root
    ./untab.py --leading-only   # only tabs in the indentation of a line
"""

import argparse
import sys
from pathlib import Path


def expand(text: str, width: int, leading_only: bool) -> str:
    """Return text with tabs replaced by `width` spaces each."""
    spaces = " " * width
    if not leading_only:
        return text.replace("\t", spaces)

    lines = []
    for line in text.splitlines(keepends=True):
        stripped = line.lstrip("\t")
        indent = len(line) - len(stripped)
        lines.append(spaces * indent + stripped)
    return "".join(lines)


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__,
                                     formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument("root", nargs="?", default=Path(__file__).parent, type=Path,
                        help="directory to walk (default: the directory of this script)")
    parser.add_argument("-w", "--width", type=int, default=4,
                        help="spaces per tab (default: 4)")
    parser.add_argument("-n", "--dry-run", action="store_true",
                        help="report files that would change, without writing")
    parser.add_argument("--leading-only", action="store_true",
                        help="expand only tabs used for indentation, "
                             "preserving gofmt's mid-line alignment tabs")
    args = parser.parse_args()

    if not args.root.is_dir():
        print(f"untab: {args.root}: not a directory", file=sys.stderr)
        return 1

    changed = 0
    for path in sorted(args.root.rglob("*.go")):
        old = path.read_text(encoding="utf-8")
        new = expand(old, args.width, args.leading_only)
        if new == old:
            continue
        changed += 1
        tabs = old.count("\t")
        print(f"{'would rewrite' if args.dry_run else 'rewrote'} {path} ({tabs} tabs)")
        if not args.dry_run:
            path.write_text(new, encoding="utf-8")

    print(f"{changed} file(s) {'to change' if args.dry_run else 'changed'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
