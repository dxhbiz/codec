package codec

import (
	"bytes"
	"runtime"
	"reflect"
	"fmt"
	"encoding/binary"
)

func Decode(data []byte, v interface {}) error {
	de := new(decoder)
	de.buf = bytes.NewBuffer(data)
	if err := de.unmarshal(v); err != nil {
		return err
	}
	return nil
}

type decoder struct {
	buf *bytes.Buffer
}

func (this *decoder) unmarshal(v interface {}) (err error) {
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

func (this *decoder) reflect(v interface {}) {
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		this.error(fmt.Errorf("Interface's kind must be reflect.Ptr."))
	}
	rv := reflect.ValueOf(v).Elem()
	if !rv.CanSet() {
		this.error(fmt.Errorf("Reflect value can't set."))
	}
	rvCount := rv.NumField()
	for i := 0; i < rvCount; i++ {
		this.decode(rv.Field(i))
	}
}

func (this *decoder) decode(rv reflect.Value) {
	switch rv.Kind() {
	case reflect.Array:
		this.decodeArray(rv)
	case reflect.Slice:
		this.decodeSlice(rv)
	case reflect.Struct:
		this.decodeStruct(rv)
	default:
		this.decodeValue(rv)
	}
}

func (this *decoder) decodeArray(rv reflect.Value) {
	rvLen := rv.Len()
	if rvLen <= 0 {
		return
	}

	firstElem := rv.Index(0)
	switch firstElem.Kind() {
	case reflect.Uint8:
		valArr := this.buf.Next(rvLen)
		reflect.Copy(rv, reflect.ValueOf(valArr))
	case reflect.Struct:
		for i := 0; i < rvLen; i++ {
			this.decodeStruct(rv.Index(i))
		}
	}
}

func (this *decoder) decodeSlice(rv reflect.Value) {
	var rvLen uint32
	binary.Read(this.buf, binary.LittleEndian, &rvLen)

	if rvLen <= 0 {
		return
	}

	newRv := reflect.MakeSlice(rv.Type(), int(rvLen), int(rvLen))
	rv.Set(newRv)
	this.decodeArray(rv)
}

func (this *decoder) decodeStruct(rv reflect.Value) {
	rvCount := rv.NumField()
	for i := 0; i < rvCount; i++ {
		this.decode(rv.Field(i))
	}
}

func (this *decoder) decodeValue(rv reflect.Value) {
	switch rv.Kind() {
	case reflect.Int8:
		var val int8
		binary.Read(this.buf, binary.LittleEndian, &val)
		rv.SetInt(int64(val))
	case reflect.Uint8:
		var val uint8
		binary.Read(this.buf, binary.LittleEndian, &val)
		rv.SetUint(uint64(val))
	case reflect.Int16:
		var val int16
		binary.Read(this.buf, binary.LittleEndian, &val)
		rv.SetInt(int64(val))
	case reflect.Uint16:
		var val uint16
		binary.Read(this.buf, binary.LittleEndian, &val)
		rv.SetUint(uint64(val))
	case reflect.Int32:
		var val int32
		binary.Read(this.buf, binary.LittleEndian, &val)
		rv.SetInt(int64(val))
	case reflect.Uint32:
		var val uint32
		binary.Read(this.buf, binary.LittleEndian, &val)
		rv.SetUint(uint64(val))
	case reflect.Float32:
		var val float32
		binary.Read(this.buf, binary.LittleEndian, &val)
		rv.SetFloat(float64(val))
	case reflect.Int64:
		var val int64
		binary.Read(this.buf, binary.LittleEndian, &val)
		rv.SetInt(val)
	case reflect.Uint64:
		var val uint64
		binary.Read(this.buf, binary.LittleEndian, &val)
		rv.SetUint(val)
	default:
		this.error(fmt.Errorf("Can't decode %v type", rv.Kind()))
	}
}

func (this *decoder) error(err error) {
	panic(err)
}
