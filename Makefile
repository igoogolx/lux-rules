define download =
.distfiles/$(1):
	mkdir -p .distfiles
	if ! curl -L#o $$@.unverified $(2); then rm -f $$@.unverified; exit 1; fi
	if ! echo "$(3)  $$@.unverified" | sha256sum -c; then rm -f $$@.unverified; exit 1; fi
	if ! mv $$@.unverified $$@; then rm -f $$@.unverified; exit 1; fi
endef

$(eval $(call download,geoip.dat,https://github.com/Loyalsoldier/v2ray-rules-dat/releases/download/202208132211/geoip.dat,cc64fd248239e216545e3a93b287548560be2dcfff5403edae34c8c3655be204))
$(eval $(call download,geosite.dat,https://github.com/Loyalsoldier/v2ray-rules-dat/releases/download/202208132211/geosite.dat,7a5c2af16b7ad7298e0d24d15206adcda6c7e5c44ecf8a0a4fa072a14037ae16))

all: clean .deps/geoip/prepared .deps/geosite/prepared
	$(MAKE) clean

.deps/geoip/prepared: .distfiles/geoip.dat
	cp .distfiles/geoip.dat .

.deps/geosite/prepared: .distfiles/geosite.dat
	cp .distfiles/geosite.dat .

clean:
	rm -rf .deps .distfiles

