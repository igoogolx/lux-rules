package main

import (
	"fmt"
	geodata "github.com/igoogolx/lux-geo-data/geo-data"
	router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"net/netip"
	"os"
	"path/filepath"
)

var (
	ruleDir      = filepath.Join(".", "rules")
	ipFileName   = "geoip.dat"
	siteFileName = "geosite.dat"
)

type Policy string

const (
	PolicyDirect Policy = "DIRECT"
	PolicyProxy  Policy = "PROXY"
	PolicyReject Policy = "REJECT"
)

func getDomainType(rType router.Domain_Type) (string, error) {
	switch rType {
	case router.Domain_Plain:
		return "DOMAIN-KEYWORD", nil
	case router.Domain_Regex:
		return "DOMAIN-REGEX", nil
	case router.Domain_RootDomain:
		return "DOMAIN-SUFFIX", nil
	case router.Domain_Full:
		return "DOMAIN", nil
	}
	return "", fmt.Errorf("invalid domain type")
}

func writeIpFile(filePath string, ips []*router.CIDR, policy Policy) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("fail to close ip file: %v", filePath)
		}
	}(f)
	if err != nil {
		log.Fatal(err)
	}
	for _, cidr := range ips {
		ipA := net.IPNet{
			IP:   cidr.Ip,
			Mask: net.CIDRMask(int(cidr.Prefix), 8*len(cidr.Ip)),
		}
		line := "IP-CIDR," + ipA.String() + "," + string(policy) + "\n"
		_, err := f.Write([]byte(line))
		if err != nil {
			return fmt.Errorf("fail to write ip:%v to %v", line, filePath)
		}
	}
	return nil
}

func writeDomainFile(filePath string, domains []*router.Domain, policy Policy) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("fail to close domain file: %v", filePath)
		}
	}(f)
	if err != nil {
		log.Fatal(err)
	}
	for _, domain := range domains {
		domainType, err := getDomainType(domain.Type)
		if err != nil {
			return err
		}
		line := domainType + "," + domain.Value + "," + string(policy) + "\n"
		_, err = f.Write([]byte(line))
		if err != nil {
			return fmt.Errorf("fail to write domain:%v to %v", line, filePath)
		}
	}
	return nil
}

func genIpFile(fileName string, countries []string, policy Policy, name string) error {
	geoList, err := geodata.LoadGeoIpFile(fileName)
	if err != nil {
		return err
	}
	for _, geoData := range geoList {
		for _, country := range countries {
			if geoData.CountryCode == country {
				err := writeIpFile(filepath.Join(ruleDir, name), geoData.Cidr, policy)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func genSiteFile(filename string, countries []string, policy Policy, name string) error {
	geositeBytes, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", filename)
	}
	var geositeList router.GeoSiteList
	if err := proto.Unmarshal(geositeBytes, &geositeList); err != nil {
		return err
	}

	for _, site := range geositeList.Entry {
		for _, country := range countries {
			if site.CountryCode == country {
				err := writeDomainFile(filepath.Join(ruleDir, name), site.Domain, policy)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func createDirIfNotExist(dir string) {
	newPath := filepath.Join(dir)
	_ = os.MkdirAll(newPath, os.ModePerm)
}

func createBypassCn() {
	name := "bypass_cn"
	err := genIpFile(ipFileName, []string{"PRIVATE", "CN"}, PolicyDirect, name)
	if err != nil {
		log.Fatalf("fail to gen geo ip file,error:%v", err)
	}
	err = genSiteFile(siteFileName, []string{"CN"}, PolicyDirect, name)
	if err != nil {
		log.Fatalf("fail to gen geo site file,error:%v", err)
	}
}
func createProxyAll() {
	name := "proxy_all"
	err := genIpFile(ipFileName, []string{"PRIVATE"}, PolicyDirect, name)
	if err != nil {
		log.Fatalf("fail to gen geo ip file,error:%v", err)
	}
}

func createBypassAll() {
	name := "bypass_all"
	err := writeDomainFile(filepath.Join(ruleDir, name), []*router.Domain{
		{Type: router.Domain_Regex, Value: ".*"},
	}, PolicyDirect)
	if err != nil {
		log.Fatalf("fail to write domain file,error:%v", err)
	}
	allAddr, err := netip.ParseAddr("0.0.0.0")
	if err != nil {
		log.Fatalf("fail to parse ip,error:%v", err)
	}
	err = writeIpFile(filepath.Join(ruleDir, name), []*router.CIDR{
		{IpAddr: "0.0.0.0", Prefix: 32, Ip: allAddr.AsSlice()},
	}, PolicyDirect)
	if err != nil {
		log.Fatalf("fail to write domain file,error:%v", err)
	}
}

func createProxyGfw() {
	name := "proxy_gfw"
	err := genSiteFile(siteFileName, []string{"GFW"}, PolicyProxy, name)
	if err != nil {
		log.Fatalf("fail to gen geo site file,error:%v", err)
	}
	err = writeDomainFile(filepath.Join(ruleDir, name), []*router.Domain{
		{Type: router.Domain_Regex, Value: ".*"},
	}, PolicyDirect)
	if err != nil {
		log.Fatalf("fail to write domain file,error:%v", err)
	}
	allAddr, err := netip.ParseAddr("0.0.0.0")
	if err != nil {
		log.Fatalf("fail to parse ip,error:%v", err)
	}
	err = writeIpFile(filepath.Join(ruleDir, name), []*router.CIDR{
		{IpAddr: "0.0.0.0", Prefix: 32, Ip: allAddr.AsSlice()},
	}, PolicyDirect)
	if err != nil {
		log.Fatalf("fail to write domain file,error:%v", err)
	}
}

func main() {
	_ = os.RemoveAll(ruleDir)
	createDirIfNotExist(ruleDir)
	createProxyAll()
	createBypassCn()
	createBypassAll()
	createProxyGfw()
}
