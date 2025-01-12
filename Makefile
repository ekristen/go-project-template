initialise:
	pre-commit --version || brew install pre-commit
	pre-commit install --install-hooks
	pre-commit run -a

build:
	goreleaser --clean --snapshot --skip sign

docs-build:
	docker run --rm -it -p 8000:8000 -v ${PWD}:/docs squidfunk/mkdocs-material build

docs-serve:
	docker run --rm -it -p 8000:8000 -v ${PWD}:/docs squidfunk/mkdocs-material

docs-seed:
	cp README.md docs/index.md
