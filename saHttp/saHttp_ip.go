package saHttp

import (
	"fmt"
	"github.com/saxon134/go-utils/saData"
	"net"
	"regexp"
)

func IsIpAddress(str string) bool {
	if str == "" || len(str) < 7 || len(str) > 15 {
		return false
	}

	reg := `^\d{1,3}[\.]\d{1,3}[\.]\d{1,3}[\.]\d{1,3}$`
	ok, err := regexp.MatchString(reg, str)
	if err != nil {
		return false
	}
	return ok
}

func GetIpRegion(ip string) *map[string]string {
	if IsIpAddress(ip) {
		reg := map[string]string{}
		if _res, err := Get("https://api.map.baidu.com/location/ip?ak=Tp7rCYFxLmiTf0EZRpc55AgdvExlLePI&coor=bd09ll&ip="+ip, nil); err == nil {
			if res, err := saData.ToMap(_res); err == nil && res != nil {
				if d, _ := saData.ToMap(res["content"]); d != nil {
					if address, _ := saData.ToMap(d["address_detail"]); address != nil {
						reg["province"], _ = saData.ToStr(d["province"])
						reg["city"], _ = saData.ToStr(d["city"])
						reg["district"], _ = saData.ToStr(d["street"])
					}

					if point, _ := saData.ToMap(d["point"]); point != nil {
						reg["long"], _ = saData.ToStr(point["x"])
						reg["lat"], _ = saData.ToStr(point["y"])
					}

					return &reg
				}
			}
		}
	}
	return nil
}

func GetLocalIP() (ips []string, err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("get ip interfaces error:", err)
		return nil, err
	}

	for _, i := range ifaces {
		addrs, errRet := i.Addrs()
		if errRet != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				if ip.IsGlobalUnicast() {
					ips = append(ips, ip.String())
				}
			}
		}
	}
	return
}
