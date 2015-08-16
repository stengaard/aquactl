// +build linux

package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"syscall"
)

// Docs:  https://github.com/torvalds/linux/blob/master/Documentation/gpio/sysfs.txt

type sysFsPin struct {
	pinid uint
	name  string
	dir   *os.File
	val   *os.File
}

const (
	gpioPath = "/sys/class/gpio"
)

func newDigitalPin(pinid uint) (DigitalPin, error) {

	p := &sysFsPin{
		pinid: pinid,
		name:  fmt.Sprintf("gpio%d", pinid),
	}

	err := p.export(true)
	if err != nil {
		return nil, err
	}

	pinPath := path.Join(gpioPath, p.name)
	p.val, err = os.OpenFile(path.Join(pinPath, "value"), os.O_WRONLY|os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	p.dir, err = os.OpenFile(path.Join(pinPath, "direction"), os.O_WRONLY, 0666)
	if err != nil {
		p.dir.Close()
		return nil, err
	}

	return p, nil
}

func (p *sysFsPin) export(dir bool) error {
	ctlFile := "export"
	if !dir {
		ctlFile = "unexport"
	}
	f, err := os.OpenFile(path.Join(gpioPath, ctlFile), os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(f, "%d", p.pinid)
	if perr, ok := err.(*os.PathError); ok {
		// already exported during export - not an error
		if dir && perr.Err == syscall.EBUSY {
			return nil
		}
		// TODO: ok?
		// reserved and can't be unexported.
		//if !dir && perr.Err == syscall.EINVAL {
		//  return nil
		//}
	}
	return err
}

func read(f io.ReadSeeker) (string, error) {
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	_, err = f.Seek(0, 0)
	return string(buf), err
}

func write(f io.WriteSeeker, s string) error {
	_, err := f.Write([]byte(s))
	if err != nil {
		return err
	}
	_, err = f.Seek(0, 0)
	return err
}

func (p *sysFsPin) SetDirection(newDir PinDirection) (PinDirection, error) {
	curDir, err := read(p.dir)
	if newDir != PinIn && newDir != PinOut && err == nil {
		err = errors.New("illegal pin direction")
	}
	err = write(p.dir, string(newDir))

	return PinDirection(curDir), err
}

func (p *sysFsPin) Close() error {
	p.dir.Close() // ignoring errs
	p.val.Close()
	return p.export(false)
}

func (p *sysFsPin) WriteBool(v bool) error {
	b := "0"
	if v {
		b = "1"
	}
	return write(p.val, b)
}

func (p *sysFsPin) ReadBool() (bool, error) {
	out, err := read(p.val)
	if err != nil {
		return false, err
	}
	switch out {
	case "1":
		return true, nil
	case "0":
		return false, nil
	}

	return false, errors.New("bad value")
}
