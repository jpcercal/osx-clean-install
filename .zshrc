export PATH="/opt/homebrew/bin:~/dotfiles/bin:$PATH"

###########################################################
# flags

export PATH="$(brew --prefix gettext)/bin:$PATH"
export CPPFLAGS="-I$(brew --prefix gettext)/include $CPPFLAGS" 
export LDFLAGS="-L$(brew --prefix gettext)/lib $LDFLAGS"

export PATH="$(brew --prefix curl)/bin:$PATH"
export CPPFLAGS="-I$(brew --prefix curl)/include $CPPFLAGS" 
export LDFLAGS="-L$(brew --prefix curl)/lib $LDFLAGS"

export PATH="$(brew --prefix openssl)/bin:$PATH"
export CPPFLAGS="-I$(brew --prefix openssl)/include $CPPFLAGS" 
export LDFLAGS="-L$(brew --prefix openssl)/lib $LDFLAGS"

export CPPFLAGS="-I$(brew --prefix readline)/include $CPPFLAGS" 
export LDFLAGS="-L$(brew --prefix readline)/lib $LDFLAGS"

export PATH="$(brew --prefix sqlite)/bin:$PATH"
export CPPFLAGS="-I$(brew --prefix sqlite)/include $CPPFLAGS" 
export LDFLAGS="-L$(brew --prefix sqlite)/lib $LDFLAGS"

export PATH="$(brew --prefix ncurses)/bin:$PATH"
export CPPFLAGS="-I$(brew --prefix ncurses)/include $CPPFLAGS" 
export LDFLAGS="-L$(brew --prefix ncurses)/lib $LDFLAGS"

export PATH="$(brew --prefix ruby)/bin:$(ruby -e 'puts Gem.bindir'):$PATH"
export CPPFLAGS="-I$(brew --prefix ruby)/include $CPPFLAGS" 
export LDFLAGS="-L$(brew --prefix ruby)/lib $LDFLAGS"

###########################################################
# python

export PYENV_ROOT="$HOME/.pyenv"
command -v pyenv >/dev/null || export PATH="$PYENV_ROOT/bin:$PATH"
eval "$(pyenv init -)"

###########################################################
# go

export GOPATH="$HOME/go"
export PATH="$GOPATH/bin:$PATH"

###########################################################
# nvm 

export NVM_DIR="$HOME/.nvm"
  [ -s "/usr/local/opt/nvm/nvm.sh" ] && . "/usr/local/opt/nvm/nvm.sh"  # This loads nvm
  [ -s "/usr/local/opt/nvm/etc/bash_completion" ] && . "/usr/local/opt/nvm/etc/bash_completion"  # This loads nvm bash_completion

###########################################################
# Pre configuration

# Define the environment variable ZPLUG_HOME
export ZPLUG_HOME="/opt/homebrew/opt/zplug"

# Loads zplug
source $ZPLUG_HOME/init.zsh

# Clear packages
zplug clear

# Load fzf
[ -f ~/.fzf.zsh ] && source ~/.fzf.zsh

###########################################################
# Packages

ZSH_AUTOSUGGEST_HIGHLIGHT_STYLE="fg=#ff00ff,bg=cyan,bold,underline"
ZPLUG_HOOK_LOAD_WD="wd() { . $ZPLUG_REPOS/mfaerevaag/wd/wd.sh }"

autoload -Uz compinit
compinit

zplug "plugins/history", from:oh-my-zsh # https://github.com/ohmyzsh/ohmyzsh/blob/master/lib/history.zsh
zplug "wfxr/forgit", from:github
zplug "mfaerevaag/wd", defer:2, as:command, use:"wd.sh", hook-load:"$ZPLUG_HOOK_LOAD_WD"
zplug "peterhurford/up.zsh", defer:2, from:github
zplug "tysonwolker/iterm-tab-colors", defer:2, from:github
zplug "zsh-users/zsh-autosuggestions", defer:2, from:github
zplug "zsh-users/zsh-completions", defer:2, from:github
zplug "zsh-users/zsh-syntax-highlighting", defer:2, from:github

###########################################################
# Install packages that have not been installed yet

if ! zplug check --verbose; then
    printf "Install? [y/N]: "
    if read -q; then
        echo; zplug install
    else
        echo
    fi
fi

zplug load

###########################################################
# Post configuration

# Added by serverless binary installer
export PATH="$HOME/.serverless/bin:$PATH"

# iTerm
test -e "${HOME}/.iterm2_shell_integration.zsh" && source "${HOME}/.iterm2_shell_integration.zsh"

# ###########################################################
# # Aliases

alias l="exa -lbF"
alias ll="exa -lbhHigUmuSa@ --sort=modified"
alias lt="exa --tree"

alias -- -='cd -'
alias 1='up 1'
alias 2='up 2'
alias 3='up 3'
alias 4='up 4'
alias 5='up 5'
alias 6='up 6'
alias 7='up 7'
alias 8='up 8'
alias 9='up 9'

# ###########################################################
# # Custom ENV vars

export GOCACHE=/tmp 
export GOSUMDB=off

###########################################################
# Starts https://starship.rs

eval "$(starship init zsh)"
[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"  # This loads nvm bash_completion
