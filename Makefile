.PHONY: no_targets__ help

help:
	sh -c "$(MAKE) -p no_targets__ | awk -F':' '/^[a-zA-Z0-9][^\$$#\/\\t=]*:([^=]|$$)/ {split(\$$1,A,/ /);for(i in A)print A[i]}' | grep -v '__\$$' | sort"

no_targets__:

pkg/parliamentarian_bot:
	cd go_bot/parliamentarian_bot && GOOS=linux GOARCH=amd64 go build -o ../../pkg/parliamentarian_bot

pkg/metrosp_bot:
	cd go_bot/metrosp_bot && GOOS=linux GOARCH=amd64 go build -o ../../pkg/metrosp_bot

clean:
	rm -Rf pkg/*

parliamentarian_bot: clean pkg/parliamentarian_bot

metrosp_bot: clean pkg/metrosp_bot

deploy_go: parliamentarian_bot metrosp_bot
	rsync -Pavh pkg/parliamentarian_bot $(DATABR_BOT_MACHINE):/usr/local/bin/parliamentarian_bot
	rsync -Pavh pkg/metrosp_bot $(DATABR_BOT_MACHINE):/usr/local/bin/metrosp_bot

deploy: deploy_go
