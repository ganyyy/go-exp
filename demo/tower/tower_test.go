package tower

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type Echo struct {
}

func (e Echo) Call(ctx context.Context, req string) (string, error) {
	fmt.Println("Echo service called with request:", req)
	return req, nil
}

func TestTower(t *testing.T) {
	var echoService Echo

	srv := NewServiceBuilder[string, string]().
		WithLayer(NewRecovery[string, string]()).
		WithLayer(NewLogger[string, string]("before")).
		WithLayer(NewTimeout[string, string](1000)).
		WithLayer(NewLogger[string, string]("after")).
		Build(echoService)

	resp, err := srv.Call(context.Background(), "Hello, World!")
	require.NoError(t, err)
	require.Equal(t, "Hello, World!", resp)
}
