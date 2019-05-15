### example

example 采用的是 hashjoin 的方式来进行 join的。
对于 hashjoin 的 join 方式来讲，hashjoin 只支持等值连接，题目只有等值连接，符合题意但适用范围较小。

example 在执行 join 的过程中，没有判断两个文件的大小，直接对 table1 进行建立 hashtable 。最好使用更小的表进行 hashtable 来构建，用大一点的表进行探测。

### 感觉 join 算法都太难实现了，给大家表演一个循环嵌套查询
