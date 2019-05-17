.PHONY: all
all: gocc seqdiag

.PHONY: test
test: gocc seqdiag
	cd seqdiag/test ; \
	go test -v ; \
	cd ../..

.PHONY: gocc
gocc:
	if ! which gocc > /dev/null; then \
	  go get github.com/goccmack/gocc ; \
	fi

.PHONY: seqdiag
seqdiag:
	cd seqdiag ; \
	rm -rf errors lexer parser token util ; \
	gocc grammar.bnf ; \
	cd ..
