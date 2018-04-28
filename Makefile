# Install has to come before lint, otherwise the linter fails due to outside
# # dependancies
.PHONY: testutils
testutils:
	go build -o cmd/testutils/runtestapp \
		github.com/stiganik/gap/cmd/testutils/runtest
	go build -o cmd/testutils/selectortestapp \
		github.com/stiganik/gap/cmd/testutils/selectortest
	go build -o cmd/testutils/combinationtestapp \
		github.com/stiganik/gap/cmd/testutils/combinationtest

.PHONY: install
install:
	go install ./...

.PHONY: lint
lint: install
	if which gometalinter > /dev/null; then \
		env TABWIDTH=8 gometalinter \
			--enable-all --disable misspell --disable test --disable testify \
			--disable safesql --disable nakedret --cyclo-over=25 \
			--line-length=120 --deadline=5m --enable-gc -t; \
	fi

