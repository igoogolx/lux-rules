package geo_data

import (
	"fmt"
	router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"google.golang.org/protobuf/proto"
	"os"
)

func LoadGeoIpFile(filename string) ([]*router.GeoIP, error) {
	geoipBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v ", filename)
	}
	var geoipList router.GeoIPList
	if err := proto.Unmarshal(geoipBytes, &geoipList); err != nil {
		return nil, err
	}
	return geoipList.Entry, nil
}

func GetGeoIp(fileName string, countries []string) ([]*router.GeoIP, error) {
	geoList, err := LoadGeoIpFile(fileName)
	if err != nil {
		return nil, err
	}
	var ips = make([]*router.GeoIP, 0)
	for _, ip := range geoList {
		for _, country := range countries {
			if ip.CountryCode == country {
				ips = append(ips, ip)
			}
		}
	}
	return ips, nil
}

func GetGeoSites(filename string, countries []string) ([]*router.GeoSite, error) {
	geositeBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", filename)
	}
	var geositeList router.GeoSiteList
	if err := proto.Unmarshal(geositeBytes, &geositeList); err != nil {
		return nil, err
	}

	var sites []*router.GeoSite
	for _, site := range geositeList.Entry {
		for _, country := range countries {
			if site.CountryCode == country {

				sites = append(sites, site)
			}
		}
	}
	return sites, nil
}
