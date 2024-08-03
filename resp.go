package main

import (
	"bufio"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

type Resp struct {
	reader *bufio.Reader
}

func newResp(rd io.Reader) *Resp {
	return &Resp{
		reader: bufio.NewReader(rd),
	}
}

func (r *Resp) readLine() (line []byte, n int, err error) {

	// Read until you get to \r
	b, err := r.reader.ReadBytes('\r')

	if err != nil {
		return nil, 0, err
	}
	// Read the \n as well
	_, err = r.reader.ReadByte()
	if err != nil {
		return nil, 0, err
	}
	return b[0 : len(b)-1], len(b) + 1, nil
}

func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}

	x64, err := strconv.ParseInt(string(line), 0, 64)
	if err != nil {
		return 0, 0, err
	}

	return int(x64), n, err

}

func (r *Resp) Read() (Value, error) {
	var val Value
	valType, err := r.reader.ReadByte()
	if err != nil {
		return val, err
	}

	switch valType {

	case BULK:
		return r.readBulk()
	case ARRAY:
		return r.readArray()
	default:
		return val, nil
	}

}

func (r *Resp) readBulk() (Value, error) {
	var val Value
	_, _, err := r.readInteger()

	if err != nil {
		return val, err
	}
	line, _, err := r.readLine()
	if err != nil {
		return val, err
	}
	val.typ = "bulk"
	val.bulk = string(line)
	return val, nil
}

func (r *Resp) readArray() (Value, error) {
	var val Value
	_, n, err := r.readInteger()

	if err != nil {
		return val, nil
	}
	val.typ = "array"
	val.array = make([]Value, 0)
	i := 0
	for i < n {
		readVal, err := r.Read()
		if err != nil {
			return val, nil
		}
		// Append read value to array
		val.array = append(val.array, readVal)
		i++
	}

	return val, nil
}
