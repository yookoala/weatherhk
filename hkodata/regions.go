package hkodata

import (
	"io"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/go-restit/lzjson"
)

// I18nName contains names of different languages
type I18nName struct {
	Zh string `json:"zh_HK,omitempty"`
	En string `json:"en,omitempty"`
}

var regionNames map[string]I18nName

// RegionName finds a region's name by its HKO shortname
func RegionName(shortName string) (name I18nName) {
	name, _ = regionNames[shortName]
	return
}

func init() {
	regionNames = map[string]I18nName{
		"cch": {
			Zh: "長洲",
			En: "Cheung Chau",
		},
		"cen": {
			Zh: "中環",
			En: "",
		},
		"gi": {
			Zh: "青洲",
			En: "Green Island",
		},
		"hka": {
			Zh: "赤鱲角",
			En: "Chek Lap Kok",
		},
		"hko": {
			Zh: "香港天文台",
			En: "Hong Kong Observatory",
		},
		"hkp": {
			Zh: "香港公園",
			En: "Hong Kong Park",
		},
		"hks": {
			Zh: "黃竹坑",
			En: "Wong Chuk Hang",
		},
		"hpv": {
			Zh: "跑馬地",
			En: "Happy Valley",
		},
		"jkb": {
			Zh: "將軍澳",
			En: "Tseung Kwan O",
		},
		"klt": {
			Zh: "九龍城",
			En: "Kowloon City",
		},
		"kp": {
			Zh: "京士柏",
			En: "King's Park",
		},
		"ksc": {
			Zh: "滘西洲",
			En: "Kau Sai Chau",
		},
		"ktg": {
			Zh: "觀塘",
			En: "Kwun Tong",
		},
		"lfs": {
			Zh: "流浮山",
			En: "Lau Fau Shan",
		},
		"ngp": {
			Zh: "昂 坪",
			En: "Ngong Ping",
		},
		"pen": {
			Zh: "坪洲",
			En: "Peng Chau",
		},
		"plc": {
			Zh: "大美督",
			En: "Tai Mei Tuk",
		},
		"sc": {
			Zh: "沙洲",
			En: "Sha Chau",
		},
		"se": {
			Zh: "啟德",
			En: "Kai Tak",
		},
		"sek": {
			Zh: "石崗",
			En: "Shek Kong",
		},
		"sf": {
			Zh: "天星碼頭",
			En: "Star Ferry",
		},
		"sha": {
			Zh: "沙田",
			En: "Sha Tin",
		},
		"skg": {
			Zh: "西貢",
			En: "Sai Kung",
		},
		"skw": {
			Zh: "筲箕灣",
			En: "Shau Kei Wan",
		},
		"ssh": {
			Zh: "上水",
			En: "Sheung Shui",
		},
		"ssp": {
			Zh: "深水埗",
			En: "Sham Shui Po",
		},
		"sty": {
			Zh: "赤柱",
			En: "Stanley",
		},
		"swh": {
			Zh: "西灣河",
			En: "",
		},
		"tap": {
			Zh: "塔門",
			En: "Tap Mun",
		},
		"tc": {
			Zh: "大老山",
			En: "Tate's Cairn",
		},
		"tkl": {
			Zh: "打鼓嶺",
			En: "Ta Kwu Ling",
		},
		"tms": {
			Zh: "大帽山",
			En: "Tai Mo Shan",
		},
		"tpk": {
			Zh: "大埔滘",
			En: "Tai Po Kau",
		},
		"tpo": {
			Zh: "大埔",
			En: "Tai Po",
		},
		"tun": {
			Zh: "屯門",
			En: "Tuen Mun",
		},
		"tw": {
			Zh: "荃灣城門谷",
			En: "Tsuen Wan Shing Mun Valley",
		},
		"twn": {
			Zh: "荃灣可觀",
			En: "Tsuen Wan Ho Koon",
		},
		"ty1": {
			Zh: "青衣",
			En: "Tsing Yi",
		},
		"tyw": {
			Zh: "北潭涌",
			En: "Pak Tam Chung",
		},
		"vp1": {
			Zh: "山頂",
			En: "The Peak",
		},
		"wgl": {
			Zh: "橫瀾島",
			En: "Waglan Island",
		},
		"wlp": {
			Zh: "濕地公園",
			En: "Wetland Park",
		},
		"wts": {
			Zh: "黃大仙",
			En: "Wong Tai Sin",
		},
		"ylp": {
			Zh: "元朗公園",
			En: "Yuen Long Park",
		},
	}
}

type reflectField struct {
	Type  string
	Index int
}

var regionDataFields map[string]reflectField

func init() {
	// reflect the Region struct for comparison
	typ := reflect.TypeOf(Region{})
	regionDataFields = make(map[string]reflectField)
	for i, numField := 0, typ.NumField(); i < numField; i++ {
		if typ.Field(i).Type.Kind() != reflect.Ptr {
			regionDataFields[typ.Field(i).Tag.Get("hkodata")] = reflectField{
				Index: i,
				Type:  typ.Field(i).Type.Name(),
			}
		} else {
			regionDataFields[typ.Field(i).Tag.Get("hkodata")] = reflectField{
				Index: i,
				Type:  "*" + typ.Field(i).Type.Elem().Name(),
			}
		}
	}
}

// Region represents data of a region from HKO non-public API endpoint
// `region_json.xml`
type Region struct {
	Name             I18nName          `hkodata:"-"`
	ShortName        string            `hkodata:"region"`
	CurrentTemp      *Temperature      `hkodata:"temp" json:"CurrentTemp,omitempty"`
	RelativeHumidity *RelativeHumidity `hkodata:"rh" json:"RelativeHumidity,omitempty"`
	WindDirection    string            `hkodata:"wind" json:"WindDirection,omitempty"`
	WindSpeed        *Speed            `hkodata:"speed" json:"WindSpeed,omitempty"`
	MaxTemp          *Temperature      `hkodata:"maxtemp" json:"MaxTemp,omitempty"`
	MinTemp          *Temperature      `hkodata:"mintemp" json:"MinTemp,omitempty"`
}

// Regions represents data from HKO non-public API endpoint `region_json.xml`
// (2015 API)
type Regions struct {
	PubDate time.Time
	Regions []Region
}

// Expires implements Expirer interface
func (regions Regions) Expires() time.Time {
	return regions.PubDate.Add(10 * time.Minute)
}

// DecodeRegionJSON decodes non-public API endpoint `region_json.xml` of
// HKO website (2015 API)
func DecodeRegionJSON(r io.Reader) (regions *Regions, err error) {
	regions = &Regions{}
	regions.Regions = make([]Region, 0, 20)

	json := lzjson.Decode(r)
	regions.PubDate, _ = time.Parse("200601021504-0700", json.Get("btime").String()+"+0800")

	// parse the field names from JSON
	var fieldNames []string
	json.Get("fields").Unmarshal(&fieldNames)

	// read data of all regions
	jsonData := json.Get("datas")
	for i, jsonDataLen := 0, jsonData.Len(); i < jsonDataLen; i++ {
		var region Region
		var fields []string
		jsonData.GetN(i).Unmarshal(&fields)
		regionVal := reflect.ValueOf(&region).Elem()

		// 1. loop each data fields of a region
		// 2. find the name of a field
		// 3. find equivlant field in the Region struct and set it
		for j, length := 0, len(fields); j < length; j++ {
			fieldName := fieldNames[j]

			if structFieldDef, ok := regionDataFields[fieldName]; ok {
				if fields[j] != "" {
					switch structFieldDef.Type {
					case "*Temperature":
						val, _ := strconv.ParseFloat(fields[j], 64)
						regionVal.Field(structFieldDef.Index).Set(reflect.ValueOf(NewTemperature(val)))
					case "*RelativeHumidity":
						val, _ := strconv.ParseFloat(fields[j], 64)
						regionVal.Field(structFieldDef.Index).Set(reflect.ValueOf(NewRelativeHumidity(val / 100)))
					case "*Speed":
						val, _ := strconv.ParseFloat(fields[j], 64)
						regionVal.Field(structFieldDef.Index).Set(reflect.ValueOf(NewSpeed(val)))
					case "string":
						regionVal.Field(structFieldDef.Index).Set(reflect.ValueOf(fields[j]))
					default:
						// note: should not be running here
						log.Printf("unhandled type: %s", structFieldDef.Type)
					}
				}
			}
		}

		// parse the short name into longer verbose region name
		region.Name, _ = regionNames[region.ShortName]

		// append the result list
		regions.Regions = append(regions.Regions, region)
	}
	return
}
