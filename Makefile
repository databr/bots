.PHONY: no_targets__ help

help:
	sh -c "$(MAKE) -p no_targets__ | awk -F':' '/^[a-zA-Z0-9][^\$$#\/\\t=]*:([^=]|$$)/ {split(\$$1,A,/ /);for(i in A)print A[i]}' | grep -v '__\$$' | sort"

no_targets__:

camarabook-data:
	go get
	go build

all: deputies
	echo "Finished"
	make clean

deputies: deputies_from_search deputies_from_xml deputies_about deputies_quotas deputies_info_from_transparencia_brasil

senators: senators-from-index

senators-from-index: camarabook-data
	./camarabook-data --save-senators-from-index

deputies_from_search: camarabook-data
	./camarabook-data --save-deputies-from-search

deputies_from_xml: camarabook-data
	./camarabook-data --save-deputies-from-xml

deputies_about: camarabook-data
	./camarabook-data --save-deputies-about

deputies_quotas: camarabook-data
	./camarabook-data --save-deputies-quotas

deputies_info_from_transparencia_brasil: camarabook-data
	./camarabook-data --save-deputies-info-from-transparencia-brasil

clean:
	rm -Rf ./camarabook-data
	rm -Rf pkg/*

build_all: clean pkg/camarabook-linux-amd64

pkg/camarabook-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o pkg/camarabook-linux-amd64
