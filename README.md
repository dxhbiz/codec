# golang编解码工具（结构体转换字节流，字节流转换结构体）

> A go codec project

## 示例

```go
package main

import (
	"github.com/dxhbiz/codec"
	"fmt"
)

type Interest struct {
	Name string
}

type Person struct {
	Name string
	Sex int8
	Interests []Interest
	Age int8
}

func main()  {
	personEncode := Person{
		Name: "test",
		Sex: 1,
		Interests: []Interest{
			{
				Name: "basketball",
			},
			{
				Name: "football",
			},
		},
		Age: 30,
	}
	buf, err := codec.Encode(&personEncode)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf)

	personDecode := Person{}
	err = codec.Decode(buf, &personDecode)
	fmt.Printf("%+v\n", personDecode)
}
```
