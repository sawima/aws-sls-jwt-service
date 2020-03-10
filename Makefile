.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/jwt-auth functions/jwtauth/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/jwt-verify functions/jwtverify/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/defaultapp functions/defaultapp/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/api-proxy functions/apiproxy/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/jwtaccount functions/jwtaccount/main.go


# help:
# 	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# functions := $(shell find functions -name \*main.go | awk -F'/' '{print $$2}')


# build: gomodgen
# 	export GO111MODULE=on
# 	@for function in $(functions) ; do \
# 		env GOOS=linux go build -ldflags="-s -w" -o bin/$$function functions/$$function/main.go ; \
# 	done


clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
