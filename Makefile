.PHONY: docs
docs:
	@docker run -v $$PWD/:/docs pandoc/latex -f markdown /docs/README.md -o /docs/build/output/README.pdf