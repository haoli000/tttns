# tttns

A command-line tool to inspect 3GPP TS 32.297 CDR (Charging Data Record) files. The name is derived from the first letters of "32297".

To decode the dumped CDR contents, the companion tool [xchf](https://github.com/haoli000/xchf) can be used.

## Installation

### Install Script

```bash
curl -L https://raw.githubusercontent.com/haoli000/tttns/main/generated/install.sh | bash
```

Specific version:

```bash
curl -L https://raw.githubusercontent.com/haoli000/tttns/main/generated/install.sh | bash -s 0.2.4
```

The script installs into `$HOME/bin` by default. Override with `$CUSTOM_INSTALL`.

### Manual Download

Get the archive for your platform from [Releases](https://github.com/haoli000/tttns/releases) and place the binary in your `$PATH`.

### Build from Source

Requires Go 1.26+:

```bash
git clone https://github.com/haoli000/tttns.git
cd tttns
make build
```

## Usage

```bash
tttns [file|-] [flags]
tttns [command]
```

All commands accept either a filename or stdin (via pipe or `-`). Output defaults to YAML; use `-j` for JSON.

## Commands

| Command | Description |
|---------|-------------|
| `tttns [file]` | Print file header + all CDR headers |
| `tttns file [file]` | Print file header only |
| `tttns cdr [file]` | Print all CDR headers |
| `tttns cdr count [file]` | Print number of CDRs |
| `tttns cdr header [file] [index]` | Print a specific CDR header (default: 1) |
| `tttns cdr dump [file] [index]` | Dump raw CDR content to stdout (default: 1) |
| `tttns version` | Print version |

For `header` and `dump`, if a single argument is provided it is treated as an index (reading from stdin) if numeric, or as a filename (with index defaulting to 1) otherwise.

## Flags

- `-j, --json`: Output in JSON format (does not apply to `cdr dump`)
- `-h, --help`: Display help

## Examples

Print file summary:

```bash
tttns example.cdr
tttns -j example.cdr
```

Count CDRs:

```bash
tttns cdr count example.cdr
```

Inspect a specific CDR header:

```bash
tttns cdr header example.cdr 3
cat example.cdr | tttns cdr header 3
```

Dump raw CDR content (e.g., pipe to [xchf](https://github.com/haoli000/xchf) for BER decoding):

```bash
tttns cdr dump example.cdr 2 | xchf
cat example.cdr | tttns cdr dump 2 | xchf
```

JSON output:

```bash
tttns cdr header -j example.cdr 1
```

## Sample Output

```yaml
header_info:
  file_length: 49982
  header_length: 54
  high_release_version: "15.10"
  low_release_version: "15.10"
  file_opening_timestamp: 18/5 03:14:00+0000
  last_cdr_append_timestamp: 18/5 03:14:00+0000
  number_of_cdrs_in_file: 100
  file_sequence_number: 2808926
  file_closure_trigger_reason: 0 - Normal closure (Undefined normal closure reason)
  node_ip_address: 10.244.190.107
  lost_cdr_indicator: No CDRs have been lost
cdr_info:
  number_of_cdrs: 100
  cdr_headers:
    - cdr_length: 501
      release_version: "15.10"
      data_record_format: BER
      ts_number: TS32.274
```

## Development

```bash
make check    # fmt + vet + staticcheck + test + govulncheck
make build    # build binary
make test     # run tests
make lint     # fmt + vet + staticcheck
```

## License

Apache-2.0.

## Acknowledgements

Scaffolded with [kleiner](https://github.com/can3p/kleiner).
