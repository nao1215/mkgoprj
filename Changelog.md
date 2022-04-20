# Changelog
All notable changes to this project will be documented in this file.  
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/). This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

# [2.1.0] - 2022-04-19
## Added
- Added function to generate Contributors.md
# [2.0.2] - 2022-04-19
## Changed
- Changed function used for abnormal termination from exitError() to print.Fatal()
# [2.0.0] - 2022-04-19
## Changed
- Rename project name from ubume to mkgoprj
- Interactive mode and Simple Application mode eliminated.
- Introduced cobra library and enabled shell completion.
# [1.5.4] - 2022-03-21
## Changed
- Fixed a bug that could not be built in CLI command project
- Fixed a mistake in the goreleaser settings. There was a mistake in the hierarchy where main.go exists
- Fixed indent in release.yml and .goreleaser.yml
# [1.5.0] - 2022-03-06
## Added
- Added goreleaser setting for ubume.
- Added issue template for ubume.
- For CLI command projects, added the function to generate files related to goreleaser
- Added function to generate code of conduct file
- Added function to generate dependabot.yml
- Added function to generate GitHub issue template files
## Changed
- Changed to be able to make clean even for library projects.
# [1.4.2] - 2022-02-22
## Added
- Added reviewdog.yml in .github/workflow
## Changed
- Changed the CLI package configuration for cobra.
# [1.4.1] - 2022-01-31
## Added
- Added GOARCH to the generated Makefile
# [1.4.0] - 2022-01-29
## Added
- --interactive option: create a project interactively.
# [1.3.0] - 2022-01-25
## Added
- --cli option: create the CLI project with cobra.
# [1.2.1] - 2022-01-24
## Changed
- Changed to output the message to the user.

# [1.2.0] - 2022-01-19
## Added
- --no-root option. Create files in the current directory without creating the project root directory.
# [1.1.0] - 2022-01-19
## Added
- Generate .github/workflow/{build.yml, unit_test.yml}
# [1.0.1] - 2022-01-16
## Added
- AUTHORS.md
- CODE_OF_CONDUCT.md
## Changed
 Restored canMake(). The necessary check process had disappeared, so I restored it.

# [1.0.0] - 2022-01-15
## Added
 - --library option. This option generate library project template.
 - README in japanese.
# [0.5.1] - 2022-01-11
## Added
 - Generate file for test code.
 - "$ make test" support genarating coverage report file.
# [0.5.0] - 2022-01-08
- First release.