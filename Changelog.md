# Changelog
All notable changes to this project will be documented in this file.  
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/). This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
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