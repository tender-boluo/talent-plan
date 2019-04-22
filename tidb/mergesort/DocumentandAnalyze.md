## 心路历程

### 单线程归并排序
首先考虑完成一个单线程归并排序

遇到的问题如下

- `src slice` 并没有与如期一样更改为排序后的数组

解决方式如下

- 增加了一个 `MergeSort` 函数并获取他的返回值，第一次用 `go` ，可能写法上不够优美，但得到了预期的结果

`benchmark` 的结果如下

- 时间为 `sort.slice` 的 2 倍 

### goroutine 并发归并排序
尝试使用多协程并发来优化单线程归并排序的时间

#### 第一次考虑使用 `channel` 直接并发单线程归并排序的代码

`benchmark` 的结果

- 时间为 `sort.slice` 的 4 倍

考虑的解决方式如下

- 猜测原因是协程之间的调度导致了时间的增加——减少并发的 `goroutine` 数量

结果

- `goroutine` 的并发和函数的递归调用不能很好的结合起来，由于函数的递归调用所以不能减少 `goroutine` 的数量

#### 第二次考虑使用协程池来减少协程的数量

结果

- 递归调用时，不达到递归的叶子节点时，所有的 `goroutine` 都会处于阻塞状态，所以实际上递归时不能减少 `goroutine` 的数量

#### 第三次考虑不使用递归来实现归并排序
即归并时，传入当前待归并的两个数组的起始位置与结束位置来完成归并

`benchmark` 结果

- `benchmark` 失败了... 打印log后发现内存不足

原因如下

- `merge` 的时候，开启的辅助数组由于没有初始化大小，导致 `apend` 执行时 `slice` 不停的申请新的内存，而旧的内存占用在等待 `gc` 又没有还给操作系统导致内存不足

解决办法

- 初始化了辅助数组的大小解决了内存占用的问题

`benchmark` 结果

- 时间为 `sort.slice` 的 4 倍

根据 `cpuprofile` 的结果进行排查

排查结果如下

- 4 倍的时间开销中有 `1/2` 用于了内存分配

解决办法

- 考虑到内存分配主要用于辅助数组的声明，把辅助数组开于 `MergeSort` 的函数中，并传给 `Merge` 的函数使用，以此减少内存的申请

`benchmark` 结果

- 时间略微多于 `sort.slice` 

根据 `cpuprofile` 的结果进行排查

排查结果如下

- `web` 图片中显示开销主要为底层调用
- `list` 显示开销主要为 `slice` 之间的赋值... 赋值？

考虑原因

- 考虑主要原因为 `go` 中没有引用传递，只有值传递，在 `slice` 传递的过程中有大量的浅拷贝，导致时间开销比较大

解决办法

- 对于初始数组先分块进行插入排序，各部分相对有序，大量减少拷贝导致的开销

`benchmark` 结果

- 时间为 `sort.slice` 的 `1/3`

优缺点

- 对于数组内元素比较多的情况，优化比较明显
- 对于数组内元素比较少时，时间复杂度反而增高