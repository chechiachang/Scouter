DOCKERHUB_USER="chechiachang"

# Build
#

get:
	go get -v ./...

.PHONY: build
build: get
	go build ./...

db:
	docker run -d --name mongo -p 27017:27017 mongo 

migrate:
	docker cp data/mongodb/scouter mongo:.
	docker exec -it mongo bash -c "mongorestore -d scouter scouter"

# Test & Run
#

PYTHON = $(shell which python)

.PHONY: test
test:
	go test ./...

.PHONY: apiserver
apiserver:
	${PYTHON} ./face_recognition/apiserver.py

# Build & ship
#

.PHONY: encodings
encoding:
	rm -f face_recognition/encodings face_recognition/index
	${PYTHON} ./face_recognition/encoding_file_generator.py

.PHONY: base
base:
	time docker build \
    --tag $(DOCKERHUB_USER)/scouter-apiserver-base \
		--file face_recognition/base/Dockerfile .

.PHONY: image
image:
	time docker build \
    --tag $(DOCKERHUB_USER)/scouter-apiserver \
		--file face_recognition/Dockerfile .

.PHONY: unity
add-unity:
	cp -rf /Users/Shared/Unity/scouter/Assets/Scouter unity/Assets/

unity:
	cp -rf unity/Assets/Scouter /Users/Shared/Unity/scouter/Assets

clean:
	rm -f avatar_downloader contribution_fetcher user_detail_fetcher user_fetcher 

# Prerequisite
#

UNAME := $(shell uname)
PORT := $(shell which port)
BREW := $(shell which brew)

ifeq ($(UNAME), Linux)

prerequisite:
	sudo apt-get update && \
	sudo apt-get install -y python3 python3-pip

else ifeq ($(UNAME), Darwin)

prerequisite:
	if [[ "$(PORT)" != "" ]]; then sudo port install coreutils; fi
	if [[ "$(BREW)" != "" ]]; then brew install coreutils; fi
	rehash

endif
