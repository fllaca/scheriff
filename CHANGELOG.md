# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Add changes in one of these sections:
* `Added` for new features.
* `Changed` for changes in existing functionality.
* `Deprecated` for soon-to-be removed features.
* `Removed` for now removed features.
* `Fixed` for any bug fixes.
* `Security` in case of vulnerabilities.

## [Unreleased]

## v0.0.1-rc1 - 2020-08-13

This is the very first _SchemaSheriff_ release candidate! :tada:

Please use this version to start playing around with the tool. If you find something to be fixed/improved please you are super welcome to raise an issue in [Scheriff Issues](https://github.com/fllaca/scheriff/issues), your feedback and Pull Requests will be very appreciated.

### Added

- Offline validation of Kubernetes configuration
- Validate multiple input folders/file by setting `-f` multiple times
- Add extra CRD schemas via `--crd` flag (multiple values allowed)
- Recursively validate folders with `-R` flag

