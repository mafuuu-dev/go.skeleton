.SILENT:
.ONESHELL:

THIS_FILE := $(lastword $(MAKEFILE_LIST))

.PHONY: help
help:
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

# -----------------------------------------------------------------------
# Configuration:
# -----------------------------------------------------------------------

d := COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker

dc := $(d) compose $(args)
de := $(d) exec -ti

# -----------------------------------------------------------------------
# Common:
# -----------------------------------------------------------------------

# $1 - message
define common.info
	echo "\033[1m==> $1\033[0m"
endef

# $1 - command; $2 - subject
define common.run
	@$(call common.info,$2)
	make $1 || true
endef

# $1 - path, $2 - command
define common.runFrom
	make -C $1 $2 || true
endef

# -----------------------------------------------------------------------
# DockerCompose:
# -----------------------------------------------------------------------

# $1 - service; $2 - command
define compose.exec
	$(de) $(project)-$1-1 $2
endef

# $1 - command; $2 - params; $3 - services
define compose.use
	$(dc) $1 $2 $3
endef
