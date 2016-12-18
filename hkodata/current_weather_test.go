package hkodata_test

import (
	"os"
	"testing"
	"time"

	"github.com/yookoala/weatherhk/hkodata"
)

func TestCurrentWeather(t *testing.T) {
	file, err := os.Open("./test/CurrentWeather.201612172144.xml")
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	defer file.Close()

	cw, err := hkodata.DecodeCurrentWeather(file)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	HKT, _ := time.LoadLocation("Asia/Hong_Kong")
	if want, have := time.Date(2016, time.December, 17, 21, 2, 0, 0, HKT), cw.PubDate; !want.Equal(have) {
		t.Errorf("expected %s, got %s", want, have)
	}
	if want, have := HKT.String(), cw.PubDate.Location().String(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := hkodata.Temperature(18), cw.AirTemperature; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := float64(0.71), cw.RelativeHumidity; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := hkodata.Temperature(18), cw.DistrictsTemperature.HongKongObservatory; want != have {
		t.Errorf("expected %d, got %d", want, have)
	}
	if want, have := hkodata.Temperature(17), cw.DistrictsTemperature.KingsPark; want != have {
		t.Errorf("expected %d, got %d", want, have)
	}
	if want, have := hkodata.Temperature(18), cw.DistrictsTemperature.WongChukHang; want != have {
		t.Errorf("expected %d, got %d", want, have)
	}
	if want, have := hkodata.Temperature(16), cw.DistrictsTemperature.TaKwuLing; want != have {
		t.Errorf("expected %d, got %d", want, have)
	}
	if want, have := hkodata.Temperature(17), cw.DistrictsTemperature.LauFauShan; want != have {
		t.Errorf("expected %d, got %d", want, have)
	}
	if want, have := hkodata.Temperature(17), cw.DistrictsTemperature.TaiPo; want != have {
		t.Errorf("expected %d, got %d", want, have)
	}
	if want, have := hkodata.Temperature(18), cw.DistrictsTemperature.ShaTin; want != have {
		t.Errorf("expected %d, got %d", want, have)
	}
	if want, have := hkodata.Temperature(18), cw.DistrictsTemperature.TuenMun; want != have {
		t.Errorf("expected %d, got %d", want, have)
	}
	if want, have := hkodata.Temperature(17), cw.DistrictsTemperature.TseungKwanO; want != have {
		t.Errorf("expected %d, got %d", want, have)
	}
	if want, have := hkodata.Temperature(18), cw.DistrictsTemperature.SaiKung; want != have {
		t.Errorf("expected %d, got %d", want, have)
	}
	if want, have := hkodata.Temperature(17), cw.DistrictsTemperature.CheungChau; want != have {
		t.Errorf("expected %d, got %d", want, have)
	}
	if want, have := hkodata.Temperature(19), cw.DistrictsTemperature.ChekLapKok; want != have {
		t.Errorf("expected %d, got %d", want, have)
	}
	if want, have := hkodata.Temperature(18), cw.DistrictsTemperature.ShekKong; want != have {
		t.Errorf("expected %d, got %d", want, have)
	}
	if want, have := hkodata.Temperature(16), cw.DistrictsTemperature.TsuenWanHoKoon; want != have {
		t.Errorf("expected %d, got %d", want, have)
	}
	if want, have := hkodata.Temperature(17), cw.DistrictsTemperature.TsuenWanShingMunValley; want != have {
		t.Errorf("expected %d, got %d", want, have)
	}
	if want, have := hkodata.Temperature(18), cw.DistrictsTemperature.HongKongPark; want != have {
		t.Errorf("expected %d, got %d", want, have)
	}
	if want, have := hkodata.Temperature(18), cw.DistrictsTemperature.ShauKeiWan; want != have {
		t.Errorf("expected %d, got %d", want, have)
	}
	if want, have := hkodata.Temperature(17), cw.DistrictsTemperature.KowloonCity; want != have {
		t.Errorf("expected %d, got %d", want, have)
	}
	if want, have := hkodata.Temperature(19), cw.DistrictsTemperature.HappyValley; want != have {
		t.Errorf("expected %d, got %d", want, have)
	}
	if want, have := hkodata.Temperature(18), cw.DistrictsTemperature.WongTaiSin; want != have {
		t.Errorf("expected %d, got %d", want, have)
	}
	if want, have := hkodata.Temperature(18), cw.DistrictsTemperature.Stanley; want != have {
		t.Errorf("expected %d, got %d", want, have)
	}
	if want, have := hkodata.Temperature(18), cw.DistrictsTemperature.KwunTong; want != have {
		t.Errorf("expected %d, got %d", want, have)
	}
	if want, have := hkodata.Temperature(18), cw.DistrictsTemperature.ShamShuiPo; want != have {
		t.Errorf("expected %d, got %d", want, have)
	}
	if want, have := hkodata.Temperature(19), cw.DistrictsTemperature.KaiTakRunwayPark; want != have {
		t.Errorf("expected %d, got %d", want, have)
	}
	if want, have := hkodata.Temperature(18), cw.DistrictsTemperature.YuenLongPark; want != have {
		t.Errorf("expected %d, got %d", want, have)
	}
}
