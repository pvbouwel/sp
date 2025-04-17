# sp

 Unix style stream processing utility

## Installation

### Install the sp binary

### Install aliases

The sp tool allows for configuration to fine tune behavior to your liking. As a result there are likely too many flags to be convenient to type every time. To overcome this you can install aliases in your rc file (e.g. `~/.zshrc` or `~/.bashrc` or another based on your shell)

Some examples:

```
alias sp-rainbow="sp color --color-type rotating --rotating-type random --stride-length 5-15"
alias sp-json-traffic="sp color --color-type JSON --colors level.info --colors info.0.255.0,warning.255.128.0,error.255.0.0"
alias sp-stdouterr="sp color"
alias sp-epoch="sp epoch"
```

## Usage

The source of truth on how to use `sp` is reachable by invoking `sp --help`.

If you installed aliases you can use aliases in a similar manner. For example:
- `cat /tmp/myfile | sp-rainbow`
- `sp-rainbow -- cat /tmp/myfile`
