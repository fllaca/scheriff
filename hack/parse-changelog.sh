#!/usr/bin/env bash
set -eo pipefail

version=$1
changelog_file=$2
changelog_entry_prefix="^## "


if [[ -z "$version" ]] || [[ -z "$changelog_file" ]]; then
  echo "usage: $0 <version-name> <changelog-file>"
  exit 1
fi

changelog_file_length=$(wc -l < $changelog_file)

# get the start line of the changelog entry of the version
version_entry=$(grep -n "$changelog_entry_prefix${version}" $changelog_file || [[ $? == 1 ]])
release_notes_start=$(echo "$version_entry" | cut -d : -f1)

if [[ -z "$release_notes_start" ]]; then
  echo "Entry with pattern '$changelog_entry_prefix$version' not found in changelog file '$changelog_file'"
  exit 1
fi

# Remove line containing the release title
release_notes_start=$((release_notes_start+1))

# get line number of previous version entry in changelog
previous_version_entry=$(cat $changelog_file | awk "{if (NR>$release_notes_start) print}" | grep -n $changelog_entry_prefix || [[ $? == 1 ]])
release_notes_length=$(echo "$previous_version_entry" | cut -d : -f1)

release_notes_end=$changelog_file_length
if [[ ! -z "$release_notes_length" ]]; then
  release_notes_end=$(($release_notes_start+$release_notes_length-1))
fi

# Print section of Changelog file containing the version release notes
sed "$release_notes_start,$release_notes_end!d" $changelog_file
