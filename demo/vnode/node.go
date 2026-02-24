package vnode

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"golang.org/x/exp/constraints"
)

/*

etcd 存储路径

sharding/
		meta/
			lock, 		 int  // 锁节点，防止并发修改
			real_num, 	 int  // 真实节点期望的数量
			virtual_num, int  // 虚拟节点总数
		nodes/
			{node_id,   [int] // 虚拟节点列表
services/
		pre-alloc/
			[service_id]      // 预分配服务，服务 ID
		running/
			[service_id]      // 服务运行中，服务 ID

*/

type Sharding struct {
	VirtualNum uint // 虚拟节点总数
	RealNum    uint // 真实节点期望的数量

	CurrentNodes uint                // 当前在线的物理节点数量
	Nodes        [MaxRealNum]string  // 分片节点列表
	VirtualNodes [MaxVirtualNum]int8 // 虚拟节点列表，索引为虚拟节点 ID，值为物理节点索引（从 0 开始），NotAssigned 表示未分配

	MetaRealNumKey    string // 元数据真实节点数量的键
	MetaVirtualNumKey string // 元数据虚拟节点数量的键
	NodesPrefixKey    string // 节点前缀键
}

const (
	MaxRealNum    = 64   // 真实节点期望的数量，不能超过 64
	MaxVirtualNum = 1024 // 虚拟节点总数，不能超过 1024

	NotAssigned = -1 // 未分配的虚拟节点标记
)

const (
	Split                  = "/"
	ShardingRoot           = Split + "sharding"
	ShardingMeta           = ShardingRoot + Split + "meta"
	ShardingLock           = ShardingMeta + Split + "lock"
	ShardingMetaRealNum    = ShardingMeta + Split + "real_num"
	ShardingMetaVirtualNum = ShardingMeta + Split + "virtual_num"

	ShardingNodes      = ShardingRoot + Split + "nodes"
	ShardingNodesSplit = ","
)

func NewSharding(client *clientv3.Client, root string) (*Sharding, error) {

	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	resp, err := client.Get(ctx,
		root+ShardingRoot,
		clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend),
	)

	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, errors.New("sharding not found")
	}

	metaRealNumKey := root + ShardingMetaRealNum
	metaVirtualNumKey := root + ShardingMetaVirtualNum
	nodesPrefixKey := root + ShardingNodes + Split

	var sharding = &Sharding{
		VirtualNum: 0,
		RealNum:    0,

		MetaRealNumKey:    metaRealNumKey,
		MetaVirtualNumKey: metaVirtualNumKey,
		NodesPrefixKey:    nodesPrefixKey,
	}

	for _, kv := range resp.Kvs {
		if err := sharding.Update(kv, mvccpb.PUT); err != nil {
			return nil, fmt.Errorf("update sharding from key %q (%q) failed: %w", kv.Key, kv.Value, err)
		}
	}

	return sharding, nil
}

func bindNumber(val []byte, target *uint) error {
	if len(val) == 0 {
		*target = 0
		return nil
	}
	var err error
	var ret uint64
	ret, err = strconv.ParseUint(string(val), 10, 32)
	if err != nil {
		return err
	}
	*target = uint(ret)
	return nil
}

func (shard *Sharding) ClearRealVirtualNode(idx int) {
	if idx >= len(shard.Nodes) {
		return
	}
	shard.Nodes[idx] = ""
	for virtualNode := range shard.VirtualNodes[:shard.VirtualNum] {
		if shard.VirtualNodes[virtualNode] == int8(idx) {
			shard.VirtualNodes[virtualNode] = NotAssigned // 标记为未分配
		}
	}
}

func (shard *Sharding) GetNextEmptyNodeIndex() int {
	for idx, node := range shard.Nodes {
		if node == "" {
			return idx
		}
	}
	return NotAssigned
}

func (shard *Sharding) GetNodeByID(nodeID string) int {
	if len(nodeID) == 0 {
		return NotAssigned
	}
	for idx, node := range shard.Nodes {
		if node == nodeID {
			return idx
		}
	}
	return NotAssigned
}

func (shard *Sharding) GetVirtualNodes(nodeId string) ([]uint16, bool) {
	if len(nodeId) == 0 {
		return nil, false
	}
	var nodeIdx = shard.GetNodeByID(nodeId)
	if nodeIdx == NotAssigned {
		return nil, false
	}
	var virtualNodes []uint16
	for virtualNode, nodeIndex := range shard.VirtualNodes[:shard.VirtualNum] {
		if int(nodeIndex) == nodeIdx {
			virtualNodes = append(virtualNodes, uint16(virtualNode))
		}
	}
	return virtualNodes, true
}

func (shard *Sharding) Update(kv *mvccpb.KeyValue, typ mvccpb.Event_EventType) error {
	if kv == nil {
		return nil
	}
	key, val := string(kv.Key), kv.Value
	if key == shard.MetaRealNumKey {
		var num uint
		if err := bindNumber(val, &num); err != nil {
			return fmt.Errorf("parse real_num failed: %w", err)
		}
	} else if key == shard.MetaVirtualNumKey {
		if err := bindNumber(val, &shard.VirtualNum); err != nil {
			return fmt.Errorf("parse virtual_num failed: %w", err)
		}
	} else if nodeID, ok := strings.CutPrefix(key, shard.NodesPrefixKey); ok {
		if len(nodeID) == 0 {
			return fmt.Errorf("invalid node ID %q, should not be empty", nodeID)
		}
		if strings.Contains(nodeID, Split) {
			return fmt.Errorf("invalid node ID %q, should not contain %q", nodeID, Split)
		}

		isDelete := typ == mvccpb.DELETE

		var nodeIdx = shard.GetNodeByID(nodeID)
		if isDelete {
			if nodeIdx == NotAssigned {
				return nil
			}
			shard.ClearRealVirtualNode(nodeIdx)
			return nil
		}

		var isNew bool
		if nodeIdx == NotAssigned {
			// 找到一个为空的节点位置
			nodeIdx = shard.GetNextEmptyNodeIndex()
			if nodeIdx == NotAssigned {
				return fmt.Errorf("shard is full, cannot add new node %q", nodeID)
			}
			isNew = true
		}

		var tmpVirtualNodes [MaxVirtualNum]int8
		for i := range tmpVirtualNodes {
			tmpVirtualNodes[i] = NotAssigned // 初始化为未分配
		}
		var maxIdx int
		for part := range bytes.SplitSeq(val, []byte(ShardingNodesSplit)) {
			var virtualNode uint
			if err := bindNumber(part, &virtualNode); err != nil {
				return fmt.Errorf("parse node %q(%q) virtual node %q failed: %w", nodeID, val, part, err)
			}
			if virtualNode >= min(shard.VirtualNum, MaxVirtualNum) {
				return fmt.Errorf("node %q(%q) virtual node %d is out of range, should be less than %d(%d)", nodeID, val, virtualNode, MaxVirtualNum, shard.VirtualNum)
			}
			if tmpVirtualNodes[virtualNode] != NotAssigned {
				return fmt.Errorf("node %q(%q) virtual node %d is duplicated", nodeID, val, virtualNode)
			}
			tmpVirtualNodes[virtualNode] = int8(nodeIdx) // 分配虚拟节点到物理节点
			maxIdx = max(maxIdx, int(virtualNode))       // 记录最大虚拟节点索引
		}

		if isNew {
			shard.Nodes[nodeIdx] = nodeID
			shard.CurrentNodes++
		} else {
			shard.ClearRealVirtualNode(nodeIdx) // 清除旧的虚拟节点分配
		}
		for virtualNode, nodeIndex := range tmpVirtualNodes[:maxIdx+1] {
			if nodeIndex == NotAssigned {
				continue // 跳过未分配的虚拟节点
			}
			shard.VirtualNodes[virtualNode] = nodeIndex // 分配虚拟节点到物理节点
		}
	}

	concurrency.NewSTM(nil, func(stm concurrency.STM) error {
		stm.Get()
		return nil
	}, concurrency.WithPrefetch())

	return nil
}

func (shard *Sharding) baseCheck() error {
	if shard.RealNum == 0 || shard.RealNum > MaxRealNum {
		return errors.New("real_num is invalid")
	}
	if shard.VirtualNum == 0 || shard.VirtualNum > MaxVirtualNum {
		return errors.New("virtual_num is invalid")
	}
	if shard.VirtualNum < shard.RealNum {
		return fmt.Errorf("virtual_num %d is less than real_num %d", shard.VirtualNum, shard.RealNum)
	}
	if shard.RealNum < shard.CurrentNodes {
		return fmt.Errorf("real_num %d is less than the number of nodes %d", shard.RealNum, shard.CurrentNodes)
	} else if shard.RealNum > shard.CurrentNodes {
		return ErrNeedWaitRealNodeOnline // 上层需要等待 real_num 个物理节点上线
	}

	// 注意: shard.Nodes[:shard.RealNum] 不能包含空字符串, 因为它们是物理节点的标识符
	// nodeId 唯一性检查
	var totalNodes = 0
	for i, node := range shard.Nodes[:shard.RealNum] {
		if node == "" {
			return fmt.Errorf("node at index %d is empty", i)
		}
		for j, otherNode := range shard.Nodes[i+1 : shard.RealNum] {
			if node == otherNode {
				return fmt.Errorf("node %q at index %d is duplicated with node %q at index %d", node, i, otherNode, i+j+1)
			}
		}
		totalNodes++
	}
	if totalNodes != int(shard.CurrentNodes) {
		return fmt.Errorf("total nodes %d does not match current nodes %d", totalNodes, shard.CurrentNodes)
	}
	return nil
}

func (shard *Sharding) IsValid() error {

	if err := shard.baseCheck(); err != nil {
		return err
	}

	// 虚拟节点分配检查
	for virtualNode, nodeIndex := range shard.VirtualNodes[:shard.VirtualNum] {
		if nodeIndex == NotAssigned {
			return fmt.Errorf("virtual node %d is not assigned to any real node", virtualNode)
		}
		if int(nodeIndex) >= int(shard.CurrentNodes) || shard.Nodes[nodeIndex] == "" {
			return fmt.Errorf("virtual node %d is assigned to an invalid real node index %d", virtualNode, nodeIndex)
		}
	}
	// 平衡性检查
	var virtualNodePerNode [MaxRealNum]int // 每个物理节点分配的虚拟节点数量
	for _, nodeIndex := range shard.VirtualNodes[:shard.VirtualNum] {
		virtualNodePerNode[nodeIndex]++ // 前边已经校验过了, 这里直接用
	}
	// 最大值和最小值之间的差距应该不超过 1
	var maxVirtualNodes, minVirtualNodes = 0, int(shard.VirtualNum)
	for _, count := range virtualNodePerNode[:shard.RealNum] {
		if count > maxVirtualNodes {
			maxVirtualNodes = count
		}
		if count < minVirtualNodes {
			minVirtualNodes = count
		}
	}
	if maxVirtualNodes-minVirtualNodes > 1 {
		return fmt.Errorf("%w: evenly distributed, max: %d, min: %d", ErrVirtualNodeNeedBalance, maxVirtualNodes, minVirtualNodes)
	}
	return nil
}

// String returns a string representation of the Sharding.
func (shard *Sharding) String() string {
	var sb strings.Builder

	if err := shard.baseCheck(); err != nil {
		return "Invalid Sharding: " + err.Error()
	}

	sb.WriteString("Sharding:\n")
	sb.WriteString("  VirtualNum: " + strconv.FormatUint(uint64(shard.VirtualNum), 10) + "\n")
	sb.WriteString("  RealNum: " + strconv.FormatUint(uint64(shard.RealNum), 10) + "\n")
	sb.WriteString("  CurrentNodes: " + strconv.FormatUint(uint64(shard.CurrentNodes), 10) + "\n")
	var tmpVirtualNodesCache [MaxRealNum]bytes.Buffer
	for vid, nodeId := range shard.VirtualNodes[:shard.VirtualNum] {
		if nodeId == NotAssigned {
			continue
		}
		tmpVirtualNodesCache[nodeId].WriteString(strconv.FormatUint(uint64(vid), 10))
		tmpVirtualNodesCache[nodeId].WriteByte(' ')
	}
	for i, node := range shard.Nodes[:shard.RealNum] {
		sb.WriteString(fmt.Sprintf("  %s ( %s)\n", node, tmpVirtualNodesCache[i].String()))
	}

	return sb.String()
}

// Balance 虚拟节点平衡, 应该只有在 ErrVirtualNodeNeedBalance 错误时, 并且成功拿到锁时才会调用
// 这个函数会尝试平衡虚拟节点到物理节点上,
func (shard *Sharding) Balance() error {
	if err := shard.baseCheck(); err != nil {
		return err
	}

	var tmpVirtualNodes [MaxVirtualNum]int8
	for i := range tmpVirtualNodes[:shard.VirtualNum] {
		tmpVirtualNodes[i] = int8(i)
	}
	rand.Shuffle(int(shard.VirtualNum), func(i, j int) {
		tmpVirtualNodes[i], tmpVirtualNodes[j] = tmpVirtualNodes[j], tmpVirtualNodes[i]
	})
	clear(shard.VirtualNodes[:])                              // 清除原有的虚拟节点分配
	var virtualNodePerNode = shard.VirtualNum / shard.RealNum // 每个物理节点分配的虚拟节点数量
	var remainder = shard.VirtualNum % shard.RealNum          // 余数, 需要额外分配给前 remainder 个物理节点

	var startVirtualIdx uint = 0 // 虚拟节点的开始索引
	for idx := range shard.CurrentNodes {
		endVirtualIdx := startVirtualIdx + virtualNodePerNode
		if idx < remainder {
			endVirtualIdx++ // 前 remainder 个物理节点多分配一个虚拟节点
		}
		for virtualNode := startVirtualIdx; virtualNode < endVirtualIdx; virtualNode++ {
			shard.VirtualNodes[tmpVirtualNodes[virtualNode]] = int8(idx)
		}
		startVirtualIdx = endVirtualIdx
	}

	// 写入分片信息到 etcd

	return nil
}

func Numbers2String[N constraints.Integer](numbers []N, split string) string {
	if len(numbers) == 0 {
		return ""
	}
	var sb strings.Builder
	for i, num := range numbers {
		if i > 0 {
			sb.WriteString(split)
		}
		sb.WriteString(strconv.FormatInt(int64(num), 10))
	}
	return sb.String()
}

func (shard *Sharding) Clone() *Sharding {
	if shard == nil {
		return nil
	}
	sv := *shard
	return &sv
}

var (
	ErrNeedWaitRealNodeOnline = errors.New("need to wait for real nodes to be online")
	ErrVirtualNodeNeedBalance = errors.New("virtual nodes need to be balanced")
)

type ShardingManager struct {
	Sharding atomic.Pointer[Sharding] // 当前分片信息
}
