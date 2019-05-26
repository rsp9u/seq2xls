ifeq ($(GOOS), windows)
	BINARY = bin/seq2xls.exe
else
	BINARY = bin/seq2xls
endif

.PHONY: all
all: seq2xls

.PHONY: test
test: seqdiag
	cd seqdiag/test ; \
	go test -v ; \
	cd ../..
	cd seqdiag ; \
	go test -v ; \
	cd ..

seq2xls: seqdiag *.go cmd/main.go
	go build -o $(BINARY) cmd/main.go

.PHONY: seqdiag
seqdiag: gocc
	cd seqdiag ; \
	rm -rf errors lexer parser token util ; \
	gocc grammar.bnf ; \
	cd ..

.PHONY: clean
clean:
	rm -f seq2xls
	rm -rf seqdiag/errors seqdiag/lexer seqdiag/parser seqdiag/token seqdiag/util

.PHONY: gocc
gocc:
	if ! which gocc > /dev/null; then \
	  go get github.com/goccmack/gocc ; \
	fi
