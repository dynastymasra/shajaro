.DEFAULT_GOAL := build-docker

REPOSITORY ?= gcr.io/solid-terra-207000/actor
VERSION ?= develop
IMAGE ?= $(REPOSITORY):$(VERSION)

build-docker:
	docker build -t $(IMAGE) -f Dockerfile .

push-docker:
	gcloud docker -- push $(IMAGE)