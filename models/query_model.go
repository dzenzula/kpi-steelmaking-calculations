package models

type Query struct {
	IdMeasuring int
	TimeStamp   string
	Value       *string
	Quality     int
	BatchId     string
}
