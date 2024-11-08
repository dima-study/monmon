package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v1 "github.com/dima-study/monmon/pkg/api/proto/stats/v1"
)

func main() {
	hostport := flag.String("connect", "localhost:50051", "gRPC service connection string")
	stat := flag.String("stat", "", "stat name to receive")
	period := flag.String("period", "1", "period in seconds how ofter get stats")
	interval := flag.String("interval", "1", "interval length in seconds of stats")

	flag.Parse()

	switch {
	case hostport == nil || *hostport == "":
		flag.Usage()
		exitError("connect argument is empty")
	case stat == nil || *stat == "":
		listProviders(*hostport)
	default:
		switch {
		case period == nil || *period == "":
			flag.Usage()
			exitError("period argument is empty")
		case interval == nil || *interval == "":
			flag.Usage()
			exitError("interval argument is empty")
		}

		getStat(*hostport, *stat, *interval, *period)
	}
}

func exitError(err any) {
	fmt.Fprintf(os.Stderr, "\nerror: %s\n", err)
	os.Exit(1)
}

func listProviders(connect string) {
	conn, err := grpc.NewClient(connect, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		exitError(err)
	}
	defer conn.Close()

	client := v1.NewStatsServiceClient(conn)

	resp, err := client.GetSupportedStats(context.Background(), &v1.GetSupportedStatsRequest{})
	if err != nil {
		exitError(fmt.Errorf("can't get supported stats: %w", err))
	}

	fmt.Println("Supported stats:")
	for _, p := range resp.Providers {
		fmt.Printf("  %s\n", p.ProviderID)

		if state := p.AvailabilityDetails.GetState(); state != v1.Available {
			fmt.Printf("    [%s] %s (%s)\n", state, p.ProviderName, p.Platform)

			if p.AvailabilityDetails.Details != "" {
				fmt.Printf("      %s\n", p.AvailabilityDetails.Details)
			}
		} else {
			fmt.Printf("    %s (%s)\n", p.ProviderName, p.Platform)
		}
	}
}

func getStat(connect string, provider string, interval string, period string) {
	// TODO: impl
}
