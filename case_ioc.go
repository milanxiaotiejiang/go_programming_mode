package main

import (
	"errors"
	"fmt"
)

/*********************************************** 结构体嵌入 */

type Widget struct {
	X, Y int
}

type Label struct {
	Widget        // Embedding (delegation)
	Text   string //Aggregation
}

type Button struct {
	Label // Embedding (delegation)
}

type ListBox struct {
	Widget          // Embedding (delegation)
	Texts  []string // Aggregation
	Index  int      // Aggregation
}

/*********************************************** 方法重写 */

type Painter interface {
	Paint()
}

type Clicker interface {
	Click()
}

/**
对于 Lable 来说，只有 Painter ，没有Clicker；对于 Button 和 ListBox来说，Painter 和Clicker都有
*/

func (label Label) Paint() {
	fmt.Printf("%p:Label.Paint(%q)\n", &label, label.Text)
}

// Paint 因为这个接口可以通过 Label 的嵌入带到新的结构体，
//所以，可以在 Button 中重载这个接口方法
func (button Button) Paint() { // Override
	fmt.Printf("Button.Paint(%s)\n", button.Text)
}
func (button Button) Click() {
	fmt.Printf("Button.Click(%s)\n", button.Text)
}

func (listBox ListBox) Paint() {
	fmt.Printf("ListBox.Paint(%q)\n", listBox.Texts)
}
func (listBox ListBox) Click() {
	fmt.Printf("ListBox.Click(%q)\n", listBox.Texts)
}

/*********************************************** 嵌入结构多态 */

/*********************************************** 反转控制 */

//type IntSet struct {
//	data map[int]bool
//}
//
//func NewIntSet() IntSet {
//	return IntSet{make(map[int]bool)}
//}
//func (set *IntSet) Add(x int) {
//	set.data[x] = true
//}
//func (set *IntSet) Delete(x int) {
//	delete(set.data, x)
//}
//func (set *IntSet) Contains(x int) bool {
//	return set.data[x]
//}

/*********************************************** 实现 Undo 功能 */

/**
我们在 UndoableIntSet 中嵌入了IntSet ，然后 Override 了 它的 Add()和 Delete() 方法；
Contains() 方法没有 Override，所以，就被带到 UndoableInSet 中来了。
在 Override 的 Add()中，记录 Delete 操作；
在 Override 的 Delete() 中，记录 Add 操作；
在新加入的 Undo() 中进行 Undo 操作。
*/

//type UndoableIntSet struct { // Poor style
//	IntSet    // Embedding (delegation)
//	functions []func()
//}
//
//func NewUndoableIntSet() UndoableIntSet {
//	return UndoableIntSet{NewIntSet(), nil}
//}
//
//func (set *UndoableIntSet) Add(x int) { // Override
//	if !set.Contains(x) {
//		set.data[x] = true
//		set.functions = append(set.functions, func() { set.Delete(x) })
//	} else {
//		set.functions = append(set.functions, nil)
//	}
//}
//
//func (set *UndoableIntSet) Delete(x int) { // Override
//	if set.Contains(x) {
//		delete(set.data, x)
//		set.functions = append(set.functions, func() { set.Add(x) })
//	} else {
//		set.functions = append(set.functions, nil)
//	}
//}
//
//func (set *UndoableIntSet) Undo() error {
//	if len(set.functions) == 0 {
//		return errors.New("no functions to undo")
//	}
//	index := len(set.functions) - 1
//	if function := set.functions[index]; function != nil {
//		function()
//		set.functions[index] = nil // For garbage collection
//	}
//	set.functions = set.functions[:index]
//	return nil
//}

/*********************************************** 反转依赖 */

/**
声明一种函数接口，表示我们的 Undo 控制可以接受的函数签名是什么样的
*/

type Undo []func()

/**
Undo 控制逻辑就可以写成下面这样
*/

func (undo *Undo) Add(function func()) {
	*undo = append(*undo, function)
}

func (undo *Undo) Undo() error {
	functions := *undo
	if len(functions) == 0 {
		return errors.New("no functions to undo")
	}
	index := len(functions) - 1
	if function := functions[index]; function != nil {
		function()
		functions[index] = nil // For garbage collection
	}
	*undo = functions[:index]
	return nil
}

/**
在 IntSet 里嵌入 Undo，接着在 Add() 和 Delete() 里使用刚刚的方法
*/

type IntSet struct {
	data map[int]bool
	undo Undo
}

func NewIntSet() IntSet {
	return IntSet{data: make(map[int]bool)}
}

func (set *IntSet) Undo() error {
	return set.undo.Undo()
}

func (set *IntSet) Contains(x int) bool {
	return set.data[x]
}

func (set *IntSet) Add(x int) {
	if !set.Contains(x) {
		set.data[x] = true
		set.undo.Add(func() {
			set.Delete(x)
		})
	} else {
		set.undo.Add(nil)
	}
}

func (set *IntSet) Delete(x int) {
	if set.Contains(x) {
		delete(set.data, x)
		set.undo.Add(func() { set.Add(x) })
	} else {
		set.undo.Add(nil)
	}
}

func main() {
	//结构体嵌入
	label := Label{Widget{10, 10}, "State:"}

	label.X = 11
	label.Y = 12

	//嵌入结构多态
	button1 := Button{Label{Widget{10, 70}, "OK"}}
	button2 := Button{Label{Widget{50, 70}, "Cancel"}}
	listBox := ListBox{Widget{10, 40},
		[]string{"AL", "AK", "AZ", "AR"}, 0}

	for _, painter := range []Painter{label, listBox, button1, button2} {
		painter.Paint()
	}

	for _, widget := range []interface{}{label, listBox, button1, button2} {
		widget.(Painter).Paint()
		if clicker, ok := widget.(Clicker); ok {
			clicker.Click()
		}
		fmt.Println() // print a empty line
	}

	set := NewIntSet()
	set.Add(1)
	set.Undo()
}
