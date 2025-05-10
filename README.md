# sp

 Unix style stream processing utility

## Installation

### Install the sp binary

### Install aliases

The sp tool allows for configuration to fine tune behavior to your liking. As a result there are likely too many flags to be convenient to type every time. To overcome this you can install aliases in your rc file (e.g. `~/.zshrc` or `~/.bashrc` or another based on your shell)

Some examples:

```
# Rainbow colours
alias sp-rainbow="sp color --color-type rotating --rotating-type random --stride-length 15-25"

# Colour JSON depending on values of the field called levelname and have alternating colours if subsequent lines match
sp-json-traffic='sp color --color-type JSON --json-key levelname --colors INFO.0.255.0,INFO.0.155.0,WARNING.255.128.0,WARNING.155.128.0,ERROR.255.0.0,ERROR.155.0.0'

# Allow colouring of stoud and stderr differently
alias sp-stdouterr="sp color"

# Replace epoch occurrences with human readable time
alias sp-epoch="sp epoch"
```

## Usage

The source of truth on how to use `sp` is reachable by invoking `sp --help`.

If you installed aliases you can use aliases in a similar manner. For example:
- `cat /tmp/myfile | sp-rainbow`
- `sp-rainbow -- cat /tmp/myfile`
