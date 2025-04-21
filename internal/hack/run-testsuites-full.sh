#!/bin/bash

export TESTING_OXIGRAPH_EXEC="${PWD}/tmp/bin/oxigraph"

if ! [ -x "$TESTING_OXIGRAPH_EXEC" ]; then
  unset TESTING_OXIGRAPH_EXEC
fi

run() {
  go test -json . | jq -rs \
    --argjson resultEmojis '{ "pass": "✅", "fail": "❌", "skip": "❔" }' \
    '
      map(select(.Test and ( .Test | split("/") | length > 2 ))) | group_by(.Test)
      | map(
        (.[0].Test | split("/") | ({ testGroup: .[1], testName: .[2] }))
        + (map(select(.Output) | .Output | match("^(    ([^\\s]+: |    ))(.+)"; "m") | .captures[2].string) | if length > 0 then {"output": join("")} else {} end)
        + ({ result: map(select(.Action == "pass" or .Action == "fail" or .Action == "skip"))[0].Action })
      )
      | sort_by(.testGroup, .testName)
      | (
          "# Test Suite"
          , ""
          , (
            group_by(.result) | map({ key: .[0].result, value: length }) | from_entries
              | "\(.pass // 0) passed, \(.fail // 0) failed, \(.skip // 0) skipped"
          )
          , ""
          , "| Result | Test Group | Test Name |"
          , "|:------ |:---------- |:--------- |"
          , map(
              "| \($resultEmojis[.result])&nbsp;\(.result | ascii_upcase) | \(.testGroup) | \(.testName) |"
            )[]
          , ""
          , (map(select(.output)) | (if length > 0 then (
            "## Output\n"
            + "\n"
            + (map(
              "### \($resultEmojis[.result]) \(.testGroup), \(.testName)\n"
              + "\n"
              + "```\n"
              + .output
              + "```\n"
            ) | join("\n"))
          ) else "" end))
      )
    ' \
    > RESULTS.md
}

while read -r dir
do
  echo "$dir"
  pushd "$dir" > /dev/null
  run
  popd > /dev/null
done < <(
  find . -type f -name testsuite_test.go -exec dirname {} \;
)