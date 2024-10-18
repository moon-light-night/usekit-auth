package suite

import (
	"context"
	authv1 "github.com/moon-light-night/usekit-proto/gen/go/auth.v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
	"testing"
	"usekit-auth/internal/config"
)

const (
	grpcHost = "localhost"
)

type Suite struct {
	*testing.T                   // instance of object for executing testing functions inside test suite
	Cfg        *config.Config    // app config
	AuthClient authv1.AuthClient // client for interaction with grpc server
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("../config/config.yaml")
	// create context for child tests
	ctx, cancelContext := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	// call Cleanup when tests finish
	t.Cleanup(func() {
		t.Helper()
		// cancel context when tests finish
		cancelContext()
	})

	cc, err := grpc.DialContext(context.Background(),
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials())) // use insecure connection for tests
	if err != nil {
		t.Fatalf("grpc server connection failed: %v:", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: authv1.NewAuthClient(cc),
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}
