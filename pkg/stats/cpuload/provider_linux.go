//go:build linux

package cpuload

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type DataProvider struct {
	prev rawValue
}

const (
	statPath         = "/proc/stat"
	providerPlatform = "linux"
)

func NewDataProvider() *DataProvider {
	p := DataProvider{}

	if err := p.Available(); err == nil {
		// Первый запуск
		p.Data()
	}

	return &p
}

func (p *DataProvider) Available() error {
	_, err := os.Stat(statPath)
	if err != nil {
		return err
	}

	file, err := os.Open(statPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

func (p *DataProvider) Data() (Value, error) {
	file, err := os.Open(statPath)
	if err != nil {
		return Value{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var val Value
	for scanner.Scan() {
		v, err := p.parse(scanner.Text())
		if err == nil {
			val = v
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return Value{}, fmt.Errorf("can't read data: %w", err)
	}

	val.T = time.Now()

	return val, nil
}

func (p *DataProvider) parse(s string) (Value, error) {
	v, err := parse(s)
	if err != nil {
		return Value{}, err
	}

	// Предыдущее значение не корректно
	if p.prev.total == 0 || p.prev.total == v.total {
		p.prev = v

		return Value{}, errors.New("invalid previous value")
	}

	val := Value{
		User:   int(100_00 * (v.user - p.prev.user) / (v.total - p.prev.total)),
		System: int(100_00 * (v.system - p.prev.system) / (v.total - p.prev.total)),
		Idle:   int(100_00 * (v.idle - p.prev.idle) / (v.total - p.prev.total)),
	}

	p.prev = v

	return val, nil
}

type rawValue struct {
	user   int64
	system int64
	idle   int64
	total  int64
}

// Расчёт использования cpu в линуксе
// github.com/hightemp/docLinux/blob/master/articles/Как рассчитывается время и процент использования ЦП Linux.md .
func parse(s string) (rawValue, error) {
	if !strings.HasPrefix(s, "cpu ") {
		return rawValue{}, errors.New("invalid line")
	}

	vals := strings.Fields(s)
	if len(vals) != 11 {
		return rawValue{}, errors.New("can't split fields from string")
	}

	const user = "user"
	const nice = "nice"
	const system = "system"
	const idle = "idle"
	fields := map[string][2]int64{
		user:      {1, 0},
		nice:      {2, 0},
		system:    {3, 0},
		idle:      {4, 0},
		"iowait":  {5, 0},
		"irq":     {6, 0},
		"softirq": {7, 0},
		"steal":   {8, 0},
	}

	total := int64(0)
	for fieldName, f := range fields {
		pos := f[0]
		fieldVal, err := strconv.ParseInt(vals[pos], 10, 64)
		if err != nil {
			return rawValue{}, fmt.Errorf(
				"invalid field %s at pos %d minute ('%s'): %w",
				fieldName,
				pos,
				vals[pos],
				err,
			)
		}

		f[1] = fieldVal
		fields[fieldName] = f

		total += fieldVal
	}

	return rawValue{
		user:   fields[user][1] + fields[nice][1],
		system: fields[system][1],
		idle:   fields[idle][1],
		total:  total,
	}, nil
}
