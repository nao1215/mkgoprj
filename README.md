# ubume - Golang project generator
ubume command generate golang project at current directory. Currently, ubume can only generate application projects. However, in the future ubume will also be able to generate library projects.
  
# How to install
## Step.1 Install golang
If you don't install golang in your system, please install Golang first. Check the [Go official website](https://go.dev/doc/install) for how to install golang.
## Step2. Install ubume
```
$ go install github.com/nao1215/ubume/cmd/ubume
```
  
# How to use
In the following example, the ubume command will generate a sample project. The binary name will be sample, and build using Makefile.
```
$ ubume github.com/nao1215/sample  ※ Argument is same as "$ go mod init"
$ tree sample/
sample/
├── Changelog.md
├── Makefile
├── cmd
│   └── sample
│       └── main.go
└── go.mod

$ cd sample
$ make build
$ ls
Changelog.md  Makefile  cmd  go.mod  sample

$ ./sample 
Hello, World
```