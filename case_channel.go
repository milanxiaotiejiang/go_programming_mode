package main

import (
	"fmt"
	"sync"
)

/*********************************************** Channel 转发函数 */

func echo(nums []int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

/*********************************************** 平方函数 */

func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

/*********************************************** 过滤奇数函数 */

func odd(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			if n%2 != 0 {
				out <- n
			}
		}
		close(out)
	}()
	return out
}

/*********************************************** 求和函数 */

func sum(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		var sum = 0
		for n := range in {
			sum += n
		}
		out <- sum
		close(out)
	}()
	return out
}

/*********************************************** 代理函数 */

type EchoFunc func([]int) <-chan int
type PipeFunc func(<-chan int) <-chan int

func pipeline(nums []int, echo EchoFunc, pipeFns ...PipeFunc) <-chan int {
	ch := echo(nums)
	for i := range pipeFns {
		ch = pipeFns[i](ch)
	}
	return ch
}

/*********************************************** Fan in/Out */
/**
动用 Go 语言的 Go Routine 和 Channel 还有一个好处，就是可以写出 1 对多，或多对 1 的 Pipeline，也就是 Fan In/ Fan Out。
*/

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func is_prime(value int) bool {
	//for i := 2; i <= int(math.Floor(float64(value)/2)); i++ {
	//	if value%i == 0 {
	//		return false
	//	}
	//}
	//return value > 1
	return value%2 == 0
}

func prime(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			if is_prime(n) {
				out <- n
			}
		}
		close(out)
	}()
	return out
}

func merge(cs []<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)
	println("cs ", len(cs))
	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan int) {
			for n := range c {
				out <- n
			}
			wg.Done()
		}(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	//var nums = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	//for n := range sum(sq(odd(echo(nums)))) {
	//	fmt.Println(n)
	//}
	//
	///**
	//类似于我们执行了 Unix/Linux 命令： echo $nums | sq | sum
	//*/
	//for n := range pipeline(nums, echo, odd, sq, sum) {
	//	fmt.Println(n)
	//}

	/**
	1. 首先，我们制造了从 1 到 10000 的数组；
	2. 然后，把这堆数组全部 echo到一个 Channel 里—— in；
	3. 此时，生成 5 个 Channel，接着都调用 sum(prime(in)) ，于是，每个 Sum 的 Go Routine 都会开始计算和；
	4. 最后，再把所有的结果再求和拼起来，得到最终的结果。
	*/
	nums := makeRange(1, 10)
	for _, num := range nums {
		println(num, is_prime(num))
	}
	in := echo(nums)

	var chans [5]<-chan int
	for i := range chans {
		ints := prime(in)
		chans[i] = sum(ints)
	}
	for n := range sum(merge(chans[:])) {
		fmt.Println(n)
	}

	//var chans1 <-chan int
	//chans1 = sum(prime(in))
	//chanSum := sum(chans1)
	//for i := range chanSum {
	//	fmt.Println(i)
	//}

}
