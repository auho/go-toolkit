
allModules	:= $(shell find . -not \( \
                       		\( \
							-path './output' \
							-o -path './.git' \
							-o -path '*/vendor/*' \
                       		\) -prune \
                       		\) -name 'go.mod' -print0 | xargs -0 -I {} dirname {})

.PHONY: help
help:
	@echo "help"
	@echo "  vet"
	@echo "  tidy"

.PHONY: vet
vet:
	$(shell for mod in ${allModules}; \
		do \
			pushd $${mod} > /dev/null; \
			go vet ./...; \
			if [ $$? -ne 0 ]; then echo "$${mod} vet error"; exit; fi; \
			popd > /dev/null || exit; \
		done \
	)
	@echo "go vet finished"

.PHONY: tidy
tidy:
	$(shell for mod in ${allModules}; \
		do \
			pushd $${mod} > /dev/null; \
			echo "go mod tidy $(sed -n 1p go.mod | cut -d ' ' -f2)"; \
            go mod tidy; \
			popd > /dev/null || exit; \
		done \
	)
	@echo "go tidy finished"
