#!/bin/sh

# inject examples/tour/main_test.go into README.md

set -eu

README='README.md'
USAGE='examples/tour/main_test.go'
START='<!-- USAGE-START -->'
END='<!-- USAGE-END -->'

# create temp file
TMP=$(mktemp /tmp/update-readme.XXXXXX) || exit 1

awk -v start="$START" -v end="$END" -v usage="$USAGE" '
  {
    if ($0 == start) {
      print start
      print "```go"
      # read & indent usage file
      while (getline line < usage > 0) {
        print line
      }
      close(usage)
      print "```"
      print end
      skip = 1
    }
    else if ($0 == end) {
      skip = 0
    }
    else if (!skip) {
      print
    }
  }
' "$README" >"$TMP"

# overwrite README
mv "$TMP" "$README"

printf '✔ README.md updated\n'
