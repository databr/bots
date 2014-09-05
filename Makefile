.PHONY: no_targets__ help

help:
	sh -c "$(MAKE) -p no_targets__ | awk -F':' '/^[a-zA-Z0-9][^\$$#\/\\t=]*:([^=]|$$)/ {split(\$$1,A,/ /);for(i in A)print A[i]}' | grep -v '__\$$' | sort"

no_targets__:

pkg/go-bot-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o pkg/go-bot-linux-amd64

clean:
	rm -Rf pkg/*

go_bots: clean pkg/go-bot-linux-amd64

deploy_go: go_bots
	rsync -Pavh pkg/go-bot-linux-amd64 $(DATABR_BOT_MACHINE):/usr/local/bin/go-bot
	ssh $(DATABR_BOT_MACHINE) 'restart go-bot'

deploy: deploy_go


