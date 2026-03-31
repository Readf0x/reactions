# build settings
DEBUG    ?= false
PACKAGES ?= gtk4

PKG_FLAGS := $(addprefix --pkg ,$(PACKAGES))

# installation settings
PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin
MANDIR ?= $(PREFIX)/share/man

reactions: $(wildcard src/*.vala) makefile
	valac -X -w$(if $(filter true,$(DEBUG)), -g) \
		$(PKG_FLAGS) \
		--output=reactions \
		$(wildcard src/*.vala)

build: reactions compile_commands.json

run: build
	$(if $(filter true,$(DEBUG)),GTK_DEBUG=interactive )./reactions

install: build
	install -Dm755 reactions "$(DESTDIR)$(BINDIR)/reactions"

compile_commands.json: $(wildcard src/*.vala) makefile
	@echo '[' > $@
	$(foreach f,$(wildcard src/*.vala),\
		echo '  {' >> $@ && \
		echo '    "directory": "$(PWD)",' >> $@ && \
		echo '    "command": "valac $(PKG_FLAGS) --output=reactions $(f)",' >> $@ && \
		echo '    "file": "$(f)"' >> $@ && \
		echo '  },' >> $@)
	@sed -i '$$ s/,$$//' $@
	@echo ']' >> $@
