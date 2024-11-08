package main

import (
	"strconv"
)
type Output struct {
	Value int
	Index int
	Script string
}

func NewOutput(value int, index int, script string) *Output {
	return &Output{
		Value: value, 
		Index: index, 
		Script: script,
	}
}

func (output *Output) toString() string {
	str := strconv.Itoa(output.Value) + strconv.Itoa(output.Index) + output.Script
	return str
}

func outputsString(outputs []*Output) string {
	str := ""
	for i:=0; i<len(outputs); i++ {
		str += outputs[i].toString()
	}
	return str
}