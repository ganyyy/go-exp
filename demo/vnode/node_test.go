package vnode

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

func TestVirtualNode(t *testing.T) {
	var shard Sharding
	shard.VirtualNum = 35
	shard.RealNum = 4
	shard.CurrentNodes = shard.RealNum
	shard.NodesPrefixKey = ShardingNodes + Split
	for i := 0; i < int(shard.RealNum); i++ {
		shard.Nodes[i] = "node" + strconv.Itoa(i)
	}

	require.NoError(t, shard.Balance(), "should balance without error")

	t.Log(shard.String())

	require.NoError(t, shard.IsValid())

	var newShard Sharding
	newShard.VirtualNum = shard.VirtualNum
	newShard.RealNum = shard.RealNum
	newShard.NodesPrefixKey = shard.NodesPrefixKey
	for _, nodeId := range shard.Nodes[:shard.CurrentNodes] {
		virtualNodes, ok := shard.GetVirtualNodes(nodeId)
		require.True(t, ok, "should get virtual nodes for node %s", nodeId)
		var kv = mvccpb.KeyValue{
			Key:   []byte(newShard.NodesPrefixKey + nodeId),
			Value: []byte(Numbers2String(virtualNodes, ShardingNodesSplit)),
		}
		require.NoError(t, newShard.Update(&kv, mvccpb.PUT), "should add node %s without error", nodeId)
	}

	t.Logf("new shard: %s", newShard.String())

	require.Equal(t, shard, newShard)
	require.NoError(t, newShard.IsValid(), "new shard should be valid")
}
