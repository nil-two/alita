### v0.8.0 - 2017-02-21

- Support mixed flag like `-rd'[:/]+'`.
- Interpret an equal sign immediately after a short option as the value of the option
  - later:   -d== #=> DELIM = "=="
  - earlier: -d== #=> DELIM = "="

### v0.7.1 - 2016-01-12

- Change the format of version from "v0.7.1" to "0.7.1".

### v0.7.0 - 2015-11-08

- Allow using `-c`, `--count` even if the delimiter is default.

### v0.6.1 - 2015-07-12

- Allow specifying input-files before flags.

### v0.6.0 - 2015-06-11

- Support `-c`, `--count` to allow specifying how many times to delimit.

### v0.5.1 - 2015-05-23

- Modify messages.

### v0.5.0 - 2015-05-21

- Leave shortest leading spaces.

### v0.4.0 - 2015-05-21

- Remove short option for --version.
- Release compiled binary for Windows, OSX, and Linux.

### v0.3.0 - 2015-05-20

- Support `-j`, `--justify` to allow changing method of justification.
- Suppress trailing null character.

### v0.2.0 - 2015-05-18

- Support `-m`, `--margin` to allow changing margin between each cells.
- Support `-r`, `--regexp` to allow specifying delimiter by regexp.
- Support `-d`, `--delimiter` to allow changing delimiter.
- Support short option for --help.
- Support short option for --version.
- Fix typo of command name.

### v0.1.0 - 2015-05-16

- Initial release.
