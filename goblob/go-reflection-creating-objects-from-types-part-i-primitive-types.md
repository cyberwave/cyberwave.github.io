[原文](https://medium.com/kokster/go-reflection-creating-objects-from-types-part-i-primitive-types-6119e3737f5d)

Go反射：根据类型创建对象-第一部分(基本类型)
> 本文是用Golang根据类型创建对象系列两部分博客的第一部分
> 该部分讨论创建基本类型
![Golang Gopher](https://cdn-images-1.medium.com/max/1000/1*65iXGLup5igJDZXA1oOFXw.png)

Golang的反射包提供了基本的API来改变程序的控制流，该程序基于将要处理的对象的类型。

反射包提供两个重要的结构 - `Type` 和 `Value` 。 `Type` 是任何Go类型的一种表现。即它可以被用来编码为任何Go类型(例如 `int`,`string`,`bool`,`myCustomType` 等等)。`Value` 是任何Go值的代表，即它可以被用来编码和操作任何的Go值。

## Types vs Kinds
Golang有一个 `type` 和 `kind` 间区别的隐藏的，鲜为人知的惯例。这个区别可以使用例子阐明。
考虑这个结构体：
```golang
type example struct {
    field1 type1
    field2 type2
}
```
这个结构体的一个对象的 `type` 将会是 `example`。这个对象的 `kind` 将会是 `struct`。`kind`可以被认为是 `type` 的 `type`。
> Golang 中所有的结构体是同样的 `kind`，但是不是同样的 `type`。

复合类型像 `Pointer`,`Array`,`Slice`,`Map` 等等。在 `type` 和 `kind` 是有区别的。

对照中，基本类型像 `int`,`float`,`string` 等等。在 `kind` 和 `type` 间是没有区别的。即一个 `int` 型变量的 `kind` 是 `int`，一个 `int` 类型变量的 `type` 也是 `int`。

## 从类型中创建对象
为了从一个类型签名中创建一个对象，对象的 `type` 和 `kind` 都是必需的。此后，当我使用术语 '类型签名'，我的意思是指Golang 中该类型的 `reflect.Type` 对象。

### 根据原始类型创建原始对象
原始对象可以使用它们的零值从它们的类型签名中创建。
> 一个类型的零值是这个类型未初始化对象的值

这儿是一个所有Golang中原始类型的列表
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
使用 `reflect.Zero`函数，原始类型的对象可以被创建
```golang
func CreatePrimitiveObjects(t reflect.Type) reflect.Value {
  return reflect.Zero(t)
}
```
这将创建一个预期类型的对象并返回一个 `reflect.Value` 相应底层零值的对象。为了这个对象起作用，它的值需要被解析。

对于每个原始类型，对象的值可以使用合适的方法被解析。

## 解析整型值
在Golang中有5个整型类型：
```golang
    Int
    Int8
    Int16
    Int32
    Int64
```
`Int` 类型表示由平台定义的默认大小的整数。下面4个类型是大小分别为8，16，32和64（bit）大小的整数。

为了解析每个整型类型，`reflect.Value` 对象代表需要被转换为合适类型的整数。

这儿是 `int32` 应该如何被解析的：
```golang
// Extract Int32
func extractInt32(v reflect.Value) (int32, error) {
  if reflect.Kind() != reflect.Int32 {
    return int32(0), errors.New("Invalid input")
  }
  var intVal int64
  intVal = v.Int()
  return int32(intVal), nil
}
```

注意 `reflect.Int()` 总是返回 `int64` 是重要的。这是为什么 `int64` 可以编码为所有其他整型的原因。

下面是剩下的整型类型如何被解析：
```golang
// Extract Int64
func extractInt64(v reflect.Value) (int64, error) {
  if reflect.Kind() != reflect.Int64 {
    return int64(0), errors.New("Invalid input")
  }
  var intVal int64
  intVal = v.Int()
  return intVal, nil
}
// Extract Int16
func extractInt16(v reflect.Value) (int16, error) {
  if reflect.Kind() != reflect.Int16 {
    return int16(0), errors.New("Invalid input")
  }
  var intVal int64
  intVal = v.Int()
  return int16(intVal), nil
}
// Extract Int8
func extractInt8(v reflect.Value) (int8, error) {
  if reflect.Kind() != reflect.Int8 {
    return int8(0), errors.New("Invalid input")
  }
  var intVal int64
  intVal = v.Int()
  return int8(intVal), nil
}
// Extract Int
func extractInt(v reflect.Value) (int, error) {
  if reflect.Kind() != reflect.Int {
    return int(0), errors.New("Invalid input")
  }
  var intVal int64
  intVal = v.Int()
  return int(intVal), nil
}
```

## 解析布尔值
布尔值由 reflect 包的常量 `Bool` 表示。
它们可以使用 `Bool()` 方法 `reflect.Value` 对象中被解析：
```golang
// Extract Bool
func extractBool(v reflect.Value) (bool, error) {
  if reflect.Kind() != reflect.Bool {
    return false, errors.New("Invalid input")
  }
  return v.Bool(), nil
}
```

## 解析无符号整数

在Golang中有5类无符号整数类型：
```golang
    Uint
    Uint8
    Uint16
    Uint32
    Uint64
```
`Uint` 类型由平台指定的默认大小的无符号整数。下面的 4 个类型分别是大小(bit)为8,16,32和64的无符号整数。

为了解析每个无符号整数类型。`reflect.Value` 对象表示需要被转换为合适类型的无符号整数。

下面是 `Uint32` 应该如何被解析：
```golang
// Extract Uint32
func extractUint32(v reflect.Value) (uint32, error) {
  if reflect.Kind() != reflect.Uint32 {
    return uint32(0), errors.New("Invalid input")
  }
  var uintVal uint64
  uintVal = v.Uint()
  return uint32(uintVal), nil
}
``` 
重要的是注意 `reflect.Uint()` 总是返回 `uint64`。这是为什么 `uint64` 可以编码为所有其他整型类型的原因。

下面是剩下的无符号整数类型应该如何被解压的：
```golang
// Extract Uint64
func extractUint64(v reflect.Value) (uint64, error) {
  if reflect.Kind() != reflect.Uint64 {
    return uint64(0), errors.New("Invalid input")
  }
  var uintVal uint64
  uintVal = v.Uint()
  return uintVal, nil
}
// Extract Uint16
func extractUint16(v reflect.Value) (uint16, error) {
  if reflect.Kind() != reflect.Uint16 {
    return uint16(0), errors.New("Invalid input")
  }
  var uintVal uint64
  uintVal = v.Uint()
  return uint16(uintVal), nil
}
// Extract Uint8
func extractUint8(v reflect.Value) (uint8, error) {
  if reflect.Kind() != reflect.Uint8 {
    return uint8(0), errors.New("Invalid input")
  }
  var uintVal uint64
  uintVal = v.Uint()
  return uint8(uintVal), nil
}
// Extract Uint
func extractUint(v reflect.Value) (uint, error) {
  if reflect.Kind() != reflect.Uint {
    return uint(0), errors.New("Invalid input")
  }
  var uintVal uint64
  uintVal = v.Uint()
  return uint(uintVal), nil
}
```

## 解析浮点数
在 Golang 中有两种少点数类型
```golang
    Float32
    Float64
```
`Float32` 类型表示大小为 32bit 的浮点数
`Float64` 类型表示大小为 64bit 的浮点数

为了解析每个类型的浮点数，`reflect.Value` 对象表示需要被转换不适当类型的浮点数对象。

下面是 `Float32` 应该如何解析的：
```golang
// Extract Float32
func extractFloat32(v reflect.Value) (float32, error) {
  if reflect.Kind() != reflect.Float32 {
    return float32(0), errors.New("Invalid input")
  }
  var floatVal float64
  floatVal = v.Float()
  return float32(floatVal), nil
}
```
重要的是注意 `reflect.Float()` 总是返回 `float64`。这是为什么 `float64` 能编码所有其他浮点数类型的原因。

下面是64-bit浮点数应该如何被解析：
```golang
// Extract Float64
func extractFloat64(v reflect.Value) (float64, error) {
  if reflect.Kind() != reflect.Float64 {
    return float64(0), errors.New("Invalid input")
  }
  var floatVal float64
  floatVal = v.Float()
  return floatVal, nil
}
```

## 解析复数值
在 Golang 中有 2 种复数类型：
```golang
    Complex64
    Complex128
```
`Complex64` 类型代表64bit大小的复数。
`Complex128` 类型代表128bit大小的复数。
为了解析每个复数类型，`reflect.Value` 对象代表需要被转换为适当类型的复数。

下面是 `Complex64` 应该如何被解析：
```golang
// Extract Complex64
func extractComplex64(v reflect.Value) (complex64, error) {
  if reflect.Kind() != reflect.Complex64 {
    return complex64(0), errors.New("Invalid input")
  }
  var complexVal complex128
  complexVal = v.Complex()
  return complex64(complexVal), nil
}
```
重要的是注意 `reflect.Complex()` 总是返回 `complex128`。这是为什么 `complex128` 能编码为所有其他复数的原因。

这儿是 `128位复数值应该怎样解析的：
```golang
// Extract Complex128
func extractComplex128(v reflect.Value) (complex128, error) {
  if reflect.Kind() != reflect.Complex128 {
    return complex128(0), errors.New("Invalid input")
  }
  var complexVal complex128
  complexVal = v.Complex()
  return complexVal, nil
}
```
## 解析字符串值
字符串值由 reflect 包中的常量 `string` 表示
它们可以通过使用对象的 `String()` 方法从 `reflect.Value` 对象中解析。

下面是 `String` 应该如何被解析：
```golang
// Extract String
func extractString(v reflect.Value) (string, error) {
  if reflect.Kind() != reflect.String {
    return "", errors.New("Invalid input")
  }
  return v.String(), nil
}
```

## 解析指针值
在 Golang 中有 2 种指针：
```golang
    Uintptr
    UnsafePointer
```
`Uintptr` 和 `UnsafePointer` 仅仅是在进程内存中虚拟地址 `uint` 值。它可以代表本地的一个变量或一个函数。

`Uintptr` 和 `UnsafePointer` 间的区别是 `Uintptr` 是被 Go 运行时类型检查的，而 `UnsafePointer` 不是。`UnsafePointer` 可以用来转换任何 Go 类型到任何其他类型，因为它们的内存布局是兼容的。如果你希望探索更多，请在下面评论，我会根据你的评论写更多。

`Uintptr` 和 `UnsafePointer` 可以分别通过使用 `Addr()` 和 `UnsafeAddr()` 从 `reflect.Value` 对象中解析。这儿有个示例展示 `Uintptr` 应该如何被解析：
```golang
Uintptr should be extracted this way
// Extract Uintptr
func extractUintptr(v reflect.Value) (uintptr, error) {
  if reflect.Kind() != reflect.Uintptr {
    return uintptr(0), errors.New("Invalid input")
  }
  var ptrVal uintptr
  ptrVal = v.Addr()
  return ptrVal, nil
}
```
下面是 `UnsafePointer` 应该如何被解析
```golang
// Extract UnsafePointer
func extractUnsafePointer(v reflect.Value) (unsafe.Pointer, error) {
  if reflect.Kind() != reflect.UnsafePointer {
    return unsafe.Pointer(0), errors.New("Invalid input")
  }
  var unsafeVal unsafe.Pointer
  unsafeVal = unsafe.Pointer(v.UnsafeAddr())
  return unsafeVal, nil
}
```
值得注意的是上面的 `v.UnsafeAddr()` 返回一个 `uintptr` 的值。它应该在同一行中进行类型转换，否则unsafe.Pointer值可能不会指向预期的位置。

## 接下来是什么
请注意，在没有检查它们的 `kind` 的情况下不要使用任何 `reflect.Value` 结构体中的方法，因为它会很容易导致一个 `painc`。

我将在下一篇博客文章中撰写关于创建更复杂的类型，比如 `struct`,`pointer`,`chan`,`map`,`slice`,`array`等待，敬请关注！
