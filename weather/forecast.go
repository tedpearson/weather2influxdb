package weather

import (
	"github.com/iancoleman/strcase"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/tedpearson/weather2influxdb/http"
	"reflect"
	"strconv"
	"time"
)


type WriteOptions struct {
	ForecastSource  string
	MeasurementName string
	Location        string
	ForecastTime    *int64
}

type Record struct {
	Time                     time.Time
	Temperature              *float64
	Dewpoint                 *float64
	FeelsLike                *float64
	SkyCover                 *float64
	WindDirection            *int
	WindSpeed                *float64
	WindGust                 *float64
	PrecipitationProbability *float64
	PrecipitationAmount      *float64
	SnowAmount               *float64
	IceAmount                *float64
}

type Records struct {
	Values []Record
}

func (rs Records) ToPoints(options WriteOptions) []*write.Point {
	ps := make([]*write.Point, len(rs.Values))
	for i, r := range rs.Values {
		ps[i] = toPoint(r.Time, r, options)
	}
	return ps
}

type AstroEvent struct {
	Time                     time.Time
	SunUp                    *int
	MoonUp                   *int
	// this is hard to name. It's not "how bright is the moon" - it's "ratio of current moon phase to the full moon".
	FullMoonRatio            *float64
}

type AstroEvents struct {
	Values []AstroEvent
}

func (as AstroEvents) ToPoints(options WriteOptions) []*write.Point {
	ps := make([]*write.Point, len(as.Values))
	for i, a := range as.Values {
		ps[i] = toPoint(a.Time, a, options)
	}
	return ps
}

func toPoint(t time.Time, i interface{}, options WriteOptions) *write.Point {
	e := reflect.ValueOf(i)
	p := influxdb2.NewPointWithMeasurement(options.MeasurementName).
		AddTag("source", options.ForecastSource).
		AddTag("location", options.Location).
		SetTime(t)
	if options.ForecastTime != nil {
		p.AddField("forecast_time", *options.ForecastTime)
		p.AddTag("forecast_time_tag", strconv.FormatInt(*options.ForecastTime, 10))
	}
	for i := 0; i < e.NumField(); i++ {
		name := strcase.ToSnake(e.Type().Field(i).Name)
		// note: skip time field already added above
		if name == "time" {
			continue
		}
		ptr := e.Field(i)
		if ptr.IsNil() {
			// don't dereference nil pointers
			continue
		}
		val := ptr.Elem().Interface()
		p.AddField(name, val)
	}
	return p
}

type Initer interface {
	Init(lat string, lon string, retryer http.Retryer) error
}

type Forecaster interface {
	Initer
	GetWeather() (Records, error)
}

type Astrocaster interface {
	Initer
	GetAstrocast() (AstroEvents, error)
}

func SetTemperature(r *Record, v float64) {
	r.Temperature = &v
}

func SetDewpoint(r *Record, v float64) {
	r.Dewpoint = &v
}

func SetFeelsLike(r *Record, v float64) {
	r.FeelsLike = &v
}

func SetSkyCover(r *Record, v float64) {
	r.SkyCover = &v
}

func SetWindDirection(r *Record, v float64) {
	i := int(v)
	r.WindDirection = &i
}

func SetWindSpeed(r *Record, v float64) {
	r.WindSpeed = &v
}

func SetWindGust(r *Record, v float64) {
	r.WindGust = &v
}

func SetPrecipitationProbability(r *Record, v float64) {
	r.PrecipitationProbability = &v
}

func SetPreciptationAmount(r *Record, v float64) {
	r.PrecipitationAmount = &v
}

func SetSnowAmount(r *Record, v float64) {
	r.SnowAmount = &v
}

func SetIceAmount(r *Record, v float64) {
	r.IceAmount = &v
}