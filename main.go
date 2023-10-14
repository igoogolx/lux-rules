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
	"strconv"
)

var (
	ruleDir      = filepath.Join(".", "rules")
	ipFileName   = "geoip.dat"
	siteFileName = "geosite.dat"
)

func writeIpFile(filePath string, ips []*router.CIDR, policy string) error {
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
		line := "IP-CIDR," + ipA.String() + "," + policy + "\n"
		_, err := f.Write([]byte(line))
		if err != nil {
			return fmt.Errorf("fail to write ip:%v to %v", line, filePath)
		}
	}
	return nil
}

func writeDomainFile(filePath string, domains []*router.Domain, policy string) error {
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
		line := "DOMAIN," + domain.Value + "/" + strconv.Itoa(int(domain.Type)) + "," + policy + "\n"
		_, err := f.Write([]byte(line))
		if err != nil {
			return fmt.Errorf("fail to write domain:%v to %v", line, filePath)
		}
	}
	return nil
}

func genIpFile(fileName string, countries []string, policy string, name string) error {
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

func genSiteFile(filename string, countries []string, policy string, name string) error {
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
	err := genIpFile(ipFileName, []string{"PRIVATE", "CN"}, "bypass", name)
	if err != nil {
		log.Fatalf("fail to gen geo ip file,error:%v", err)
	}
	err = genSiteFile(siteFileName, []string{"CN"}, "bypass", name)
	if err != nil {
		log.Fatalf("fail to gen geo site file,error:%v", err)
	}
}
func createProxyAll() {
	name := "proxy_all"
	err := genIpFile(ipFileName, []string{"PRIVATE"}, "bypass", name)
	if err != nil {
		log.Fatalf("fail to gen geo ip file,error:%v", err)
	}
}

func createBypassAll() {
	name := "bypass_all"
	err := writeDomainFile(filepath.Join(ruleDir, name), []*router.Domain{
		{Type: router.Domain_Regex, Value: ".*"},
	}, "bypass")
	if err != nil {
		log.Fatalf("fail to write domain file,error:%v", err)
	}
	allAddr, err := netip.ParseAddr("0.0.0.0")
	if err != nil {
		log.Fatalf("fail to parse ip,error:%v", err)
	}
	err = writeIpFile(filepath.Join(ruleDir, name), []*router.CIDR{
		{IpAddr: "0.0.0.0", Prefix: 32, Ip: allAddr.AsSlice()},
	}, "bypass")
	if err != nil {
		log.Fatalf("fail to write domain file,error:%v", err)
	}
}

func createProxyGfw() {
	name := "proxy_gfw"
	err := genSiteFile(siteFileName, []string{"GFW"}, "proxy", name)
	if err != nil {
		log.Fatalf("fail to gen geo site file,error:%v", err)
	}
	err = writeDomainFile(filepath.Join(ruleDir, name), []*router.Domain{
		{Type: router.Domain_Regex, Value: ".*"},
	}, "bypass")
	if err != nil {
		log.Fatalf("fail to write domain file,error:%v", err)
	}
	allAddr, err := netip.ParseAddr("0.0.0.0")
	if err != nil {
		log.Fatalf("fail to parse ip,error:%v", err)
	}
	err = writeIpFile(filepath.Join(ruleDir, name), []*router.CIDR{
		{IpAddr: "0.0.0.0", Prefix: 32, Ip: allAddr.AsSlice()},
	}, "bypass")
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
