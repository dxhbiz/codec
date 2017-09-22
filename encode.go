package codec

import (
	"bytes"
	"runtime"
	"reflect"
	"encoding/binary"
	"fmt"
)

func Encode(v interface {}) ([]byte, error) {
	en := new(encoder)
	en.buf = new(bytes.Buffer)
	if err := en.marshal(v); err != nil {
		return nil, err
	}
	return en.buf.Bytes(), nil
}

func reflectValue(v interface {}) reflect.Value {
	var val reflect.Value
	if reflect.TypeOf(v).Kind() == reflect.Ptr {
		val = reflect.ValueOf(v).Elem()
	} else {
		val = reflect.ValueOf(v)
	}
	return val
}

type encoder struct {
	buf *bytes.Buffer
}

func (this *encoder) marshal(v interface {}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			if s, ok := r.(string); ok {
				panic(s)
			}
			err = r.(error)
		}
	}()
	this.reflect(v)
	return nil
}

func (this *encoder) reflect(v interface {}) {
	rv := reflectValue(v)
	rvCount := rv.NumField()

	for i := 0; i < rvCount; i++ {
		this.encode(rv.Field(i))
	}
}

func (this *encoder) encode(rv reflect.Value) {
	switch rv.Kind() {
	case reflect.Array:
		this.encodeArray(rv)
	case reflect.Struct:
		this.encodeStruct(rv)
	case reflect.Slice:
		this.encodeSlice(rv)
	default:
		this.encodeValue(rv)
	}
}

func (this *encoder) encodeArray(rv reflect.Value) {
	rvLen := rv.Len()

	if rvLen <= 0 {
		return
	}

	firstElem := rv.Index(0)
	switch firstElem.Kind() {
	case reflect.Struct:
		for i:= 0; i < rvLen; i++ {
			this.encodeStruct(rv.Index(i))
		}
	case reflect.Uint8:
		valArr := make([]byte, rvLen)
		for i := 0; i < rvLen; i++ {
			valArr[i] = byte(rv.Index(i).Uint())
		}
		this.buf.Write(valArr)
	}
}

func (this *encoder) encodeStruct(rv reflect.Value) {
	rvCount := rv.NumField()
	for i := 0; i < rvCount; i++ {
		this.encode(rv.Field(i))
	}
}

func (this *encoder) encodeSlice(rv reflect.Value) {
	rvLen := uint32(rv.Len())
	binary.Write(this.buf, binary.LittleEndian, rvLen)
	this.encodeArray(rv)
}

func (this *encoder) encodeValue(rv reflect.Value) {
	switch rv.Kind() {
	case reflect.Int8:
		val := int8(rv.Int())
		binary.Write(this.buf, binary.LittleEndian, val)
	case reflect.Uint8:
		val := uint8(rv.Uint())
		binary.Write(this.buf, binary.LittleEndian, val)
	case reflect.Int16:
		val := int16(rv.Int())
		binary.Write(this.buf, binary.LittleEndian, val)
	case reflect.Uint16:
		val := uint16(rv.Uint())
		binary.Write(this.buf, binary.LittleEndian, val)
	case reflect.Int32:
		val := int32(rv.Int())
		binary.Write(this.buf, binary.LittleEndian, val)
	case reflect.Uint32:
		val := uint32(rv.Uint())
		binary.Write(this.buf, binary.LittleEndian, val)
	case reflect.Float32:
		val := float32(rv.Float())
		binary.Write(this.buf, binary.LittleEndian, val)
	case reflect.Int64:
		val := rv.Int()
		binary.Write(this.buf, binary.LittleEndian, val)
	case reflect.Uint64:
		val := rv.Uint()
		binary.Write(this.buf, binary.LittleEndian, val)
	default:
		this.error(fmt.Errorf("Can't encode %v type.", rv.Kind()))
	}
}

func (this *encoder) error(err error) {
	panic(err)
}
