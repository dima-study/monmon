package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v1 "github.com/dima-study/monmon/pkg/api/proto/stats/v1"
)

func main() {
	hostport := flag.String("connect", "localhost:50051", "gRPC service connection string")
	stat := flag.String("stat", "", "stat name to receive")
	period := flag.Int64("period", 1, "period in seconds how ofter get stats")
	interval := flag.Int64("interval", 1, "interval length in seconds of stats")

	flag.Parse()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-shutdown
		cancel()
	}()

	switch {
	case hostport == nil || *hostport == "":
		flag.Usage()
		exitError("connect argument is empty")
	case stat == nil || *stat == "":
		listProviders(ctx, *hostport)
	default:
		switch {
		case period == nil || *period < 1:
			flag.Usage()
			exitError("period argument must be >=1")
		case interval == nil || *interval < 1:
			flag.Usage()
			exitError("interval argument must be >=1")
		}

		getStat(ctx, *hostport, *stat, *interval, *period)
	}
}

func exitError(err any) {
	fmt.Fprintf(os.Stderr, "\nerror: %s\n", err)
	os.Exit(1)
}

func getConn(connect string) *grpc.ClientConn {
	conn, err := grpc.NewClient(connect, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		exitError(err)
	}

	return conn
}

func getProviders(ctx context.Context, connect string) []*v1.Provider {
	conn := getConn(connect)
	defer conn.Close()

	client := v1.NewStatsServiceClient(conn)

	resp, err := client.GetSupportedStats(ctx, &v1.GetSupportedStatsRequest{})
	if err != nil {
		exitError(fmt.Errorf("can't get supported stats: %w", err))
	}

	return resp.Providers
}

func getAvailableProvider(ctx context.Context, connect string, providerID string) *v1.Provider {
	providers := getProviders(ctx, connect)
	for _, p := range providers {
		if p.ProviderID != providerID {
			continue
		}

		if state := p.AvailabilityDetails.GetState(); state != v1.Available {
			exitError(
				fmt.Errorf("provider %s is not available: %s %s", providerID, state, p.AvailabilityDetails.Details),
			)
		}

		return p
	}

	exitError(fmt.Errorf("server doesn't support provider %s", providerID))

	return nil
}

func listProviders(ctx context.Context, connect string) {
	providers := getProviders(ctx, connect)

	fmt.Println("Supported stats:")
	for _, p := range providers {
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

func getStat(ctx context.Context, connect string, providerID string, interval int64, period int64) {
	_ = getAvailableProvider(ctx, connect, providerID)

	conn := getConn(connect)
	defer conn.Close()

	client := v1.NewStatsServiceClient(conn)

	req := v1.GetStatsRequest{
		Interval: interval,
		Period:   period,
	}
	stream, err := client.GetStats(ctx, &req)
	if err != nil {
		exitError(fmt.Errorf("can't get supported stats: %w", err))
	}

	out := strings.Builder{}
	for {
		resp, err := stream.Recv()

		// err может быть вызван отменой контекста, првоеряем
		select {
		case <-ctx.Done():
			return
		default:
		}

		if err != nil {
			exitError(fmt.Errorf("error while receive: %w", err))
		}

		rec := resp.Record
		if rec.Provider.ProviderID != providerID {
			continue
		}

		args := []any{}

		out.WriteString("at %s === [%s | %s]\n")
		args = append(
			args,
			rec.Time.AsTime().String(),
			rec.Provider.ProviderName,
			rec.Provider.Platform,
		)

		for _, v := range rec.Value {
			out.WriteString("%s: %s\n") // value.Name: value.Value
			args = append(args, v.Name, v.Value)
		}

		log.Printf(out.String()+"\n", args...)
		out.Reset()
	}
}
