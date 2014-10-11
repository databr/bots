# Bots

> [Trello com Tarefas](https://trello.com/b/3WLlqXpX/databr)

O projeto **Bots** é responsável pela coletar, tratar e relacionar os dados das diversas fontes.

Inicialmente o projeto foi escrito em Ruby, foi reescrito algumas coisas em golang, porém não existe limitações de linguagens para esse projeto. Dado que algumas premissas sejam respeitas e tudo funciona:

## Desenvolvendo um bot
``` shell
// Example
$ git clone git@github.com:databr/bots.git bots
$ cd bots
$ git branch -b dukex-nome_do_meu_bot master
```

Primeiramente faça um fork do projeto [databr/bots](https://github.com/databr/bots/), crie um branch relacionado ao seu bot ou a alteração que deseja fazer, e depois nos envie um pull-request




<br />

<br />

### Acesso ao banco de dados

``` shell
$ echo $MONGO_URL
mongodb://localhost
$ echo $DATABASE_NAME
databr
```

Ao desenvolver seu bot tenha em mente que utilizamos **mongodb** e o acesso ao mongo de ser feito via a variavel de ambiente `$MONGO_URL`, a banco de dados no mongo utilizado é feito via a variavel de ambiente `$DATABASE_NAME`, a estrutura do banco depende muito do bot e do dado, bots que atualiza ou agrega dados existente devem sempre seguir a estrutura existente, qualquer duvida consulte sempre as [issues no github](https://github.com/databr/bots/issues), sinta-se livre pra enviar sugestões ou duvidas por lá.


### Criando bot

``` shell
$ mkdir js_bot && cd js_bot
$ touch package.json
$ touch runner.js
```

Por conviniencia é recomendado criar um diretorio \_bot com a linguagem como prefixo, exemplos: `ruby_bot`, `go_bot`, `python_bot`, etc, na raiz do projeto, a forma que o bot é feito é criterio do desenvolvedor do mesmo escolher, caso já exista bots na linguagem que você pretende contrubuir é altamente recomendado seguir os padrões já estabelecidos, caso contrario um padrão podera ser definido, altamente recomendado que documente o mesmo para futuros desenvolvedores.


### Bot Deploy

``` Makefile
// Golang Example

pkg/go-bot-linux-amd64:
	cd go_bot && GOOS=linux GOARCH=amd64 go build -o ../pkg/go-bot-linux-amd64

go_bots: clean pkg/go-bot-linux-amd64

deploy_go: go_bots
	rsync -Pavh pkg/go-bot-linux-amd64 $(DATABR_BOT_MACHINE):/usr/local/bin/go-bot
	ssh $(DATABR_BOT_MACHINE) 'restart go-bot'
```

O deploy é feito atraves do comando ```make deploy```, para cada linguagem o comando deploy tem a responsabilidade de criar binarios/executaveis e subir isso para o servidor de bots em produção, caso precise o servidor de bots em produção está armazenado na variavel de ambiente ```$DATABR_BOT_MACHINE``` e o binarios/executaveis deve ser enviado para ```/usr/local/bin/``` nesse servidor, o exemplo ao lado é forma que está sendo feito com golang, um binario é criado no diretorio ```pkg``` e enviado via rsync para o server, bem fácil e simples.


