//go:build windows

package cpuload

import (
	"errors"
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

type DataProvider struct {
	prev rawValue
}

const providerPlatform = "windows"

func NewDataProvider() *DataProvider {
	p := DataProvider{}

	if err := p.Available(); err == nil {
		// Первый запуск
		p.Data()
	}

	return &p
}

// https://learn.microsoft.com/en-us/windows/win32/api/minwinbase/ns-minwinbase-filetime
type FILETIME struct {
	DwLowDateTime  uint32
	DwHighDateTime uint32
}

func (f FILETIME) Int64() int64 {
	// DwLowDateTime + 2**32 * DwHighDateTime
	return int64(f.DwLowDateTime) + int64(4294967296)*int64(f.DwHighDateTime)
}

var (
	DllKernel32             = windows.NewLazySystemDLL("kernel32.dll")
	Win32ProcGetSystemTimes = DllKernel32.NewProc("GetSystemTimes")
)

type rawValue struct {
	user   int64
	system int64
	idle   int64
	total  int64
}

func getSystemTimes() (t rawValue, err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		if rErr, ok := r.(error); ok {
			err = errors.Join(err, rErr)
			return
		}

		err = errors.Join(err, fmt.Errorf("getSystemTimes failed: %v", r))
	}()

	var lpIdleTime FILETIME
	var lpKernelTime FILETIME
	var lpUserTime FILETIME

	// https://learn.microsoft.com/ru-ru/windows/win32/api/processthreadsapi/nf-processthreadsapi-getsystemtimes
	r, _, _ := Win32ProcGetSystemTimes.Call(
		uintptr(unsafe.Pointer(&lpIdleTime)),
		uintptr(unsafe.Pointer(&lpKernelTime)),
		uintptr(unsafe.Pointer(&lpUserTime)))
	if r == 0 {
		return rawValue{}, windows.GetLastError()
	}

	idle := lpIdleTime.Int64()
	kernel := lpKernelTime.Int64()
	user := lpUserTime.Int64()

	system := kernel - idle

	return rawValue{
		user:   user,
		system: system,
		idle:   idle,
		total:  user + system + idle,
	}, nil
}

func (p *DataProvider) Available() error {
	_, err := getSystemTimes()
	if err != nil {
		return err
	}

	return nil
}

func (p *DataProvider) Data() (Value, error) {
	v, err := getSystemTimes()
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
