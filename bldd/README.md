# bldd (Backward ldd)

A tool that shows all EXECUTABLE files that use specified shared library files. It supports multiple architectures and generates detailed reports.

## Features

- Scans directories for ELF executables
- Supports multiple architectures:
  - x86 (i386)
  - x86-64
  - armv7
  - aarch64
- Generates sorted reports by number of executable usages
- Configurable scan directory
- Output reports in text format

## Usage

```bash
bldd [options]
```

### Options

- `-dir string`: Directory to scan for ELF files (default: ".")
- `-output string`: Output file for the report (default: "report.txt")

### Examples

1. Scan current directory and generate report:

```bash
bldd
```

2. Scan specific directory and save report to custom file:

```bash
bldd -dir /path/to/dir -output my_report.txt
```

3. Scan home directory:

```bash
bldd -dir $HOME -output home_report.txt
```

## Output Format

The report is organized by architecture and shows:

- Library name and number of executables using it
- List of executables using each library
- Sorted by number of usages (high to low)

Example output:

```txt
Report on dynamic used libraries by ELF executables on /home
==================================================

---------- i386 (x86) ----------
libc.so.6 (1 execs)
-> /path/to/executable1
-> /path/to/executable2

---------- x86-64 ----------
libc.so.6 (681 execs)
-> /path/to/executable1
-> /path/to/executable2
...
```

## Requirements

- Go 1.21 or later
- Linux system (for ELF file support)
