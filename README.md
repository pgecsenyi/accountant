# Accountant

Matches files by checksum and stores the old name - new name pairs as well as the corresponding checksums in a CSV file.

## Installation and usage

Use one of the build files or _Visual Studio Code_ to build the program. This will provide you one executable in the _bin_ folder. You can also use the `go run` command of course.

The application can be used in two different modes. First, called with the following parameters it will produce a CSV file containing the checksums of the files in the given directory.

    -source "SomeDirectory" -alg sha1 -outchk "Output.csv" -bp "SomeDirectory/"

When called with the following parameters it will produce a file containing old name - new name pairs as well as a new CSV file with the updated filenames.

    -source "SomeDirectory" -alg sha1 -outchk "Output-updated.csv" -inchk "Output.csv" -outnames "Output-names.fm" -bp "SomeDirectory/"

The meanings of the command line arguments above are as follows.

  * `source`: the source directory to list the files from.
  * `alg`: the hash algorithm to use (can be `md5`, `sha1`, `sha256`, `sha512`).
  * `outchk`: the name of the output CSV containing checksums and some other meta data.
  * `inchk`: the name of the input CSV which is used to identify files and match them with their old name.
  * `bp`: base path, this will be trimmed from the path strings written in the meta files. Optional.

## Development Environment

  * Windows 10
  * Go 1.7.4 Windows amd64
  * Visual Studio Code 1.12.1
    * Extension: Go 0.6.61
