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

## [v0.0.1-rc2] - 2020-08-25

Added the ability for strict validation (don't accept warnings) and some fixes to match the Kubernetes behaviour when veriyfing `null` fields and additional properties.

Thank you so much to @LeoVerto for his contributions and to @figuerascarlos for his thorough reviews! :heart:

### Fixed

- Fix nullable fields validation ([#13](https://github.com/fllaca/scheriff/pull/13))
- Fix validating files without yaml extension (like when doing process substitution) ([#14](https://github.com/fllaca/scheriff/pull/14))
- Fix false OK validations for configuration with additional properties ([#11](https://github.com/fllaca/scheriff/pull/11))

### Added

- Added a `--strict` flag to make _scheriff_ fail when errors but also warnings are encountered ([#9](https://github.com/fllaca/scheriff/pull/9))
- Published Docker image: [quay.io/fllaca/scheriff](https://quay.io/repository/fllaca/scheriff?tab=tags)

## [v0.0.1-rc1] - 2020-08-13

This is the very first _SchemaSheriff_ release candidate! :tada:

Please use this version to start playing around with the tool. If you find something to be fixed/improved please you are super welcome to raise an issue in [Scheriff Issues](https://github.com/fllaca/scheriff/issues), your feedback and Pull Requests will be very appreciated.

### Added

- Offline validation of Kubernetes configuration
- Validate multiple input folders/file by setting `-f` multiple times
- Add extra CRD schemas via `--crd` flag (multiple values allowed)
- Recursively validate folders with `-R` flag

