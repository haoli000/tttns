# tttns - A cli tool to inspect 3GPP TS 32.297 CDR files

Based on the [TS 32.297 CDR Filer Parser](https://github.com/haoli000/TS32.297_CDR_File_Parser)
and [Kaitai Struct: compiler](https://github.com/kaitai-io/kaitai_struct_compiler) to generate the parser code.

To get started, run:
```bash
> tttns -h
A cli tool to inspect 3GPP TS 32.297 CDR files. 
The name is simply derived from the first letters of 32297.

Usage:
  tttns [file|-] [flags]
  tttns [command]

Available Commands:
  cdr         Print all CDR header info
  completion  Generate the autocompletion script for the specified shell
  file        Print CDR file header info
  help        Help about any command
  version     Print the version number of tttns

Flags:
  -h, --help   help for tttns
  -j, --json   Output in JSON format except for cdr dump

Use "tttns [command] --help" for more information about a command.
```

To get CDR info, run:
```bash
> tttns cdr -h
Print all CDR header infos

Usage:
  tttns cdr [file|-] [flags]
  tttns cdr [command]

Available Commands:
  count       Get number of CDRs file
  dump        Dump the row content of CDR to stdout
  header      Print CDR header info

Flags:
  -h, --help   help for cdr
  -j, --json   Output in JSON format

Use "tttns cdr [command] --help" for more information about a command.
```

## Installation

### Install Script

Download `tttns` and install into a local bin directory.

#### MacOS, Linux, WSL

Latest version:

```bash
curl -L https://raw.githubusercontent.com/haoli000/tttns/main/generated/install.sh | bash
```

Specific version:

```bash
curl -L https://raw.githubusercontent.com/haoli000/tttns/main/generated/install.sh | bash -s 0.0.4
```

The script will install the binary into `$HOME/bin` folder by default, you can override this by setting
`$CUSTOM_INSTALL` environment variable

### Manual download

Get the archive that fits your system from the [Releases](https://github.com/haoli000/tttns/releases) page and
extract the binary into a folder that is mentioned in your `$PATH` variable.

## Notes

The project has been scaffolded with the help of [kleiner](https://github.com/can3p/kleiner)
