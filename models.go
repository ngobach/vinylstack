package main

type song struct {
	url    string
	title  string
	artist string
	album  string
	cover  string
}

type doc struct {
	playlist []song
}
