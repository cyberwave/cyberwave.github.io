[原文](https://medium.com/kokster/go-reflection-creating-objects-from-types-part-ii-composite-types-69a0e8134f20)
Go反射：根据类型创建对象-第二部分(复合类型)

> 本文是用Golang根据类型创建对象系列两部分博客的第二部分  
> 该部分讨论创建复合类型，第一部分的博客[传送门](https://github.com/cyberwave/cyberwave.github.io/blob/master/goblob/go-reflection-creating-objects-from-types-part-i-primitive-types.md)

![Golang Gopher](https://cdn-images-1.medium.com/max/1000/1*65iXGLup5igJDZXA1oOFXw.png)

在上一个博文中，我解释了 go reflect 包中 `type` 和 `kind` 的概念。在这个博文中，我将深入到Golang中术语的使用。这是因为在复合类型中术语 `type` 和 `kind` 承载比原型类型更多的意义。

# Types 和 Kinds
"Type" 是程序员用来描述程序中关于数据或函数中元数据的术语。Go 运行时和编译器对单词 `type` 有不同的意义。

这可以通过一个例子来明白。考虑 Go 程序中这个声明
```golang
var Version string
```
当提到这个声明，程序员会说 `Version` 是一个类型为 `string` 的变量。  
考虑下面的例子，其中 `Version` 是一个复合类型。
```golang
type VersionType string
var Version VersionType
```
在这个例子中，一个叫做 `VersionType` 的新类型在第一行被创建。在第二行中，`Version` 被声明为一个类型 `VersionType` 的变量。 它的类型的类型（`VersionType`)是 `string`。这些的名字是 `kind`。

总之，`Version` 的 `Type` 是 `VersionType`。`Version` 的 `kind` 是 `string`。

> *类型是 Go 程序中用户定义的关于数据或函数的元数据。Kind是 Go 程序中编译器定义的关于数据或函数的元数据。*

Kind 被运行时和编译器用来分配一个变量的内存布局或者为一个函数分配栈布局。

# 从原始类型创建复合对象
创建 `kind` 具有以下任何值的复合对象与从原始 Kinds 中创建原始对象没有区别。
```golang
    Bool
    Int
    Int8
    Int16
    Int32
    Int64
    Uint
    Uint8
    Uint16
    Uint32
    Uint64
    Uintptr
    Float32
    Float64
    Complex64
    Complex128
    String
    UnsafePointer
```
这是一个从它的类型（`VersionType`）创建 `Version` 的例子。由于 `Version` 有一个原始 `kind`，它可以通过使用零值创建。
```golang
1  func CreatePrimitiveObjects(t reflect.Type) reflect.Value {
2    return reflect.Zero(t)
3  }
4
5  func extractVersionType(v reflect.Value) (VersionType, error) {
6    if v.Type().String() != "VersionType" {  
7      return "", errors.New("invalid input")
8    }
9    return v.String(), nil
10 }
```
请注意，一个 `Type` 变量的 `String()` 方法返回完整的路径和 `Type` 的名称。即如果 `VersionType` 被定义在一个不同的包中 `mypkg`，那么这个 `String()` 方法会返回`mypkg.VersionType`。

# 从复合 Kinds 创建复合对象
复合种类是包含其他 Kinds的 Kinds。`Map`， `Struct`， `Array` 等都是复合类型(Kinds)的例子。这里是复合种类(Kinds)的完整列表。
```golang
    Array
    Chan
    Func
    Interface
    Map
    Ptr
    Slice
    Struct
```
复合对象可以使用“零”值创建，就像基本对象一样。但是，如果不进行额外的处理，则无法使用特定类型的空值进行初始化。后面的部分提供了关于初始化复合类型(Kinds)的更多细节。

# 从他们的类型签名创建组合数组对象
空数组对象可以使用零值创建。数组的零值是一个空的数组对象。这是一个从它的类型签名创建一个数据对象的例子
```golang
func CreateCompositeObjects（t reflect.Type）reflect.Value { 
    return reflect.Zero（t）
}
```
这个函数将创建一个持有任何复合类型的空对象的 `reflect.Value` 结构体。

反射包有一个 `ArrayOf(Type, int)` 函数，可以用来创建包含类型 `Type` 和指定长度的元素的数组。这是一个显示其用法的例子。
```golang
func CreateArray（t reflect.Type，length int）reflect.Value { 
  var arrayType reflect.Type 
  arrayType = reflect.ArrayOf（length，t）
  return reflect.Zero（arrayType）
}
```
数组元素的类型(type)由传入的参数 `t` 决定。 数组的长度由传入的参数 `length` 决定。为了将值解析到数组对象中，使用反射包的最佳选项是将该数组作为接口进行处理。
```golang
func extractArray(v reflect.Value) (interface{}, error) {
  if v.Kind() != reflect.Array {
    return nil, errors.New("invalid input")
  }
  var array interface{}
  array = v.Interface()
  return array, nil
}
```
请注意，值的 `Slice()` 方法也可以用来解析数组的值。但是，此方法需要在 `reflect.Value` 计算数组之前将数组类型转换为可寻址数组。避免这个细微之处最好使用 `Interface` 方法来保持代码基础的简单性和可读性。

# 从他们的类型签名创建复合 Channel 对象
空 channel 对象可以使用零值创建。通道的零值是空的通道对象。以下是从其类型签名创建通道对象的示例
```golang
func CreateCompositeObjects（t reflect.Type）reflect.Value { 
    return reflect.Zero（t）
}
```
这个函数将创建一个持有任何复合类型的空对象的 `reflect.Value` 的结构体。

反射包有两个函数来创建一个 chan。该 `ChanOf` 函数可用于创建通道的类型签名，然后 `MakeChan(Type, int)` 函数可为通道分配内存。这是一个显示其用法的例子。
```golang
func CreateChan（t reflect.Type，buffer int）reflect.Value { 
    chanType：= reflect.ChanOf（reflect.BothDir，t）
    return reflect.MakeChan（chanType，buffer）
}
```
通道元素的类型由传入参数 `t` 决定。频道的长度由传入的参数 `buffer` 决定。为了将值解析到通道对象中，使用反射包的最佳选项是将该通道作为接口进行处理。创建频道的方向可以通过将方向作为第一个参数传递给调用 `ChanOf` 来控制。可能的值 `ChanDir` 如下。
```golang
    SendDir 
    RecvDir 
    BothDir
```
`BothDir` 表示通道可以写入和读取。`SendDir` 表示频道只能写入。`RecvDir` 表示通道只能接收数据。
```golang
func extractChan（v reflect.Value）（interface {}，error）{ 
  if v.Kind（）！= reflect.Chan { 
    return nil，errors.New（“invalid input”）
  } 
  var ch interface {} 
  ch = v。 Interface（）
  return ch，nil 
}
```
# 从他们的类型签名中创建复合函数对象
函数对象**不能**使用零值创建。

反射包有两个函数来创建一个函数。该 `FuncOf` 函数可以用来创建函数的类型签名，然后 `MakeFunc(Type, func(args []Value) (results []Value)) Value` 函数可以用来为函数分配内存。这是一个显示其用法的例子。
```golang
func CreateFunc(
                fType reflect.Type, 
                f func(args []reflect.Value) 
                (results []reflect.Value)
               ) (reflect.Value, error) 
{ 
 if fType.Kind() != reflect.Func {
  return reflect.Value{}, errors.New("invalid input")
 }
 
 var ins, outs *[]reflect.Type 
 
 ins = new([]reflect.Type)
 outs = new([]reflect.Type)
 
 for i:=0; i<fType.NumIn(); i++ {
  *ins = append(*ins, fType.In(i))
 }
 
 for i:=0; i<fType.NumOut(); i++ {
  *outs = append(*outs, fType.Out(i))
 }
 var variadic bool
 variadic = fType.IsVariadic()
 return AllocateStackFrame(*ins, *outs, variadic, f), nil
}
func AllocateStackFrame(
                        ins []reflect.Type,
                        outs []reflect.Type,
                        variadic bool,
                        f func(args []reflect.Value) 
                        (results []reflect.Value)
                       ) reflect.Value 
{
 var funcType reflect.Type
 funcType = reflect.FuncOf(ins, outs, variadic)
 return reflect.MakeFunc(funcType, f) 
}
```
`CreateFunc` 需要2个参数。第一个是你想创建的函数的 `type`，第二个是实现第一个参数类型的实际函数，但是输入和返回值被转换为 `reflect.Value` 对象。我已经创建了这个函数来为用户提供一个方便的方法，以便他们不必重新发明函数是如何在运行中创建的。

为了使用它来创建函数，需要定义函数的类型签名。这是一个类型签名的例子
```golang
type fn func（int）（int）
```
上面的函数签名的Type结构体可以使用下面的方法构造。
```golang
1 var funcVar fn 
2 var funcType reflect.Type 
3 funcType = reflect.TypeOf（funcVar）
```
`funcType` 可以作为传入 `CreateFunc` 的第一个参数。需要将满足 `funcType` 类型签名的函数作为第二个参数传递。下面的例子说明了：
```golang
func doubler（input int）（int）{ 
return input * 2 
}
```
该函数 `doubler` 将输入乘以2.它接受1个类型为 `int` 的参数并返回1个类型为 `int` 参数。此函数满足 `fn` 类型，但不满足通用函数类型：
```golang
f func（args [] reflect.Value）（results [] reflect.Value）
```
为了满足通用的 func 类型，doubler应该创建一个满足通用func类型的等效版本。如下所示。
```golang
func doublerReflect(args []reflect.Value) (result []reflect.Value) {
 if len(args) != 1 {
  panic(fmt.Sprintf("expected 1 arg, found %d", len(args)))
 }
 if args[0].Kind() != reflect.Int {
  panic(fmt.Sprintf("expected 1 arg of kind int, found 1 args of kind", args[0].Kind()))
 }
 
 var intVal int64
 intVal = args[0].Int()
 
 var doubleIntVal int
 doubleIntVal = doubler(int(intVal))
 
 var returnValue reflect.Value
 returnValue = reflect.ValueOf(doubleIntVal)
 
 return []reflect.Value{returnValue}
}
```
`doublerReflect` 在功能上等同于 `doubler` 并满足通用func类型签名。即它需要1个参数，它是一个 `reflect.Value` 对象的片段，并返回1个值，这是一个 `reflect.Value` 对象的片段。输入表示函数的输入，返回值表示新函数的返回值。

函数调用 `CreateFunc` 可以用 `funcType` 和 `doublerReflect` 作为参数。下面的例子展示了一个由新创建的函数调用的函数。
```golang
func main（）{ 
  var funcVar fn 
  var funcType reflect.Type
  funcType = reflect.TypeOf（funcVar）
  v，err：= CreateFunc（funcType，doublerReflect）
   if err！= nil { 
     fmt.Printf（“Error function％v \ n”，err）
   } 
 
   input：= 42 
 
   ins：= [ ] reflect.Value（[] reflect.Value {reflect.ValueOf（input）}）
 
   outs：= v.Call（ins）
   for i：= range outs { 
     fmt.Printf（“％+ v \ n”，outs [i ] .Interface（））
   } 
} 
//输出：84
```
# 从他们的类型签名创建复合Map对象
空 Map 对象可以使用零值创建。Map的零值是一个空的 Map 对象。这是一个从它的类型签名创建一个Map对象的例子
```golang
func CreateCompositeObjects（t reflect.Type）reflect.Value { 
    return reflect.Zero（t）
}
```
这个函数将创建一个持有任何复合类型的空对象的 `reflect.Value` 结构体。

反射包有一个 `MapOf(Type, Type)` 函数，可以用来创建包含第一个参数的 `type`的键和第二个参数的元素的 `type`的Map。这是一个显示其用法的例子。
```golang
func CreateMap（key，elem reflect.Type）reflect.Value { 
  var mapType reflect.Type 
  mapType = reflect.MapOf（key，elem）
  return reflect.MakeMap（mapType）
}
```
为了将这些值解析到 map 对象中，使用反射包的最佳可能选项是将该 map 作为接口进行处理。
```golang
func extractMap(v reflect.Value) (interface{}, error) {
  if v.Kind() != reflect.Map {
    return nil, errors.New("invalid input")
  }
  var mapVal interface{}
  mapVal = v.Interface()
  return mapVal, nil
}
```
请注意，maps也可以通过使用 `MakeMapWithSize` 分配。使用此方法创建映射的步骤与上述步骤完全相同，但 `MakeMap` 调用可以使用 `MakeMapWithSize` 和一个 *size* 参数替代。

# 从他们的类型签名创建复合Ptr对象
空的Ptr对象可以通过使用零值创建。Ptr 的零值是一个 `nil` 指针对象。下面是从它的类型签名中创建一个Ptr对象的例子
```golang
func CreateCompositeObjects（t reflect.Type）reflect.Value { 
    return reflect.Zero（t）
}
```
这个函数将创建一个持有任何复合类型的空对象 `reflect.Value` 的结构体。

反射包有一个 `PtrTo(Type)` 函数，可以用来创建指向类型元素 `Type` 的Ptrs。这是一个显示其用法的例子。
```golang
func CreatePtr（t reflect.Type）reflect.Value { 
  var ptrType reflect.Type 
  ptrType = reflect.PtrTo（t）
  return reflect.Zero（ptrType）
}
```
*请注意，上述功能也可以通过使用 `reflect.New(t)` 来实现。上面说明的技术在功能上等同于使用 `reflect.New(t)`。*

Ptr指向的元素的类型由传入参数 `t` 决定。为了将值解析到对象中，使用反射包的最佳选项是首先使用 `Indirect()`调用或`Elem()` 方法来确定值，然后将该值作为接口进行处理。
```golang
func extractElement(v reflect.Value) (interface{}, error) {
  if v.Kind() != reflect.Ptr {
    return nil, errors.New("invalid input")
  }
  var elem reflect.Value
  elem = v.Elem()
  
  var elem interface{}
  elem = v.Interface()
  return elem, nil
}
```
# 从他们的类型签名创建复合切片对象
空切片对象可以通过使用零值创建。Slice的零值是一个空的Slice对象。这是一个从它的类型签名创建一个Slice对象的例子
```golang
func CreateCompositeObjects（t reflect.Type）reflect.Value { 
    return reflect.Zero（t）
}
```
这个函数将创建一个持有任何复合类型的空对象的 `reflect.Value` 结构体。

反射包有一个 `SliceOf(Type)` 函数，可以用来创建包含类型 `Type` 元素的切片。这是一个显示其用法的例子。
```golang
func CreateSlice（t reflect.Type）reflect.Value { 
  var sliceType reflect.Type 
  sliceType = reflect.SliceOf（length，t）
  return reflect.Zero（sliceType）
}
```
Slice的元素的类型由传入的参数 `t` 决定。为了将值解析到切片对象中，使用反射包的最佳选项是将该切片作为接口进行处理。
```golang
func extractSlice(v reflect.Value) (interface{}, error) {
  if v.Kind() != reflect.Slice {
    return nil, errors.New("invalid input")
  }
  var slice interface{}
  slice = v.Interface()
  return slice, nil
}
```
# 从他们的类型签名创建复合结构对象
空结构对象可以使用零值创建。结构体的零值是一个空的结构体对象。下面是从它的类型签名中创建一个结构体对象的例子
```golang
func CreateCompositeObjects（t reflect.Type）reflect.Value { 
    return reflect.Zero（t）
}
```
这个函数将创建一个持有任何复合类型的空对象的 `reflect.Value`结构体。

反射包有一个 `StructOf([]reflect.StructFields)` 函数可以用来创建包含 `StructField` 对象中定义的类型的字段的结构体 。这是一个显示其用法的例子。
```golang
func CreateStruct（fields [] reflect.StructField）reflect.Value { 
  var structType reflect.Type 
  structType = reflect.StructOf（fields）
  return reflect.Zero（structType）
}
```
为了将值解析到数组对象中，使用反射包的最佳选项是将该结构作为接口进行处理。
```golang
func extractStruct(v reflect.Value) (interface{}, error) {
  if v.Kind() != reflect.Struct {
    return nil, errors.New("invalid input")
  }
  var st interface{}
  st = v.Interface()
  return st, nil
}
```
# 结论
这是一个关于在运行时使用反射包来即时创建任何Go类型的完整教程。我提供了创建 `Func` 类型的便捷方法，因为它比其他类型更复杂，如果不仔细设计，很容易污染你的代码库。

请继续关注我的下一篇博客文章，在 Golang 中将任何类型转换为其他类型！接下来，我将解释如何在Golang中编写JIT（Just-in-Time），然后介绍如何在 Go中使用 `reflect` 代码生成。


