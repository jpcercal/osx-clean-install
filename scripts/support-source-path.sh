#!/usr/bin/env bash

export PATH="/opt/homebrew/bin:~/dotfiles/bin:$PATH"
export PATH="$(brew --prefix ruby)/bin:$(ruby -e 'puts Gem.bindir'):$PATH"
export PATH="$PYENV_ROOT/bin:$PATH"
export PATH="$GOPATH/bin:$PATH"
export PATH="$HOME/.serverless/bin:$PATH"
