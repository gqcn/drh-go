// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// DRH算法实现的Map.
package gdrh

// DRH算法操作对象
type Map struct {
    size   int          // 分区基数
    degree int          // 分区阶数
    root   *drhTable    // 根哈希表
}

// 哈希表
type drhTable struct {
    m      *Map         // 所属Map对象
    p      *drhPart     // 关联的分区对象(哪个分区指向该哈希表)
    deep   int          // 表深度(其实没什么意义，打印层级时候可以用得上)
    size   int          // 表分区(除根节点外，必须为奇数)
    parts  []*drhPart   // 分区数组
}

// 哈希表中的分区
type drhPart struct {
    t      *drhTable    // 所属哈希表
    items  []*drhItem   // 数据项列表，必须按照值进行从小到大排序，便于二分查找
    table  *drhTable    // 深度分区标识，指向另外一个哈希表
}

// 分区中的数据项
type drhItem struct {
    key   int           // 数据项键名，这里设置为int，便于演示
    value interface{}   // 数据项键值
}

// 创建DRH对象
func New(size, degree int) *Map {
    m := &Map {
        size   : size,
        degree : degree,
    }
    m.root = &drhTable {
        m     : m,
        deep  : 0,
        size  : size,
        parts : make([]*drhPart, size),
    }
    return m
}

// 设置键值对数据
func (m *Map) Set(key int, value interface{}) {
    m.root.set(key, value)
}

// 根据键名查询键值
func (m *Map) Get(key int) interface{} {
    p := m.root.search(key)
    if p != nil {
        result, cmp := p.search(key)
        if cmp == 0 {
            return p.items[result].value
        }
    }
    return nil
}

// 根据键名删除键值对
func (m *Map) Remove(key int) {
    p := m.root.search(key)
    if p != nil {
        result, cmp := p.search(key)
        if cmp == 0 {
            p.remove(result)
        }
    }
}

// 在当前哈希表上计算分区索引号
func (t *drhTable) getPartIndexByKey(key int) int {
    return key % t.size
}

// 根据键名检索对应的分区块，返回对应分区块对象指针，找不到则返回nil
func (t *drhTable) search(key int) *drhPart {
    p := t.parts[t.getPartIndexByKey(key)]
    if p != nil {
        if p.table != nil {
            return p.table.search(key)
        } else {
            return p
        }
    }
    return nil
}

// 在当前哈希表上设置键值对数据
func (t *drhTable) set(key int, value interface{}) {
    p := t.search(key)
    if p == nil {
        t.parts[t.getPartIndexByKey(key)] = &drhPart {
            t     : t,
            items : []*drhItem{{
                key   : key,
                value : value,
            }},
        }
    } else {
        if p.table != nil {
            p.table.set(key, value)
            return
        }
        index, cmp := p.search(key)
        if cmp == 0 {
            p.items[index].value = value
        } else {
            // 首先进行进行数据项插入
            p.save(&drhItem{key : key, value : value }, index, cmp)
            // 接着再判断是否需要进行DRH算法处理
            p.checkAndDoDeepReHash()
        }
    }
}

// 在当前哈希表上设置键值对数据
func (t *drhTable) remove(key int) {
    p := t.search(key)
    if p != nil {
        if p.table != nil {
            p.table.remove(key)
            return
        }
        index, cmp := p.search(key)
        if cmp == 0 {
            p.remove(index)
            // 如果分区元素
            if len(p.items) == 0 {
                p.table = nil
            }
        }
    }
}

// 对当前分区执行DRH算法，重新散列该分区数据到新的哈希表
func (p *drhPart) checkAndDoDeepReHash() {
    // 再判断是否需要进行DRH算法处理
    if len(p.items) != p.t.m.degree {
        return
    }
    // 计算分区增量，保证数据散列(分区后在同一请求处理中不再进行二次分区)
    // 分区增量必须为奇数，保证分区数据分配均匀
    size := p.t.size + 1
    if size%2 == 0 {
        size++
    }
    parts := make(map[int][]*drhItem)
    done  := true
    for {
        for i := 0; i < len(p.items); i ++ {
            index := p.items[i].key%size
            if _, ok := parts[index]; !ok {
                parts[index] = make([]*drhItem, 0)
            }
            parts[index] = append(parts[index], p.items[i])
            if len(parts[index]) == p.t.m.degree {
                done  = false
                parts = make(map[int][]*drhItem)
                size += 2 // 奇数+2必定为奇数
                break
            }
        }
        if done {
            break
        } else {
            done = true
        }
    }

    // 分区必定会成功，这里递增哈希表，增加深度
    table := &drhTable {
        m     : p.t.m,
        p     : p,
        deep  : p.t.deep + 1,
        size  : size,
        parts : make([]*drhPart, size),
    }
    for k, v := range parts {
        table.parts[k] = &drhPart{
            t     : table,
            items : v,
        }
    }
    p.items = nil
    p.table = table
}

// 添加一项, cmp < 0往前插入，cmp >= 0往后插入
func (p *drhPart) save(item *drhItem, index int, cmp int) {
    if cmp == 0 {
        p.items[index] = item
        return
    }
    pos := index
    if cmp == -1 {
        // 添加到前面
    } else {
        // 添加到后面
        pos = index + 1
        if pos >= len(p.items) {
            pos = len(p.items)
        }
    }
    rear   := append([]*drhItem{}, p.items[pos : ]...)
    p.items = append(p.items[0 : pos], item)
    p.items = append(p.items, rear...)
}

// 删除数组元素
func (p *drhPart) remove(index int) {
    p.items = append(p.items[ : index], p.items[index + 1 : ]...)
}

// 在当前分区上进行二分检索
// 返回值1: 二分查找中最近对比的数组索引
// 返回值2: -2表示压根什么都未找到，-1表示最近一个索引对应的值比key小，1最近一个索引对应的值比key大
// 两个值的好处是即使匹配不到key,也能进一步确定插入的位置索引
func (p *drhPart) search(key int) (int, int) {
    min := 0
    max := len(p.items) - 1
    mid := 0
    cmp := -2
    for min <= max {
        mid = int((min + max) / 2)
        if key < p.items[mid].key {
            max = mid - 1
            cmp = -1
        } else if key > p.items[mid].key {
            min = mid + 1
            cmp = 1
        } else {
            return mid, 0
        }
    }
    return mid, cmp
}