package promgrpc_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/alexeyxo/promgrpc/v4"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
)

func TestNewClientMessageSentSizeStatsHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h := promgrpc.NewStatsHandler(promgrpc.NewClientMessageSentSizeStatsHandler(promgrpc.NewClientMessageSentSizeHistogramVec()))
	ctx = h.TagRPC(ctx, &stats.RPCTagInfo{
		FullMethodName: "/service/Method",
		FailFast:       true,
	})
	h.HandleRPC(ctx, &stats.OutHeader{
		Client: true,
		Header: metadata.MD{"user-agent": []string{"fake-user-agent"}},
	})
	h.HandleRPC(ctx, &stats.OutPayload{
		Client: true,
		Length: 5,
	})
	h.HandleRPC(ctx, &stats.OutPayload{
		Client: true,
		Length: 5,
	})
	h.HandleRPC(ctx, &stats.OutPayload{
		Client: true,
		Length: 5,
	})
	h.HandleRPC(ctx, &stats.OutPayload{
		Client: false,
		Length: 5,
	})

	const metadata = `
		# HELP grpc_client_message_sent_size_histogram_bytes TODO
        # TYPE grpc_client_message_sent_size_histogram_bytes histogram
	`
	expected := `
		grpc_client_message_sent_size_histogram_bytes{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.005"} 0
		grpc_client_message_sent_size_histogram_bytes{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.01"} 0
		grpc_client_message_sent_size_histogram_bytes{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.025"} 0
		grpc_client_message_sent_size_histogram_bytes{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.05"} 0
		grpc_client_message_sent_size_histogram_bytes{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.1"} 0
		grpc_client_message_sent_size_histogram_bytes{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.25"} 0
		grpc_client_message_sent_size_histogram_bytes{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="0.5"} 0
		grpc_client_message_sent_size_histogram_bytes{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="1"} 0
		grpc_client_message_sent_size_histogram_bytes{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="2.5"} 0
		grpc_client_message_sent_size_histogram_bytes{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="5"} 3
		grpc_client_message_sent_size_histogram_bytes{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="10"} 3
		grpc_client_message_sent_size_histogram_bytes{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service",le="+Inf"} 3
		grpc_client_message_sent_size_histogram_bytes_sum{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service"} 15
		grpc_client_message_sent_size_histogram_bytes_count{grpc_client_user_agent="fake-user-agent",grpc_is_fail_fast="true",grpc_method="Method",grpc_service="service"} 3
	`

	if err := testutil.CollectAndCompare(h, strings.NewReader(metadata+expected), "grpc_client_message_sent_size_histogram_bytes"); err != nil {
		t.Fatal(err)
	}
}
