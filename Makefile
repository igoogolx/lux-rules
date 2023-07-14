define download =
.distfiles/$(1):
	mkdir -p .distfiles
	if ! curl -L#o $$@.unverified $(2); then rm -f $$@.unverified; exit 1; fi
	if ! echo "$(3)  $$@.unverified" | sha256sum -c; then rm -f $$@.unverified; exit 1; fi
	if ! mv $$@.unverified $$@; then rm -f $$@.unverified; exit 1; fi
endef

$(eval $(call download,geoip.dat,https://github.com/Loyalsoldier/v2ray-rules-dat/releases/download/202307132210/geoip.dat,4a05d4bce4cc7e9ab7a453add9d44e81e62313e3a055a6db4d0ff5e734448a5c))
$(eval $(call download,geosite.dat,https://github.com/Loyalsoldier/v2ray-rules-dat/releases/download/202307132210/geosite.dat,90543f9afcb86dff4076db6dfd4e32a8fa59fe1910f525cbd786b461cabcb611))

all: clean .deps/geoip/prepared .deps/geosite/prepared
	$(MAKE) clean

.deps/geoip/prepared: .distfiles/geoip.dat
	cp .distfiles/geoip.dat .

.deps/geosite/prepared: .distfiles/geosite.dat
	cp .distfiles/geosite.dat .

clean:
	rm -rf .deps .distfiles

