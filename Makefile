PROG=rkt-inspect
SOURCEDIR=.

SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

$(PROG): $(SOURCES)
	godep go build -o $@

.PHONY: clean
clean:
	rm -f $(PROG)
