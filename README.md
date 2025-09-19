# sp

Unix style stream processing utility.

Your companion to make reading information in a shell more enjoyable.

## Installation

### Install the sp binary

For a unix-type OS the following string of commands will:
 - download the archive of the specified version for your os and machine architecture (to your working directory)
 - extract it inside your home directory
 - create a symbolic link to $HOME/bin
 - cleanup the downloaded archive

So if you have `$HOME/bin` on your `PATH` environment variable the `sp` command is resolvable.

```bash
VERSION="1.0.0" && \
  TARGET="sp_$(uname)_$(uname -m)" && \
  TARBAL="${TARGET}.tar.gz" && \
  curl -L "https://github.com/pvbouwel/sp/releases/download/v${VERSION}/${TARBAL}" -o "$TARBAL" && \
  mkdir -p "$HOME/${TARGET}" && \
  tar -C "$HOME/${TARGET}" -xzvf "$TARBAL" && \
  ln -svf "${HOME}/${TARGET}/sp" "${HOME}/bin/sp" && rm "$TARBAL"
```

### Install aliases

The sp tool allows for configuration to fine tune behavior to your liking. As a result there are likely too many flags to be convenient to type every time. To overcome this you can install aliases in your rc file.

To show the example aliases run:
```bash
sp aliases
```

The recommended way to install them is to put them in a separate file:
```bash
sp aliases > "$HOME/.sp-aliases"
```

And a once of install of the source command in your rc file:
```bash
echo 'source "$HOME/.sp-aliases"' >> ${HOME}/.$(echo $SHELL| sed 's_.*[/]\([a-z]*\)_\1_')rc
```

Tip: the comments at the start of the output will show similar commands tailored for your environment which look a little less daunting.



## Usage

The source of truth on how to use `sp` is reachable by invoking `sp --help`.

If you installed aliases you can use aliases in a similar manner. For example:
- `cat /tmp/myfile | sp-rainbow`
- `sp-rainbow -- cat /tmp/myfile`
