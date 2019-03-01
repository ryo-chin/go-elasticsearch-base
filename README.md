# Setup
1. install goenv
goenv is version management tool.
```bash
$ brew install goenv
$ echo 'export GOENV_ROOT="$HOME/.goenv"' >> ~/.bash_profile # or ~/.zshenv
$ echo 'export PATH="$GOENV_ROOT/bin:$PATH"' >> ~/.bash_profile # or ~/.zshenv
$ echo 'eval "$(goenv init -)"' >> ~/.bash_profile # or ~/.zprofile
$ source ~/.bash_profile
# if you use zsh
$ source ~/.zshenv
$ source ~/.zprofile
```
2. install dep
dep is package management tool.
```bash
$ brew install dep
# [NOTICE] GOPATH="$HOME/{anywhere ok}
$ echo 'export GOPATH="$HOME/works/go"' >> ~/.bash_profile # or ~/.zshenv
# [NOTICE] You have to make src dir under GOPATH.
$ mkdir ~/works/go/src
```
3. launch docker container
```bash
$ cd docker
$ docker-compose up -d
$ docker-compose start
```
4. run app
```bash
$ cd {project_root}
$ go run main.go
```
