package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/v2fly/v2ray-core/v4/app/router"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

var (
	ipDir   = filepath.Join(".", "geoData", "ip")
	siteDir = filepath.Join(".", "geoData", "site")

	ipFileName   = "geoip.dat"
	siteFileName = "geosite.dat"

	outIps   = []string{"CN", "PRIVATE"}
	outSites = []string{"GFW", "CN"}
)

func writeIpFile(filePath string, ips []*router.CIDR) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	for _, cidr := range ips {
		ipA := net.IPNet{
			IP:   cidr.Ip,
			Mask: net.CIDRMask(int(cidr.Prefix), 8*len(cidr.Ip)),
		}
		line := ipA.String() + "\n"
		_, err := f.Write([]byte(line))
		if err != nil {
			return fmt.Errorf("fail to write ip:%v to %v", line, filePath)
		}
	}
	return nil
}

func writeDomainFile(filePath string, domains []*router.Domain) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	for _, domain := range domains {
		line := domain.Value + "/" + strconv.Itoa(int(domain.Type)) + "\n"
		_, err := f.Write([]byte(line))
		if err != nil {
			return fmt.Errorf("fail to write domain:%v to %v", line, filePath)
		}
	}
	return nil
}

func genIpFile(fileName string, countries []string) error {
	geoList, err := LoadGroIpFile(fileName)
	if err != nil {
		return err
	}
	for _, geoData := range geoList {
		for _, country := range countries {
			if geoData.CountryCode == country {
				err := writeIpFile(filepath.Join(ipDir, country), geoData.Cidr)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func LoadGroIpFile(filename string) ([]*router.GeoIP, error) {
	geoipBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v ", filename)
	}
	var geoipList router.GeoIPList
	if err := proto.Unmarshal(geoipBytes, &geoipList); err != nil {
		return nil, err
	}
	return geoipList.Entry, nil
}

func genSiteFile(filename string, countries []string) error {
	geositeBytes, err := ioutil.ReadFile(filename)
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
				err := writeDomainFile(filepath.Join(siteDir, country), site.Domain)
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

func main() {
	createDirIfNotExist(ipDir)
	createDirIfNotExist(siteDir)
	err := genIpFile(ipFileName, outIps)
	if err != nil {
		log.Fatalf("fail to gen geo ip file,error:%v", err)
	}
	err = genSiteFile(siteFileName, outSites)
	if err != nil {
		log.Fatalf("fail to gen geo site file,error:%v", err)
	}
}
