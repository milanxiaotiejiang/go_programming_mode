package main

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"time"
)

type Person struct {
	Name   string
	Sexual string
	Age    int
}

func PrintPerson(p *Person) {
	fmt.Printf("Name=%s, Sexual=%s, Age=%d\n", p.Name, p.Sexual, p.Age)
}

func (p *Person) Print() {
	fmt.Printf("Name=%s, Sexual=%s, Age=%d\n", p.Name, p.Sexual, p.Age)
}

//type WithName struct {
//	Name string
//}

type Country struct {
	//WithName
	Name string
}
type City struct {
	//WithName
	Name string
}

type Stringable interface {
	ToString() string
}

//type Printable interface {
//	PrintStr()
//}

//func (c Country) PrintStr() {
//	fmt.Println(c.Name)
//}

//func (c City) PrintStr() {
//	fmt.Println(c.Name)
//}

//func (w WithName) PrintStr() {
//	fmt.Println(w.Name)
//}

func (c Country) ToString() string {
	return "Country = " + c.Name
}

func (c City) ToString() string {
	return "City = " + c.Name
}

func PrintStr(p Stringable) {
	fmt.Println(p.ToString())
}

type Shape interface {
	Sides() int
	Area() int
}

type Square struct {
	len int
}

func (s *Square) Sides() int {
	return 4
}

//func (s *Square) Area() int {
//	return 6
//}

func main() {
	/*********************************************** Slice */
	/**
	首先，创建一个 foo 的 Slice，其中的长度和容量都是 5；
	然后，开始对 foo 所指向的数组中的索引为 3 和 4 的元素进行赋值；
	最后，对 foo 做切片后赋值给 bar，再修改 bar[1]。

	因为 foo 和 bar 的内存是共享的，所以，foo 和 bar 对数组内容的修改都会影响到对方。
	*/
	foo := make([]int, 5)
	foo[3] = 42
	foo[4] = 100

	bar := foo[1:4]
	bar[1] = 99

	/**
	把 a[1:16] 的切片赋给 b，此时，a 和 b 的内存空间是共享的
	对 a 做了一个 append()的操作，这个操作会让 a 重新分配内存，这就会导致 a 和 b 不再共享

	append()操作让 a 的容量变成了 64，而长度是 33。
	这里你需要重点注意一下，append()这个函数在 cap 不够用的时候，就会重新分配内存以扩大容量，如果够用，就不会重新分配内存了！
	*/
	a := make([]int, 32)
	b := a[1:16]
	a = append(a, 1)
	a[2] = 42

	log.Println(b)

	/**
	dir1 和 dir2 共享内存，虽然 dir1 有一个 append() 操作，但是因为 cap 足够，于是数据扩展到了dir2 的空间

	如果要解决这个问题
	dir1 := path[:sepIndex]  --> dir1 := path[:sepIndex:sepIndex]
	*/
	path := []byte("AAAA/BBBBBBBBB")
	sepIndex := bytes.IndexByte(path, '/')

	dir1 := path[:sepIndex]
	dir2 := path[sepIndex+1:]

	fmt.Println("dir1 =>", string(dir1)) //prints: dir1 => AAAA
	fmt.Println("dir2 =>", string(dir2)) //prints: dir2 => BBBBBBBBB

	dir1 = append(dir1, "suffix"...)

	fmt.Println("dir1 =>", string(dir1)) //prints: dir1 => AAAAsuffix
	fmt.Println("dir2 =>", string(dir2)) //prints: dir2 => uffixBBBB

	/*********************************************** 深度比较 */
	/**
	当我们复制一个对象时，这个对象可以是内建数据类型、数组、结构体、Map……
	在复制结构体的时候，如果我们需要比较两个结构体中的数据是否相同，就要使用深度比较，而不只是简单地做浅度比较。
	这里需要使用到反射 reflect.DeepEqual()
	*/
	m1 := map[string]string{"one": "a", "two": "b"}
	m2 := map[string]string{"two": "b", "one": "a"}
	fmt.Println("m1 == m2:", reflect.DeepEqual(m1, m2))

	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}
	fmt.Println("s1 == s2:", reflect.DeepEqual(s1, s2))

	/*********************************************** 接口编程 */
	/**
	其中一个使用一个函数，另一个使用一个“成员函数”
	使用“成员函数”的方式叫“Receiver”，这种方式是一种封装
	*/
	var p = Person{
		Name:   "Hao Chen",
		Sexual: "Male",
		Age:    44,
	}

	PrintPerson(&p)
	p.Print()

	/**
	使用了一个 Printable 的接口，而 Country 和 City 都实现了接口方法 PrintStr() 把自己输出
	*/
	//c1 := Country{"China"}
	//c2 := City{"Beijing"}
	//c1.PrintStr()
	//c2.PrintStr()
	/**
	引入一个叫 WithName的结构体
	*/
	//c1 := Country{WithName{"China"}}
	//c2 := City{WithName{"Beijing"}}
	//c1.PrintStr()
	//c2.PrintStr()
	/**
	我们使用了一个叫Stringable 的接口，我们用这个接口把“业务类型” Country 和 City 和“控制逻辑” Print() 给解耦了。
	于是，只要实现了Stringable 接口，都可以传给 PrintStr() 来使用。
	*/
	d1 := Country{"USA"}
	d2 := City{"Los Angeles"}
	PrintStr(d1)
	PrintStr(d2)

	/*********************************************** 接口完整性检查 */
	/**
	Square 并没有实现 Shape 接口的所有方法，程序虽然可以跑通，但是这样的编程方式并不严谨

	*/
	s := Square{len: 5}
	fmt.Printf("%d\n", s.Sides())

	//var _ Shape = (*Square)(nil)

	/*********************************************** 时间 */
	/**
	一定要使用 time.Time 和 time.Duration 这两个类型
	*/
	println(time.Time{})

	/*********************************************** 性能提示 */
	/**
	1. 如果需要把数字转换成字符串，使用 strconv.Itoa() 比 fmt.Sprintf() 要快一倍左右。
	2. 尽可能避免把String转成[]Byte ，这个转换会导致性能下降。
	3. 如果在 for-loop 里对某个 Slice 使用 append()，请先把 Slice 的容量扩充到位，这样可以避免内存重新分配以及系统自动按 2 的 N 次方幂进行扩展但又用不到的情况，从而避免浪费内存。
	4. 使用StringBuffer 或是StringBuild 来拼接字符串，性能会比使用 + 或 +=高三到四个数量级。
	5. 尽可能使用并发的 goroutine，然后使用 sync.WaitGroup 来同步分片操作。
	6. 避免在热代码中进行内存分配，这样会导致 gc 很忙。尽可能使用 sync.Pool 来重用对象。
	7. 使用 lock-free 的操作，避免使用 mutex，尽可能使用 sync/Atomic包
	8. 使用 I/O 缓冲，I/O 是个非常非常慢的操作，使用 bufio.NewWrite() 和 bufio.NewReader() 可以带来更高的性能。
	9. 对于在 for-loop 里的固定的正则表达式，一定要使用 regexp.Compile() 编译正则表达式。性能会提升两个数量级。
	10. 如果你需要更高性能的协议，就要考虑使用 protobuf 或 msgp 而不是 JSON，因为 JSON 的序列化和反序列化里使用了反射。
	11. 你在使用 Map 的时候，使用整型的 key 会比字符串的要快，因为整型比较比字符串比较要快。
	*/
}
