# tttns

`tttns` is a command-line interface (CLI) tool designed to inspect 3GPP TS 32.297 CDR (Charging Data Record) files. The name "tttns" is derived from the first letters of "32297".

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

## Usage

```bash
tttns [file|-] [flags]
tttns [command]
```

## Available Commands

- `cdr`: Print all CDR header info or use its sub commands
- `completion`: Generate the autocompletion script for the specified shell
- `file`: Print CDR file header info or use its sub commands
- `help`: Help about any command
- `version`: Print the version number of tttns

## Flags

- `-h, --help`: Display help for tttns
- `-j, --json`: Output in JSON format (except for cdr dump)

## Detailed Command Information

### cdr Command

The `cdr` command is used to print and manipulate CDR (Charging Data Record) information.

Usage:

```bash
tttns cdr [file|-] [flags]
tttns cdr [command]
```

Flags:

- `-h, --help`: Display help for the cdr command
- `-j, --json`: Output in JSON format

#### Available Subcommands

1. **count**: Get the number of CDRs in a file

   ```bash
   tttns cdr count [file|-]
   ```

2. **dump**: Dump the raw content of CDR to stdout

   ```bash
   tttns cdr dump [file|-] [index|1]
   ```

3. **header**: Print CDR header info

   ```bash
   tttns cdr header [file|-] [index|1]
   ```

For more information about a specific subcommand, use:

```bash
tttns cdr [subcommand] --help
```

## Examples

1. Get the number of CDRs in a file:

   ```bash
   tttns cdr count example.cdr
   ```

2. Print CDR header info of the 1st CDR:

   ```bash
   tttns cdr header example.cdr 1
   cat example.cdr | tttns cdr header 1
   ```

3. Dump the raw content of the 2nd CDR to stdout:

   ```bash
   tttns cdr dump example.cdr 2
   cat example.cdr | tttns cdr dump 2
   ```

4. Print CDR header info in JSON format:
  
   ```bash
   tttns cdr header example.cdr 1 --json
   ```

## License

Apache-2.0.

## Notes

The project has been scaffolded with the help of [kleiner](https://github.com/can3p/kleiner).
