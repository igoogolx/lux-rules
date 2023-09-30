define download =
.distfiles/$(1):
	mkdir -p .distfiles
	if ! curl -L#o $$@.unverified $(2); then rm -f $$@.unverified; exit 1; fi
	if ! echo "$(3)  $$@.unverified" | sha256sum -c; then rm -f $$@.unverified; exit 1; fi
	if ! mv $$@.unverified $$@; then rm -f $$@.unverified; exit 1; fi
endef

$(eval $(call download,geoip.dat,https://github.com/Loyalsoldier/v2ray-rules-dat/releases/download/202309292208/geoip.dat,a86b5ecbc00f779ecf7236d72545152f1302f7d94dc5d786e55356d83f025edf))
$(eval $(call download,geosite.dat,https://github.com/Loyalsoldier/v2ray-rules-dat/releases/download/202309292208/geosite.dat,8e57dfe4d24c540f9fe2a9360ab87fd1e47885fae493aecc93b0159022560daa))

all: clean .deps/geoip/prepared .deps/geosite/prepared
	$(MAKE) clean

.deps/geoip/prepared: .distfiles/geoip.dat
	cp .distfiles/geoip.dat .

.deps/geosite/prepared: .distfiles/geosite.dat
	cp .distfiles/geosite.dat .

clean:
	rm -rf .deps .distfiles

