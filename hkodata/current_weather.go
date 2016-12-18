package hkodata

import (
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed/rss"
)

// Temperature contains Temperature in degree celcius
type Temperature int

// DistrictsTemperature contains Temperature of different districts in HK
type DistrictsTemperature struct {
	HongKongObservatory    Temperature
	KingsPark              Temperature
	WongChukHang           Temperature
	TaKwuLing              Temperature
	LauFauShan             Temperature
	TaiPo                  Temperature
	ShaTin                 Temperature
	TuenMun                Temperature
	TseungKwanO            Temperature
	SaiKung                Temperature
	CheungChau             Temperature
	ChekLapKok             Temperature
	ShekKong               Temperature
	TsuenWanHoKoon         Temperature
	TsuenWanShingMunValley Temperature
	HongKongPark           Temperature
	ShauKeiWan             Temperature
	KowloonCity            Temperature
	HappyValley            Temperature
	WongTaiSin             Temperature
	Stanley                Temperature
	KwunTong               Temperature
	ShamShuiPo             Temperature
	KaiTakRunwayPark       Temperature
	YuenLongPark           Temperature
}

// CurrentWeather contains all information of current weather in HKO's report
type CurrentWeather struct {
	PubDate              time.Time
	AirTemperature       Temperature
	RelativeHumidity     float64
	DistrictsTemperature DistrictsTemperature
	Raw                  string `json:"-"`
}

type ParseError []error

func (errs ParseError) Error() (msg string) {
	if len(errs) == 0 {
		return ""
	}

	msg = "ParseError:\n"
	for _, err := range errs {
		msg += err.Error() + "\n"
	}
	return
}

func DecodeCurrentWeather(r io.Reader) (data *CurrentWeather, err error) {

	// parse the content as RSS feed
	parser := rss.Parser{}
	feed, err := parser.Parse(r)
	if err != nil {
		return
	}

	// get description of the first item
	desc := strings.NewReader(feed.Items[0].Description)
	doc, err := goquery.NewDocumentFromReader(desc)
	if err != nil {
		return
	}

	if feed.Items[0].PubDateParsed == nil {
		err = fmt.Errorf("Failed to parse PubDate: %#v", feed.Items[0].PubDate)
		return
	}

	// prepare the parse the Temperature table
	reName := regexp.MustCompile(`[^\w]`)
	reDegree := regexp.MustCompile(`^.*?(\d+) degree.+?$`)
	data = &CurrentWeather{
		PubDate: *feed.Items[0].PubDateParsed,
		Raw:     doc.Text(),
	}
	distTempTyp := reflect.TypeOf(data.DistrictsTemperature)
	distTempVal := reflect.ValueOf(&data.DistrictsTemperature).Elem()

	// for better error reporting
	parseErrors := make([]error, 0, 20)

	// parse Temperature table
	doc.Find("table tr").Each(func(i int, s *goquery.Selection) {
		text1 := s.Find("td:nth-child(1)").Text()
		text2 := s.Find("td:nth-child(2)").Text()
		district := reName.ReplaceAllString(text1, "")
		_, ok := distTempTyp.FieldByName(district)

		if !ok {
			parseErrors = append(parseErrors, fmt.Errorf("[Warning] unknown district: %s", district))
			return
		}
		if !reDegree.MatchString(text2) {
			parseErrors = append(parseErrors, fmt.Errorf("[Error] unidentified degree string: %s (district: %s)", text2, district))
			return
		}

		field := distTempVal.FieldByName(district)
		submatches := reDegree.FindStringSubmatch(text2)
		degree, err := strconv.ParseInt(submatches[1], 10, 32)
		if err != nil {
			parseErrors = append(parseErrors, fmt.Errorf("[Error] unidentified degree number in string: %s (in %#v, district: %s)", submatches[1], text2, district))
			return
		}

		// set to the field
		field.Set(reflect.ValueOf(Temperature(degree)))
	})

	// parse air temperature
	descText := doc.Text()
	reAirTemp := regexp.MustCompile(`Air temperature\s*:\s*(\d+)\s+(degree|degrees) Celsius`)
	airTempStr := reAirTemp.FindStringSubmatch(descText)
	airTemp, _ := strconv.ParseInt(airTempStr[1], 10, 32)
	data.AirTemperature = Temperature(airTemp)

	// parse humidity
	reHumidity := regexp.MustCompile(`Relative Humidity\s*:\s*(\d+)\s+per cent`)
	humidityStr := reHumidity.FindStringSubmatch(descText)
	humidity, _ := strconv.ParseFloat(humidityStr[1], 64)
	data.RelativeHumidity = humidity / 100

	return
}
