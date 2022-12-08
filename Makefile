define download =
.distfiles/$(1):
	mkdir -p .distfiles
	if ! curl -L#o $$@.unverified $(2); then rm -f $$@.unverified; exit 1; fi
	if ! echo "$(3)  $$@.unverified" | sha256sum -c; then rm -f $$@.unverified; exit 1; fi
	if ! mv $$@.unverified $$@; then rm -f $$@.unverified; exit 1; fi
endef

$(eval $(call download,geoip.dat,https://github.com/Loyalsoldier/v2ray-rules-dat/releases/download/202212072210/geoip.dat,8c58d22cb94bf98a42b1b2dff8ac9c39f42f1e83f52dc1ab016c72e8a22c5fcb))
$(eval $(call download,geosite.dat,https://github.com/Loyalsoldier/v2ray-rules-dat/releases/download/202212072210/geosite.dat,1c7072652e017abe6c306d83a2a077d2e28f85dab912ee6df96c42c81b98d5ea))

all: clean .deps/geoip/prepared .deps/geosite/prepared
	$(MAKE) clean

.deps/geoip/prepared: .distfiles/geoip.dat
	cp .distfiles/geoip.dat .

.deps/geosite/prepared: .distfiles/geosite.dat
	cp .distfiles/geosite.dat .

clean:
	rm -rf .deps .distfiles

