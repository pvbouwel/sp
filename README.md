# sp

 Unix style stream processing utility

## Installation

### Install the sp binary

For a unix-type OS the following string of commands will:
 - download the archive of the specified version for your os and machine architecture (to your working directory)
 - extract it inside your home directory
 - create a symbolic link to $HOME/bin
 - cleanup the downloaded archive

So if you have `$HOME/bin` on your `PATH` environment variable the `sp` command is resolvable.

```bash
VERSION="0.1.0" && \
  TARGET="sp_$(uname)_$(uname -m)" && \
  TARBAL="${TARGET}.tar.gz" && \
  curl -L "https://github.com/pvbouwel/sp/releases/download/v${VERSION}/${TARBAL}" -o "$TARBAL" && \
  mkdir -p "$HOME/${TARGET}" && \
  tar -C "$HOME/${TARGET}" -xzvf "$TARBAL" && \
  ln -svf "${HOME}/${TARGET}/sp" "${HOME}/bin/sp" && rm "$TARBAL"
```

### Install aliases

The sp tool allows for configuration to fine tune behavior to your liking. As a result there are likely too many flags to be convenient to type every time. To overcome this you can install aliases in your rc file (e.g. `~/.zshrc` or `~/.bashrc` or another based on your shell)

Some examples:

```bash
# Rainbow colours
alias sp-rainbow="sp color --color-type rotating --rotating-type random --stride-length 15-25"

# Colour JSON depending on values of the field called levelname and have alternating colours if subsequent lines match
alias sp-json-traffic='sp color --ignore-case --color-type JSON --json-key levelname --colors INFO.0.255.0,INFO.0.155.0,WARNING.255.128.0,WARNING.155.128.0,ERROR.255.0.0,ERROR.155.0.0'

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
