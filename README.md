# Bots

Join us on IRC at #databr to chat with other databr maintainers! ([web access](http://webchat.freenode.net/?channels=databr))

O projeto **Bots** é responsável pela coletar, tratar e relacionar os dados das diversas fontes.

Inicialmente o projeto foi escrito em Ruby, foi reescrito algumas coisas em golang, porém não existe limitações de linguagens para esse projeto. Dado que algumas premissas sejam respeitas e tudo funciona:

## Desenvolvendo um bot

Primeiramente faça um fork desse projeto, crie um branch relacionado ao seu bot ou a alteração que deseja fazer, e depois nos envie um pull-request

``` shell
// Example
$ git clone git@github.com:databr/bots.git bots
$ cd bots
$ git branch -b dukex-nome_do_meu_bot master
```

### Acesso ao banco de dados

Ao desenvolver seu bot tenha em mente que utilizamos **mongodb** e o acesso ao mongo de ser feito via a variavel de ambiente `$MONGO_URL`, a banco de dados no mongo utilizado é feito via a variavel de ambiente `$MONGO_DATABASE_NAME`, a estrutura do banco depende muito do bot e do dado, bots que atualiza ou agrega dados existente devem sempre seguir a estrutura existente, qualquer duvida consulte sempre as [issues no github](https://github.com/databr/bots/issues), sinta-se livre pra enviar sugestões ou duvidas por lá. 

``` shell
$ echo $MONGO_URL
mongodb://localhost
$ echo $MONGO_DATABASE_NAME
databr
```

Caso você esteja utilizando **golang** use os packages [github.com/databr/api/database](https://github.com/databr/api/tree/master/database) e [github.com/databr/api/models](https://github.com/databr/api/tree/master/models), database tem as conexões do banco prontas e models tem alguns modelos já usados.

``` golang
package bot

import (
  "github.com/databr/api/database"
  "github.com/databr/api/models"
)

func Start(){ 
  DB := database.NewMongoDB()
  // ...
  DB.Create(models.State{ 
    // ..
  })
}
```

### Criando bot

``` shell
$ mkdir js_bot && cd js_bot
$ touch package.json
$ touch runner.js
```

Por conviniencia é recomendado criar um diretorio \_bot com a linguagem como prefixo, exemplos: `ruby_bot`, `go_bot`, `python_bot`, etc, na raiz do projeto, a forma que o bot é criterio do desenvolvedor do mesmo escolher, caso já exista bots na linguagem que você pretende contrubuir é altamente recomendado seguir os padrões já estabelecidos, caso contrario um padrão podera ser definido, é altamente recomendado que documente o mesmo para futuros desenvolvedores.


### Bot Deploy

O deploy é feito atraves do comando ```make deploy```, para cada linguagem o comando deploy tem a responsabilidade de criar binarios/executaveis e subir isso para o servidor de bots em produção, caso precise do servidor de bots em produção está armazenado na variavel de ambiente ```$DATABR_BOT_MACHINE``` e o binarios/executaveis deve ser enviado para ```/usr/local/bin/``` nesse servidor, o exemplo abaixo é a forma que está sendo feito com golang, um binario é criado no diretorio ```pkg``` e enviado via rsync para o server, bem fácil e simples.

``` Makefile
// Golang Example

pkg/parliamentarian_bot:
	cd go_bot/parliamentarian_bot && go build -o ../../pkg/parliamentarian_bot
pkg/metrosp_bot:
	cd go_bot/metrosp_bot && go build -o ../../pkg/metrosp_bot
pkg/ibge_bot:
	cd go_bot/ibge_bot && go build -o ../../pkg/ibge_bot

deploy_go: build_go
	rsync -Pavh pkg/parliamentarian_bot $(DATABR_BOT_MACHINE):/usr/local/bin/parliamentarian_bot
	rsync -Pavh pkg/metrosp_bot $(DATABR_BOT_MACHINE):/usr/local/bin/metrosp_bot
	rsync -Pavh pkg/ibge_bot $(DATABR_BOT_MACHINE):/usr/local/bin/ibge_bot

```

### Servidor de Bots

Os bots rodam em [bots.databr.io](http://bots.databr.io).
