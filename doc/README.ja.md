[![Build](https://github.com/nao1215/ubume/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/nao1215/ubume/actions/workflows/build.yml)  
# ubume - Golngプロジェクトテンプレートジェネレータ
ubumeコマンドは、golangプロジェクトテンプレートをカレントディレクトリに作成します。Version 1.0.0では、アプリケーションプロジェクトとライブラリプロジェクトが作成できます。自動生成するファイルには、「プロジェクト管理を簡単にするMakefile」と「GitHub Actionsのファイル（ビルド、ユニットテスト）」が含まれます。ただし、"$ git init"は実行しません。
![Screenshot](./images/sample.png) 
  
# インストール方法
## Step.1 Golangのインストール
Golangをシステムにインストールしていない場合は、まずはgolangをインストールしてください。インストール方法は、[Go公式サイト](https://go.dev/doc/install) で確認してください。  
  
## Step2. ubumeのインストール
```
$ go install github.com/nao1215/ubume/cmd/ubume@latest
```
  
# 使い方
## アプリケーションプロジェクトの作成
以下の例では、ubumeコマンドはsampleプロジェクトを作成します。バイナリ名は"sample"で、ビルドにはMakefileを使います。
```
$ ubume github.com/nao1215/sample  ※ 引数は"$ go mod init"と同じ。
ubume starts creating the 'sample' application project (import path='github.com/nao1215/sample')

[START] check if ubume can create the project
[START] create directories
[START] create files
        sample (your project root)
         ├─ Makefile
         ├─ Changelog.md
         ├─ cmd
         │  └─ sample
         │     ├─ main.go
         │     ├─ main_test.go
         │     └─ doc.go
         └─ .github
            └─ workflows
               ├─ build.yml
               └─ unit_test.yml

BUILD SUCCESSFUL in 6[ms]

$ cd sample
$ make build
$ ls
Changelog.md  Makefile  cmd  go.mod  sample

$ ./sample 
Hello, World

$ make test
env GOOS=linux go test -v -cover ./... -coverprofile=cover.out
=== RUN   TestHelloWorld
--- PASS: TestHelloWorld (0.00s)
PASS
coverage: 50.0% of statements
ok      github.com/nao1215/sample/cmd/sample    0.001s  coverage: 50.0% of statements
go tool cover -html=cover.out -o cover.html
```

## ライブラリプロジェクトの作成
```
$ ubume --library github.com/nao1215/sample
ubume starts creating the 'sample' library project (import path='github.com/nao1215/sample')

[START] check if ubume can create the project
[START] create directories
[START] create files
        sample (your project root)
         ├─ sample_test.go
         ├─ Makefile
         ├─ Changelog.md
         ├─ doc.go
         ├─ sample.go
         └─ .github
            └─ workflows
               └─ unit_test.yml

BUILD SUCCESSFUL in 6[ms]
```

# 自己文書化されたMakefile
ubumeコマンドによって生成されるMakefileは、[自己文書化](https://postd.cc/auto-documented-makefile/)されています。makeコマンドを実行した時、Makefileのターゲットリストが表示されます。ターゲット名の横には、ヘルプメッセージが表示されます。

```
$ make
build           Build binary 
clean           Clean project
fmt             Format go source code 
test            Start test
vet             Start go vet
```
新しいターゲットを追加したい場合は、ターゲットの横に**"##"**から始まるコメントを書いてください。"##"以降の文字列が抽出され、ヘルプメッセージとして利用されます。以下に例を示します。
```
build:  ## Build binary 
	env GO111MODULE=on GOOS=$(GOOS) $(GO_BUILD) $(GO_LDFLAGS) -o $(APP) cmd/sample/main.go

clean: ## Clean project
	-rm -rf $(APP) cover.out cover.html
```
# 連絡先
「バグを見つけた場合」や「機能追加要望」に関するコメントを開発者に送りたい場合は、以下の連絡先を使用してください。

- [GitHub Issue](https://github.com/nao1215/ubume/issues)

# ライセンス
ubumeプロジェクトは、[Apache License 2.0](./LICENSE)条文の下でライセンスされています。