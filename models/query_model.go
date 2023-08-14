package models

type Query struct {
	IdMeasuring int
	TimeStamp   string
	Value       *float64
	Quality     int
	BatchId     int
}
