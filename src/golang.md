# 1. `const` 中的 `iota` 问题
```golang
const (
  CHECK_ALL = -1
  UNCHECK   = iota
  CHECKED
  REJECT
```
在上述代码中，本期望获得的值分别是：`-1,0,1,2` 实际的输出却是 `-1` `1` `2` `3`

> [https://golang.org/ref/spec#Iota](https://golang.org/ref/spec#Iota)
> 在常量声明中，预先声明的标识符 iota表示连续的无类型整数 常量。它的值是该 常量声明中相应ConstSpec的索引，从零开始。
> 所以在上面的代码中，`UNCHECK` 的 `iota` 位于 `const` 中第二行，所以它的索引是1

上述代码可以改为:
```golang
  CHECK_ALL = iota - 1
  UNCHECK
  CHECKED
  REJECT
```
