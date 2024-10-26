//go:build linux

package loadavg

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	loadavg_path      = "/proc/loadavg"
	provider_platform = "linux"
)

type DataProvider struct{}

func (p *DataProvider) Available() error {
	_, err := os.Stat(loadavg_path)
	if err != nil {
		return err
	}

	file, err := os.Open(loadavg_path)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

func (p *DataProvider) Data() (Value, error) {
	file, err := os.Open(loadavg_path)
	if err != nil {
		return Value{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return Value{}, fmt.Errorf("can't read data: %w", err)
	}

	t := time.Now()

	v, err := parse(scanner.Text())
	if err != nil {
		return Value{}, fmt.Errorf("can't parse data: %w", err)
	}

	v.T = t

	return v, nil
}

func parse(s string) (Value, error) {
	s = strings.ReplaceAll(s, ".", "")
	v := strings.Fields(s)

	if len(v) < 3 {
		return Value{}, errors.New("can't split fields from string")
	}

	one, err := strconv.ParseInt(v[0], 10, 64)
	if err != nil {
		return Value{}, fmt.Errorf("invalid field 1 minute ('%s'): %w", v[0], err)
	}

	five, err := strconv.ParseInt(v[1], 10, 64)
	if err != nil {
		return Value{}, fmt.Errorf("invalid field 5 minutes ('%s'): %w", v[1], err)
	}

	fifteen, err := strconv.ParseInt(v[2], 10, 64)
	if err != nil {
		return Value{}, fmt.Errorf("invalid field 15 minutes  ('%s'): %w", v[2], err)
	}

	return Value{
		One:     int(one),
		Five:    int(five),
		Fifteen: int(fifteen),
	}, nil
}
