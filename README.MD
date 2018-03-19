<div align=center>
<img src="https://github.com/johng-cn/drh-go/blob/master/drh.png" height="600"/>
</div>


DRH(Deep-Re-Hash)深度哈希分区算法是一种针对哈希表在海量数据及磁盘存储下的一种改进算法，它的时间复杂度介于常数O(1)和对数O(d*log (n-1))之间(即：O(1) <= T(n) <= O(d*log(n-1)) ，其中n为阶数，d为深度)，提供了极高的数据检索、插入、修改、删除效率。