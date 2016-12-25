package hkodata_test

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	pretty "github.com/tonnerre/golang-pretty"
	"github.com/yookoala/weatherhk/hkodata"
)

func TestDecodeRegionJSON(t *testing.T) {
	file, err := os.Open("./test/region_json.201612191037.xml")
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	defer file.Close()

	regions, _ := hkodata.DecodeRegionJSON(file)
	if want, have := time.Date(2016, time.December, 19, 10, 20, 0, 0, hkodata.HKT), regions.PubDate; !want.Equal(have) {
		t.Errorf("expected %s, got %s", want, have)
	}

	expected := []hkodata.Region{
		hkodata.Region{
			Name:             hkodata.RegionName("hka"),
			ShortName:        "hka",
			CurrentTemp:      hkodata.NewTemperature(24.1),
			RelativeHumidity: hkodata.NewRelativeHumidity(.54),
			WindDirection:    "E",
			WindSpeed:        hkodata.NewSpeed(24),
			MaxTemp:          hkodata.NewTemperature(24.6),
			MinTemp:          hkodata.NewTemperature(19.5),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("cch"),
			ShortName:        "cch",
			CurrentTemp:      hkodata.NewTemperature(22.0),
			RelativeHumidity: hkodata.NewRelativeHumidity(.68),
			WindDirection:    "SE",
			WindSpeed:        hkodata.NewSpeed(30),
			MaxTemp:          hkodata.NewTemperature(22.0),
			MinTemp:          hkodata.NewTemperature(18.4),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("hpv"),
			ShortName:        "hpv",
			CurrentTemp:      hkodata.NewTemperature(23.1),
			RelativeHumidity: nil,
			WindDirection:    "",
			WindSpeed:        nil,
			MaxTemp:          hkodata.NewTemperature(23.1),
			MinTemp:          hkodata.NewTemperature(17.1),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("hko"),
			ShortName:        "hko",
			CurrentTemp:      hkodata.NewTemperature(20.8),
			RelativeHumidity: hkodata.NewRelativeHumidity(.74),
			WindDirection:    "",
			WindSpeed:        nil,
			MaxTemp:          hkodata.NewTemperature(20.8),
			MinTemp:          hkodata.NewTemperature(18.5),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("hkp"),
			ShortName:        "hkp",
			CurrentTemp:      hkodata.NewTemperature(21.6),
			RelativeHumidity: nil,
			WindDirection:    "",
			WindSpeed:        nil,
			MaxTemp:          hkodata.NewTemperature(21.6),
			MinTemp:          hkodata.NewTemperature(17.5),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("se1"),
			ShortName:        "se1",
			CurrentTemp:      hkodata.NewTemperature(21.8),
			RelativeHumidity: nil,
			WindDirection:    "",
			WindSpeed:        nil,
			MaxTemp:          hkodata.NewTemperature(22.5),
			MinTemp:          hkodata.NewTemperature(19.0),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("ksc"),
			ShortName:        "ksc",
			CurrentTemp:      hkodata.NewTemperature(22.2),
			RelativeHumidity: hkodata.NewRelativeHumidity(.66),
			WindDirection:    "",
			WindSpeed:        nil,
			MaxTemp:          hkodata.NewTemperature(22.2),
			MinTemp:          hkodata.NewTemperature(15.5),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("kp"),
			ShortName:        "kp",
			CurrentTemp:      hkodata.NewTemperature(21.7),
			RelativeHumidity: hkodata.NewRelativeHumidity(.66),
			WindDirection:    "SE",
			WindSpeed:        hkodata.NewSpeed(9),
			MaxTemp:          hkodata.NewTemperature(21.8),
			MinTemp:          hkodata.NewTemperature(17.7),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("klt"),
			ShortName:        "klt",
			CurrentTemp:      hkodata.NewTemperature(22.5),
			RelativeHumidity: nil,
			WindDirection:    "",
			WindSpeed:        nil,
			MaxTemp:          hkodata.NewTemperature(23.2),
			MinTemp:          hkodata.NewTemperature(17.8),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("ktg"),
			ShortName:        "ktg",
			CurrentTemp:      hkodata.NewTemperature(21.9),
			RelativeHumidity: nil,
			WindDirection:    "",
			WindSpeed:        nil,
			MaxTemp:          hkodata.NewTemperature(22.0),
			MinTemp:          hkodata.NewTemperature(18.4),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("lfs"),
			ShortName:        "lfs",
			CurrentTemp:      hkodata.NewTemperature(23.8),
			RelativeHumidity: hkodata.NewRelativeHumidity(.62),
			WindDirection:    "SE",
			WindSpeed:        hkodata.NewSpeed(8),
			MaxTemp:          hkodata.NewTemperature(23.8),
			MinTemp:          hkodata.NewTemperature(15.9),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("ngp"),
			ShortName:        "ngp",
			CurrentTemp:      hkodata.NewTemperature(19.3),
			RelativeHumidity: nil,
			WindDirection:    "E",
			WindSpeed:        hkodata.NewSpeed(39),
			MaxTemp:          hkodata.NewTemperature(20.1),
			MinTemp:          hkodata.NewTemperature(17.5),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("tyw"),
			ShortName:        "tyw",
			CurrentTemp:      hkodata.NewTemperature(22.4),
			RelativeHumidity: nil,
			WindDirection:    "",
			WindSpeed:        nil,
			MaxTemp:          hkodata.NewTemperature(22.4),
			MinTemp:          hkodata.NewTemperature(12.4),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("pen"),
			ShortName:        "pen",
			CurrentTemp:      hkodata.NewTemperature(21.8),
			RelativeHumidity: hkodata.NewRelativeHumidity(.72),
			WindDirection:    "NE",
			WindSpeed:        hkodata.NewSpeed(18),
			MaxTemp:          hkodata.NewTemperature(21.8),
			MinTemp:          hkodata.NewTemperature(18.7),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("skg"),
			ShortName:        "skg",
			CurrentTemp:      hkodata.NewTemperature(20.6),
			RelativeHumidity: hkodata.NewRelativeHumidity(.72),
			WindDirection:    "SE",
			WindSpeed:        hkodata.NewSpeed(8),
			MaxTemp:          hkodata.NewTemperature(20.6),
			MinTemp:          hkodata.NewTemperature(16.6),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("ssp"),
			ShortName:        "ssp",
			CurrentTemp:      hkodata.NewTemperature(24.1),
			RelativeHumidity: nil,
			WindDirection:    "",
			WindSpeed:        nil,
			MaxTemp:          hkodata.NewTemperature(24.3),
			MinTemp:          hkodata.NewTemperature(18.0),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("sha"),
			ShortName:        "sha",
			CurrentTemp:      hkodata.NewTemperature(22.9),
			RelativeHumidity: hkodata.NewRelativeHumidity(.59),
			WindDirection:    "SE",
			WindSpeed:        hkodata.NewSpeed(5),
			MaxTemp:          hkodata.NewTemperature(23.0),
			MinTemp:          hkodata.NewTemperature(15.8),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("skw"),
			ShortName:        "skw",
			CurrentTemp:      hkodata.NewTemperature(21.2),
			RelativeHumidity: nil,
			WindDirection:    "",
			WindSpeed:        nil,
			MaxTemp:          hkodata.NewTemperature(21.7),
			MinTemp:          hkodata.NewTemperature(17.9),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("sek"),
			ShortName:        "sek",
			CurrentTemp:      hkodata.NewTemperature(22.9),
			RelativeHumidity: hkodata.NewRelativeHumidity(.63),
			WindDirection:    "",
			WindSpeed:        nil,
			MaxTemp:          hkodata.NewTemperature(22.9),
			MinTemp:          hkodata.NewTemperature(15.4),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("ssh"),
			ShortName:        "ssh",
			CurrentTemp:      hkodata.NewTemperature(21.1),
			RelativeHumidity: hkodata.NewRelativeHumidity(.70),
			WindDirection:    "",
			WindSpeed:        nil,
			MaxTemp:          hkodata.NewTemperature(21.1),
			MinTemp:          hkodata.NewTemperature(16.1),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("sty"),
			ShortName:        "sty",
			CurrentTemp:      hkodata.NewTemperature(20.6),
			RelativeHumidity: nil,
			WindDirection:    "E",
			WindSpeed:        hkodata.NewSpeed(20),
			MaxTemp:          hkodata.NewTemperature(20.7),
			MinTemp:          hkodata.NewTemperature(18.6),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("tkl"),
			ShortName:        "tkl",
			CurrentTemp:      hkodata.NewTemperature(21.9),
			RelativeHumidity: hkodata.NewRelativeHumidity(.65),
			WindDirection:    "E",
			WindSpeed:        hkodata.NewSpeed(10),
			MaxTemp:          hkodata.NewTemperature(21.9),
			MinTemp:          hkodata.NewTemperature(15.0),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("tms"),
			ShortName:        "tms",
			CurrentTemp:      hkodata.NewTemperature(16.4),
			RelativeHumidity: nil,
			WindDirection:    "",
			WindSpeed:        nil,
			MaxTemp:          hkodata.NewTemperature(17.0),
			MinTemp:          hkodata.NewTemperature(12.6),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("tpo"),
			ShortName:        "tpo",
			CurrentTemp:      hkodata.NewTemperature(22.9),
			RelativeHumidity: hkodata.NewRelativeHumidity(.68),
			WindDirection:    "",
			WindSpeed:        nil,
			MaxTemp:          hkodata.NewTemperature(22.9),
			MinTemp:          hkodata.NewTemperature(16.7),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("vp1"),
			ShortName:        "vp1",
			CurrentTemp:      hkodata.NewTemperature(20.4),
			RelativeHumidity: nil,
			WindDirection:    "",
			WindSpeed:        nil,
			MaxTemp:          hkodata.NewTemperature(20.4),
			MinTemp:          hkodata.NewTemperature(15.5),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("jkb"),
			ShortName:        "jkb",
			CurrentTemp:      hkodata.NewTemperature(22.3),
			RelativeHumidity: hkodata.NewRelativeHumidity(.67),
			WindDirection:    "E",
			WindSpeed:        hkodata.NewSpeed(8),
			MaxTemp:          hkodata.NewTemperature(22.6),
			MinTemp:          hkodata.NewTemperature(17.0),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("ty1"),
			ShortName:        "ty1",
			CurrentTemp:      hkodata.NewTemperature(22.1),
			RelativeHumidity: hkodata.NewRelativeHumidity(.61),
			WindDirection:    "E",
			WindSpeed:        hkodata.NewSpeed(8),
			MaxTemp:          hkodata.NewTemperature(22.1),
			MinTemp:          hkodata.NewTemperature(16.6),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("twn"),
			ShortName:        "twn",
			CurrentTemp:      hkodata.NewTemperature(21.9),
			RelativeHumidity: hkodata.NewRelativeHumidity(.64),
			WindDirection:    "",
			WindSpeed:        nil,
			MaxTemp:          hkodata.NewTemperature(21.9),
			MinTemp:          hkodata.NewTemperature(16.1),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("tw"),
			ShortName:        "tw",
			CurrentTemp:      hkodata.NewTemperature(22.8),
			RelativeHumidity: hkodata.NewRelativeHumidity(.63),
			WindDirection:    "",
			WindSpeed:        nil,
			MaxTemp:          hkodata.NewTemperature(22.8),
			MinTemp:          hkodata.NewTemperature(17.6),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("tun"),
			ShortName:        "tun",
			CurrentTemp:      hkodata.NewTemperature(23.8),
			RelativeHumidity: hkodata.NewRelativeHumidity(.58),
			WindDirection:    "S",
			WindSpeed:        hkodata.NewSpeed(8),
			MaxTemp:          hkodata.NewTemperature(24.3),
			MinTemp:          hkodata.NewTemperature(17.9),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("wgl"),
			ShortName:        "wgl",
			CurrentTemp:      hkodata.NewTemperature(22.4),
			RelativeHumidity: hkodata.NewRelativeHumidity(.72),
			WindDirection:    "E",
			WindSpeed:        hkodata.NewSpeed(23),
			MaxTemp:          hkodata.NewTemperature(22.7),
			MinTemp:          hkodata.NewTemperature(18.6),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("wlp"),
			ShortName:        "wlp",
			CurrentTemp:      hkodata.NewTemperature(23.6),
			RelativeHumidity: hkodata.NewRelativeHumidity(.59),
			WindDirection:    "NE",
			WindSpeed:        hkodata.NewSpeed(6),
			MaxTemp:          hkodata.NewTemperature(24.0),
			MinTemp:          hkodata.NewTemperature(15.3),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("hks"),
			ShortName:        "hks",
			CurrentTemp:      hkodata.NewTemperature(22.3),
			RelativeHumidity: hkodata.NewRelativeHumidity(.62),
			WindDirection:    "NE",
			WindSpeed:        hkodata.NewSpeed(14),
			MaxTemp:          hkodata.NewTemperature(22.3),
			MinTemp:          hkodata.NewTemperature(17.2),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("wts"),
			ShortName:        "wts",
			CurrentTemp:      hkodata.NewTemperature(23.6),
			RelativeHumidity: nil,
			WindDirection:    "",
			WindSpeed:        nil,
			MaxTemp:          hkodata.NewTemperature(23.6),
			MinTemp:          hkodata.NewTemperature(17.5),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("ylp"),
			ShortName:        "ylp",
			CurrentTemp:      hkodata.NewTemperature(23.8),
			RelativeHumidity: nil,
			WindDirection:    "",
			WindSpeed:        nil,
			MaxTemp:          hkodata.NewTemperature(23.8),
			MinTemp:          hkodata.NewTemperature(15.1),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("tc"),
			ShortName:        "tc",
			CurrentTemp:      hkodata.NewTemperature(17.5),
			RelativeHumidity: nil,
			WindDirection:    "E",
			WindSpeed:        hkodata.NewSpeed(28),
			MaxTemp:          hkodata.NewTemperature(17.6),
			MinTemp:          hkodata.NewTemperature(14.0),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("gi"),
			ShortName:        "gi",
			CurrentTemp:      nil,
			RelativeHumidity: nil,
			WindDirection:    "NE",
			WindSpeed:        hkodata.NewSpeed(29),
			MaxTemp:          nil,
			MinTemp:          nil,
		},
		hkodata.Region{
			Name:             hkodata.RegionName("se"),
			ShortName:        "se",
			CurrentTemp:      nil,
			RelativeHumidity: nil,
			WindDirection:    "E",
			WindSpeed:        hkodata.NewSpeed(17),
			MaxTemp:          nil,
			MinTemp:          nil,
		},
		hkodata.Region{
			Name:             hkodata.RegionName("sc"),
			ShortName:        "sc",
			CurrentTemp:      nil,
			RelativeHumidity: nil,
			WindDirection:    "E",
			WindSpeed:        hkodata.NewSpeed(18),
			MaxTemp:          nil,
			MinTemp:          nil,
		},
		hkodata.Region{
			Name:             hkodata.RegionName("sf"),
			ShortName:        "sf",
			CurrentTemp:      nil,
			RelativeHumidity: nil,
			WindDirection:    "E",
			WindSpeed:        hkodata.NewSpeed(16),
			MaxTemp:          nil,
			MinTemp:          nil,
		},
		hkodata.Region{
			Name:             hkodata.RegionName("plc"),
			ShortName:        "plc",
			CurrentTemp:      nil,
			RelativeHumidity: nil,
			WindDirection:    "E",
			WindSpeed:        hkodata.NewSpeed(11),
			MaxTemp:          nil,
			MinTemp:          nil,
		},
		hkodata.Region{
			Name:             hkodata.RegionName("tpk"),
			ShortName:        "tpk",
			CurrentTemp:      nil,
			RelativeHumidity: nil,
			WindDirection:    "E",
			WindSpeed:        hkodata.NewSpeed(11),
			MaxTemp:          nil,
			MinTemp:          nil,
		},
		hkodata.Region{
			Name:             hkodata.RegionName("tap"),
			ShortName:        "tap",
			CurrentTemp:      nil,
			RelativeHumidity: nil,
			WindDirection:    "SE",
			WindSpeed:        hkodata.NewSpeed(20),
			MaxTemp:          nil,
			MinTemp:          nil,
		},
	}

	for i := 0; i < len(expected); i++ {

		// log the json results
		bytes, _ := json.Marshal(regions.Regions[i])
		t.Logf("json.regions[%d]: %s", i, bytes)

		// compare values
		if want, have := expected[i], regions.Regions[i]; !reflect.DeepEqual(want, have) {
			t.Errorf("unexpected difference in %#v / %#v (want != have)", want.ShortName, want.Name.Zh)
			for _, desc := range pretty.Diff(want, have) {
				t.Log("\tregion." + desc)
			}
		}

		// check if name actually found
		if regions.Regions[i].Name.En == "" {
			t.Errorf("The English name of %#v is empty string", regions.Regions[i].ShortName)
		}
		if regions.Regions[i].Name.Zh == "" {
			t.Errorf("The Chinese name of %#v is empty string", regions.Regions[i].ShortName)
		}
	}

}

func TestDecodeRegionJSON_dummyCase(t *testing.T) {

	file, err := os.Open("./test/region_json.dummy.xml")
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	defer file.Close()

	regions, err := hkodata.DecodeRegionJSON(file)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if want, have := time.Date(2016, time.December, 19, 10, 20, 0, 0, hkodata.HKT), regions.PubDate; !want.Equal(have) {
		t.Errorf("expected %s, got %s", want, have)
	}

	expected := []hkodata.Region{
		hkodata.Region{
			Name:             hkodata.RegionName("hka"),
			ShortName:        "hka",
			CurrentTemp:      hkodata.NewTemperature(24.1),
			RelativeHumidity: hkodata.NewRelativeHumidity(.54),
			WindDirection:    "E",
			WindSpeed:        hkodata.NewSpeed(24),
			MaxTemp:          hkodata.NewTemperature(24.6),
			MinTemp:          hkodata.NewTemperature(19.5),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("cch"),
			ShortName:        "cch",
			CurrentTemp:      hkodata.NewTemperature(0),
			RelativeHumidity: hkodata.NewRelativeHumidity(0),
			WindDirection:    "SE",
			WindSpeed:        hkodata.NewSpeed(0),
			MaxTemp:          hkodata.NewTemperature(0),
			MinTemp:          hkodata.NewTemperature(0),
		},
		hkodata.Region{
			Name:             hkodata.RegionName("hpv"),
			ShortName:        "hpv",
			CurrentTemp:      nil,
			RelativeHumidity: nil,
			WindDirection:    "",
			WindSpeed:        nil,
			MaxTemp:          nil,
			MinTemp:          nil,
		},
		hkodata.Region{
			Name:             hkodata.RegionName("hks"),
			ShortName:        "hks",
			CurrentTemp:      hkodata.NewTemperature(-10),
			RelativeHumidity: hkodata.NewRelativeHumidity(-.1),
			WindDirection:    "",
			WindSpeed:        hkodata.NewSpeed(-10),
			MaxTemp:          hkodata.NewTemperature(-10),
			MinTemp:          hkodata.NewTemperature(-10),
		},
	}

	for i := 0; i < len(expected); i++ {
		// log the json results
		bytes, _ := json.Marshal(regions.Regions[i])
		t.Logf("json.regions[%d]: %s", i, bytes)

		// compare values
		if want, have := expected[i], regions.Regions[i]; !reflect.DeepEqual(want, have) {
			t.Errorf("unexpected difference in %#v / %#v (want != have)", want.ShortName, want.Name.Zh)
			for _, desc := range pretty.Diff(want, have) {
				t.Log("\tregion." + desc)
			}
		}

		// check if name actually found
		if regions.Regions[i].Name.En == "" {
			t.Errorf("The English name of %#v is empty string", regions.Regions[i].ShortName)
		}
		if regions.Regions[i].Name.Zh == "" {
			t.Errorf("The Chinese name of %#v is empty string", regions.Regions[i].ShortName)
		}
	}

}
