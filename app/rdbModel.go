package main

type RdbModel struct {
	header   string
	metadata map[string]string
	database *Database
}

var GlobalRDB = &RdbModel{}

type Database struct {
	dbNumber int
	HashMap  map[string]Value
}
