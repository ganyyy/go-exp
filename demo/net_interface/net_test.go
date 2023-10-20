package netinterface

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNetData(t *testing.T) {
	var srcData = &IntData{Data: 123}
	bs, err := MarshalDataToNet(srcData)
	require.NoError(t, err)
	t.Logf("marshal int data: \n%s", string(bs))

	data, err := UnmarshalDataFromNet[*IntData](bs)
	require.NoError(t, err)
	require.Equal(t, srcData, data)
}
