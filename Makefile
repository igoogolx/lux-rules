define download =
.distfiles/$(1):
	mkdir -p .distfiles
	if ! curl -L#o $$@.unverified $(2); then rm -f $$@.unverified; exit 1; fi
	if ! echo "$(3)  $$@.unverified" | sha256sum -c; then rm -f $$@.unverified; exit 1; fi
	if ! mv $$@.unverified $$@; then rm -f $$@.unverified; exit 1; fi
endef

$(eval $(call download,geoip.dat,https://github.com/Loyalsoldier/v2ray-rules-dat/releases/download/202204172211/geoip.dat,5f4aa741f330c5c99eca57ce2b47b4be1bfdcc2ec561d4213365fe0377101ceb))
$(eval $(call download,geosite.dat,https://github.com/Loyalsoldier/v2ray-rules-dat/releases/download/202204172211/geosite.dat,c6beb8f04c30099799917ca5dd202e08a164a089f9e59868b540248f8a96859b))

all: clean .deps/geoip/prepared .deps/geosite/prepared
	$(MAKE) clean

.deps/geoip/prepared: .distfiles/geoip.dat
	cp .distfiles/geoip.dat .

.deps/geosite/prepared: .distfiles/geosite.dat
	cp .distfiles/geosite.dat .

clean:
	rm -rf .deps .distfiles

