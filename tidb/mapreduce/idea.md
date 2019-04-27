### 实现方案

- 考虑通过 map/reduce 实现 url 计数 （类单词计数）
- 对 map/reduce 的结果进行归并排序
- 取归并排序结果的前 10 个

#### 4 月 24 日

基本了解框架是怎么调度的，考虑优先完成框架，并跑通 `urltop10_example`

- 结果 
	
当天完成（ `utils` 里的工具函数真好用）

#### 4 月 25 日

完成了自己的 MapF 和 ReduceF 函数。

**相关逻辑**

- MapReduce 过程。

`main` 包中的 `init` 函数，开始一个 `MRCluster` 的启动。提交待处理的任务，待处理的任务经过`run`函数进行分解提交给 `worker` 线程处理。`worker` 线程对对应的文件进行处理，得到 `MapF` 和 `ReduceF` 可以处理的数据类型。其中 `Map` 函数处理 string 类型的数据，`Reduce` 处理 json 类型的数据。

- `Example` 相关

两次 `MapReduce` 过程

第一次 Map 过程，处理对应的文件，单纯的把 string 数据转换为 json 数据

第一次 Reduce 过程，把相同的 url 文件转换为 `url count` 的形式

第二次 Map 过程，把所有不同的 url 处理为 `key:"",value: "url count"`形式，保证了key相同，以便于 `Reduce` 函数处理

第二次 Reduce 过程，取所有不同的 url 中出现次数最多的前 10 个，并转换为 string 形式

一些弊端： 

`MapF` 和 `ReduceF`的结果都会被写入磁盘，第一次 `Map` 过程写入的文件，相同的 key 存储了多次，存储了不必要的数据，有对应的数据开销。

第二次 `Map` 的过程，实际上已经已经可以只保留对应数据中的 top10，但是保留了所有文件，也是有存储的开销和对应磁盘写入的开销。


- 我的代码

也是两次 `MapReduce` 过程

第一次 Map 过程，处理对应的文件，对单个文件做 url 计数，输出 json 格式的 `key:"",value: "url count"`

第一次 Reduce 过程，把其中中相同的 url 对应的 count 求和，文件转换为 string 格式的 `url count` 形式

第二次 Map 过程，取输入文件中 url 出现次数最多的前十个，并转换为 `key:"",value: "url count"`形式，保证key相同，以便于 `Reduce` 函数处理

第二次 Reduce 过程，取多个文件结果中的前 10 个，并转换为 string 形式

- 相比 `Example` 的优势

弥补了重复存储 url 的缺点，对于取出现次数最多的 url 的处理，先在 Map 过程中做了预处理，减小了 Reduce 过程中的计算量


- profile 分析

**Example**

- 通过 top10 来看，耗时比较高的几个调用中 json 相关的调用比较多，主要耗时为 json 的编码与解码

- `list worker` 以后可以看到 worker 中的主要耗时为将 Map 过程的结果编码为 json 格式写入中间文件以及 Reduce 过程从中间文件把 json 格式解析为 keyvalue 的格式

从以上两点来看，和代码层面的分析一致，已经在 mapreduce 的过程中尽量避免保留多余的数据，来避免 json 编码与解码的耗时。

**我的代码**

- 从图上来看，耗时还是 json 的相关操作耗时是性能的瓶颈（要不要考虑换个库来优化一下）