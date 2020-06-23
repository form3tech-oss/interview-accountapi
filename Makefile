.PHONY: docs
docs:
	docker run -v ${PWD}/:/data pandoc/latex -f markdown README.md -o README.pdf && open README.pdf
