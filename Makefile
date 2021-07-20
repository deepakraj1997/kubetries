.PHONY:	all unit smoke e2e

all: install-deps unit e2e

GODEPS	= \
	github.com/onsi/ginkgo/ginkgo \
	github.com/onsi/gomega/...
		

install-deps:
	for dep in $(GODEPS); do \
		go get $$dep; \
	done

unit: install-deps
	echo "Unit Tests" # for utils

smoke e2e: install-deps
	ginkgo e2e_tests/
	