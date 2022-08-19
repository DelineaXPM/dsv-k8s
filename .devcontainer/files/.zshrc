#eval "$(direnv hook zsh)"
export ZSH="$HOME/.oh-my-zsh"
ZSH_THEME=avit
plugins=(z direnv zsh-interactive-cd docker golang gh zsh-navigation-tools)
source "$ZSH/oh-my-zsh.sh"
