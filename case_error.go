package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
)

type Point struct {
	Longitude     float64
	Latitude      float64
	Distance      float64
	ElevationGain float64
	ElevationLoss float64
}

//func parse(r io.Reader) (*Point, error) {
//
//	var p Point
//
//	if err := binary.Read(r, binary.BigEndian, &p.Longitude); err != nil {
//		return nil, err
//	}
//	if err := binary.Read(r, binary.BigEndian, &p.Latitude); err != nil {
//		return nil, err
//	}
//	if err := binary.Read(r, binary.BigEndian, &p.Distance); err != nil {
//		return nil, err
//	}
//	if err := binary.Read(r, binary.BigEndian, &p.ElevationGain); err != nil {
//		return nil, err
//	}
//	if err := binary.Read(r, binary.BigEndian, &p.ElevationLoss); err != nil {
//		return nil, err
//	}
//	return &p, nil
//}

/**
要解决这个事，我们可以用函数式编程的方式

通过使用 Closure 的方式把相同的代码给抽出来重新定义一个函数，这样大量的 if err!=nil 处理得很干净了，但是会带来一个问题，那就是有一个 err 变量和一个内部的函数，感觉不是很干净
*/
//func parse(r io.Reader) (*Point, error) {
//	var p Point
//	var err error
//	read := func(data interface{}) {
//		if err != nil {
//			return
//		}
//		err = binary.Read(r, binary.BigEndian, data)
//	}
//	read(&p.Longitude)
//	read(&p.Latitude)
//	read(&p.Distance)
//	read(&p.ElevationGain)
//	read(&p.ElevationLoss)
//	if err != nil {
//		return &p, err
//	}
//	return &p, nil
//}
/**
Go 语言的 bufio.Scanner()
1. 定义一个结构体和一个成员函数
2. 改变代码结构
*/
type Reader struct {
	r   io.Reader
	err error
}

func (r Reader) read(data interface{}) {
	if r.err == nil {
		r.err = binary.Read(r.r, binary.BigEndian, data)
	}
}

func parse(input io.Reader) (*Point, error) {

	var p Point
	r := Reader{r: input}

	r.read(&p.Longitude)
	r.read(&p.Latitude)
	r.read(&p.Distance)
	r.read(&p.ElevationGain)
	r.read(&p.ElevationLoss)

	if r.err != nil {
		return nil, r.err
	}

	return &p, nil
}

/**
在 Go 语言的开发者中，更为普遍的做法是将错误包装在另一个错误中，同时保留原始内容
*/
type authorizationError struct {
	operation string
	err       error // original error
}

func (e *authorizationError) Error() string {
	return fmt.Sprintf("authorization failed during %s: %v", e.operation, e.err)
}

func main() {
	/*********************************************** 错误处理 */
	/**
	Go 语言的很多函数都会返回 result、err 两个值，于是就有这样几点：
	1. 参数上基本上就是入参，而返回接口把结果和错误分离，这样使得函数的接口语义清晰；
	2. 而且，Go 语言中的错误参数如果要忽略，需要显式地忽略，用 _ 这样的变量来忽略；
	3. 另外，因为返回的 error 是个接口（其中只有一个方法 Error()，返回一个 string ），所以你可以扩展自定义的错误处理。
	4. 另外，如果一个函数返回了多个不同类型的 error，你也可以使用下面这样的方式：

	Go 语言的错误处理的方式，本质上是返回值检查，但是它也兼顾了异常的一些好处——对错误的扩展。
	*/

	err := fmt.Errorf("")

	if err != nil {
		switch err.(type) {
		case *json.SyntaxError:
			//
		//case *ZeroDivisionError:
		//	//
		//case *NullPointerError:
		//	//
		default:
			//
		}
	}

	/*********************************************** 资源清理 */
	/**
	C 语言：使用的是 goto fail; 的方式到一个集中的地方进行清理
	C++ 语言：一般来说使用 RAII 模式，通过面向对象的代理模式，把需要清理的资源交给一个代理类，然后再析构函数来解决。
	Java 语言：可以在 finally 语句块里进行清理。
	Go 语言：使用 defer 关键词进行清理。
	*/

	/*********************************************** 包装错误 */
	if err != nil {
		err = fmt.Errorf("something failed: %v", err)
	}

	/**
	https://github.com/pkg/errors
	*/
	if err != nil {
		err = errors.Wrap(err, "read failed")
	}
}
