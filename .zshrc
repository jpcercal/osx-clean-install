
###########################################################
# homebrew

eval "$(/opt/homebrew/bin/brew shellenv)"

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
[[ -d $PYENV_ROOT/bin ]] && export PATH="$PYENV_ROOT/bin:$PATH"
eval "$(pyenv init -)"

###########################################################
# go

export GOPATH="$HOME/go"
export PATH="$GOPATH/bin:$PATH"

###########################################################
# rust

source "$HOME/.cargo/env"

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

autoload -Uz compinit
compinit

zplug "wfxr/forgit", from:github
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

###########################################################
# atuin (history)

atuin-setup() {
  ! hash atuin && return
  
  export ATUIN_NOBIND="true"
  eval "$(atuin init zsh)"
  fzf-atuin-history-widget() {
    local selected num
    setopt localoptions noglobsubst noposixbuiltins pipefail no_aliases 2>/dev/null
    selected=$(atuin search --cmd-only --limit ${ATUIN_LIMIT:-5000} | tac |
      FZF_DEFAULT_OPTS="--height ${FZF_TMUX_HEIGHT:-40%} $FZF_DEFAULT_OPTS -n2..,.. --tiebreak=index --bind=ctrl-r:toggle-sort,ctrl-z:ignore $FZF_CTRL_R_OPTS --query=${LBUFFER} +m" fzf)
    local ret=$?
    if [ -n "$selected" ]; then
      # the += lets it insert at current pos instead of replacing
      LBUFFER+="${selected}"
    fi
    zle reset-prompt
    return $ret
  }
  zle -N fzf-atuin-history-widget

  bindkey '^E' _atuin_search_widget
  bindkey '^R' fzf-atuin-history-widget
  bindkey '^[[A' fzf-atuin-history-widget
}
atuin-setup

bindkey '^ ' autosuggest-accept
_zsh_autosuggest_strategy_atuin_top() {
    suggestion=$(atuin search --cmd-only --limit 1 --search-mode prefix $1)
}
ZSH_AUTOSUGGEST_STRATEGY=atuin_top

###########################################################
# Kubectl

export KUBECONFIG=~/.kube/config
export AWS_PROFILE="TXBTurboshopDeveloperAccess-440308253360"

###########################################################
# functions

gli() {
  local filter
  if [ -n $@ ] && [ -f $@ ]; then
    filter="-- $@"
  fi

  git log \
    --graph --color=always --abbrev=7 --format='%C(auto)%h %an %C(blue)%s %C(yellow)%cr' $@ | \
    fzf \
      --ansi --no-sort --reverse --tiebreak=index \
      --preview "f() { set -- \$(echo -- \$@ | grep -o '[a-f0-9]\{7\}'); [ \$# -eq 0 ] || git show --color=always \$1 $filter; }; f {}" \
      --bind "j:down,k:up,alt-j:preview-down,alt-k:preview-up,ctrl-f:preview-page-down,ctrl-b:preview-page-up,q:abort,ctrl-m:execute:
                (grep -o '[a-f0-9]\{7\}' | head -1 |
                xargs -I % sh -c 'git show --color=always % | less -R') << 'FZF-EOF'
                {}
                FZF-EOF" \
      --preview-window=right:60% \
      --height 80%
}
