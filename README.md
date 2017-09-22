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
	codecInfo := Person{
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
	buf, err := codec.Encode(&codecInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf)

	info := Person{}
	err = codec.Decode(buf, &info)
	fmt.Printf("%+v\n", info)
}
```
