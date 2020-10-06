build:
	docker build . -t haiken -f Dockerfile.prod

build-no-cache:
	docker build . -t haiken -f Dockerfile.prod --no-cache

ikku-test:
	docker run --rm -it -v ${CURDIR}:/app haiken go test -v ./...
