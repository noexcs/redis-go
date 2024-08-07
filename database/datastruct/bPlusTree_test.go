package datastruct

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

func TestBPlusTree(t *testing.T) {
	m := make(map[string]any)
	bPlusTree := MakeBPlusTree[string]()
	output, err := os.Create("test_operations.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = fmt.Fprintln(output, "bPlusTree := MakeNew()")
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < 10_0000; i++ {
		key := strconv.Itoa(rand.Intn(1000))
		if _, ok := m[key]; ok {
			// 验证在BPlusTree中也存在
			if _, _, existInBPlusTree := bPlusTree.Get(key); !existInBPlusTree {
				t.Errorf("错误: 应该存在该key")
			}

			// 已存在，则删除
			delete(m, key)
			fmt.Fprintf(output, "bPlusTree.Delete(%v)\n", key)
			bPlusTree.Delete(key)
			if !Diagnose(bPlusTree, output) {
				t.Errorf("diagnostic failed after Delete")
			}

			// 验证在 BPlusTree 中已经不存在
			if _, _, existInBPlusTree := bPlusTree.Get(key); existInBPlusTree {
				t.Errorf("错误: 删除后仍然存在")
			}
		} else {
			// 验证在 BPlusTree 中也不存在
			if _, _, existInBPlusTree := bPlusTree.Get(key); existInBPlusTree {
				t.Errorf("错误: 不应该存在该key")
			}

			// 不存在，则添加
			m[key] = key
			fmt.Fprintf(output, "bPlusTree.Insert(%v, %v)\n", key, key)
			bPlusTree.Insert(key, key, false)
			if !Diagnose(bPlusTree, output) {
				t.Errorf("diagnostic failed after Insert")
			}

			// 验证在 BPlusTree 中已经存在
			if _, _, existInBPlusTree := bPlusTree.Get(key); !existInBPlusTree {
				t.Errorf("错误: 添加后不在")
			}
		}
	}
	err = output.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	for k := range m {
		_, _, exist := bPlusTree.Get(k)
		if !exist {
			t.Errorf("错误: 添加后不在")
		}
	}
	iterator := bPlusTree.Iterator()
	for iterator.Next() {
		k, _ := iterator.Value()
		if _, exist := m[k]; !exist {
			t.Errorf("错误: 删除后仍在存在")
		}
	}
	fmt.Println("Test end")
}

// Diagnose 检查结构是否正确
func Diagnose[T BPTKey](t *BPlusTree[T], output io.Writer) bool {
	if t.root == nil {
		return true
	}
	leafOccurred := false
	if !t.root.isLeaf {
		queue := make([]*UnionNode[T], 0, 2*MaxChildren)
		queue = append(queue, t.root)
		var cursor *UnionNode[T]
		for len(queue) > 0 {
			cursor = queue[0]
			queue = queue[1:]
			if !cursor.isLeaf {
				if leafOccurred {
					// 因为是层序遍历，所以如果出现过叶子节点，就不可能再出现内部节点
					fmt.Fprintf(output, "错误：叶子节点不在同一深度")
				}
				if cursor.parent != nil && (len(cursor.children) < MinChildren || MaxChildren < len(cursor.children)) {
					fmt.Fprintf(output, "错误: InternalNode 大小不合法。")
					return false
				}
				for idx, childPtr := range cursor.children {
					if childPtr == nil {
						fmt.Fprintf(output, "错误: 不应有空指针的子节点")
						return false
					}
					if childPtr.parent != cursor {
						fmt.Fprintf(output, "错误: 父子节点关系错误")
						return false
					}
					queue = append(queue, childPtr)
					if !childPtr.isLeaf {
						if idx == 0 {
							if childPtr.leftPtr != nil {
								fmt.Fprintf(output, "错误: 不应有左兄弟")
								return false
							}
						}
						if idx == len(cursor.children)-1 {
							if childPtr.rightPtr != nil {
								fmt.Fprintf(output, "错误: 不应有右兄弟")
								return false
							}
						}
						if childPtr.leftPtr != nil && childPtr.leftPtr.rightPtr != childPtr {
							fmt.Fprintf(output, "错误: 左右兄弟不一致1")
							return false
						}
						if childPtr.rightPtr != nil && childPtr.rightPtr.leftPtr != childPtr {
							fmt.Fprintf(output, "错误: 左右兄弟不一致2")
							return false
						}
					} else if childPtr.isLeaf {
						if childPtr.leftPtr != nil && childPtr.leftPtr.rightPtr != childPtr {
							fmt.Fprintf(output, "错误: 左右兄弟不一致3")
							return false
						}
						if childPtr.rightPtr != nil && childPtr.rightPtr.leftPtr != childPtr {
							fmt.Fprintf(output, "错误: 左右兄弟不一致4")
							return false
						}
					}

					maxKey, minKey := getMaxMinKey(childPtr)
					if idx-1 >= 0 && CompareBPTKeys[T](cursor.keys[idx-1], minKey) == 1 {
						fmt.Fprintf(output, "错误: 子节点的最小值大于了父节点对应的key")
						return false
					}
					if idx < len(cursor.keys) && CompareBPTKeys[T](cursor.keys[idx], maxKey) == -1 {
						fmt.Fprintf(output, "错误: 子节点的最大值大于了父节点对应的key的后面一个key")
						return false
					}
				}
			} else if cursor.isLeaf {
				leafOccurred = true
				if cursor.parent != nil && len(cursor.kvPairs) < MinKeys || MaxKeys < len(cursor.kvPairs) {
					fmt.Fprintf(output, "错误: LeafNode 大小不合法。")
					return false
				}
				for idx := range cursor.kvPairs {
					if idx-1 > 0 && CompareBPTKeys[T](cursor.kvPairs[idx].key, cursor.kvPairs[idx-1].key) == -1 {
						fmt.Fprintf(output, "错误: KVPair内部顺寻不一致。")
						return false
					}
				}
			}
		}

	} else {
		// root为叶节点的话，最多只能有一个leafNode
		if t.root.parent != nil {
			fmt.Fprintf(output, "错误: 不应有parent。")
			return false
		}
		if t.root.leftPtr != nil {
			fmt.Fprintf(output, "错误: 不应有leftPtr。")
			return false
		}
		if t.root.rightPtr != nil {
			fmt.Fprintf(output, "错误: 不应有rightPtr。")
			return false
		}
		for idx := range t.root.kvPairs {
			if idx > 0 && CompareBPTKeys[T](t.root.kvPairs[idx].key, t.root.kvPairs[idx-1].key) == -1 {
				fmt.Fprintf(output, "错误: KVPair内部顺寻不一致。")
				return false
			}
		}
	}
	return true
}

func getMaxMinKey[T BPTKey](node *UnionNode[T]) (max T, min T) {
	if node.isLeaf {
		if len(node.kvPairs) == 0 {
			panic("Empty leafNode.")
		}
		max, min = node.kvPairs[0].key, node.kvPairs[0].key
		for i := 1; i < len(node.kvPairs); i++ {
			if CompareBPTKeys[T](node.kvPairs[i].key, min) == -1 {
				min = node.kvPairs[i].key
			}
			if CompareBPTKeys[T](node.kvPairs[i].key, min) == 1 {
				max = node.kvPairs[i].key
			}
		}
		return max, min
	} else {
		if len(node.children) < 2 {
			panic("Illegal InternalNode.")
		}
		max, min = getMaxMinKey[T](node.children[0])
		for i := 1; i < len(node.children); i++ {
			maxKey, minKey := getMaxMinKey(node.children[i])
			if CompareBPTKeys[T](maxKey, max) == 1 {
				max = maxKey
			}
			if CompareBPTKeys[T](maxKey, max) == -1 {
				min = minKey
			}
		}
		return max, min
	}
}
