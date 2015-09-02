package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Point struct {
	Time  uint64
	Value float64
}

type Metric struct {
	//	Path      string
	StartTime uint64
	Points    []Point
}

type Partition struct {
	Key     string
	Metrics map[string]*Metric
}

type Ring struct {
	Size       uint32
	Partitions map[string]*Partition
}

func NewRing() *Ring {
	ring := new(Ring)
	ring.Size = 1

	ring.Partitions = make(map[string]*Partition)

	ring.MakePartitions()

	return ring
}

func (r *Ring) MakePartitions() {
	r.Partitions["main"] = &Partition{Key: "main"}
	r.Partitions["main"].Metrics = make(map[string]*Metric)
}

func (r *Ring) InsertPoint(path string, point *Point) {
	m := r.GetMetric(path)

	if m.Points == nil {
		m.Points = make([]Point, 0)
	}

	m.Points = append(m.Points, *point)
}

func (r *Ring) GetMetric(path string) *Metric {
	m := r.Partitions["main"].Metrics[path]

	if m == nil {
		m = &Metric{}
		r.Partitions["main"].Metrics[path] = m
	}

	return m
}

func main() {

	// ring := NewRing()

	// ring.InsertPoint("testpath", &Point{Time: 3, Value: 2.0})

	// fmt.Printf("%+v", ring)

	p1 := Point{Time: 255, Value: 1.0}
	p2 := Point{Time: 256, Value: 2.0}
	p3 := Point{Time: 257, Value: 3.0}

	m := Metric{Points: []Point{p1, p2, p3}}

	b := new(bytes.Buffer)

	if err := binary.Write(b, binary.BigEndian, m); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v", *b)
	}
}
