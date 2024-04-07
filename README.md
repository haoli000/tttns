# tttns - A cli tool to inspect 3GPP TS 32.297 CDR files

Based on the [TS 32.297 CDR Filer Parser](https://github.com/haoli000/TS32.297_CDR_File_Parser)
and [Kaitai Struct: compiler](https://github.com/kaitai-io/kaitai_struct_compiler) to generate the parser code.

To get started, run:
```bash
tttns -h
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
