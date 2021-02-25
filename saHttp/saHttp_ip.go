package saHttp

import (
	"github.com/saxon134/go-utils/saData"
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
		if _res, err := Get("http://ip.taobao.com/service/getIpInfo.php?ip="+ip, nil); err == nil {
			if res, err := saData.JsonToMap(_res); err == nil && res != nil {
				if d, _ := saData.ToMap(res["data"]); d != nil {
					if v, _ := saData.ToStr(d["region_id"]); v != "" {
						reg["provinceCode"] = v
						reg["provinceName"], _ = saData.ToStr(d["region"])
					}
					if v, _ := saData.ToStr(d["city_id"]); v != "" {
						reg["cityCode"] = v
						reg["cityName"], _ = saData.ToStr(d["city"])
					}
					if v, _ := saData.ToStr(d["area_id"]); v != "" {
						reg["districtCode"] = v
						reg["districtName"], _ = saData.ToStr(d["area"])
					}
					return &reg
				}
			}
		}
	}
	return nil
}
