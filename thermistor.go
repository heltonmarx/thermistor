package main

import (
	"fmt"
	"math"
	"strconv"
	"syscall"
)

//
//	Beaglebone Black read ADC pin (VREF 1.8)
//
func ReadADC(pin int) (float64, error) {
	path := fmt.Sprintf("/sys/devices/ocp.3/helper.15/AIN%d", pin)

	fd, err := syscall.Open(path, syscall.O_RDONLY, 0)
	if err != nil {
		return -1, err
	}
	defer syscall.Close(fd)
	var buf []byte

	n, err := syscall.Read(fd, buf)
	if err != nil || n <= 0 {
		return -1, err
	}
	t, err := strconv.Atoi(string(buf))
	return float64(t), err
}

//
//	Steinhart Hart Equation
//		Lower Temp C: 		-55
//		Midrange Temp C:	47.5
//		Upper Temp C:		150
//		Resistance @ 25° C:	10.000
//		Thermistor Curve:	Y (-3.9%/C @ 25C) Mil Ratio A
//
// 		link: [ http://www.thermistor.com/calculators?r=sheccr ]
//
func ReadTemperature(pin int) (float64, error) {
	const (
		a      = 0.000830194603253
		b      = 0.000265071184936
		c      = 0.000000108900012
		vcc    = 1.8
		r1     = 1000
		factor = 1000
	)

	var (
		rth    float64
		vsense float64
		temp   float64
		adc    float64
	)

	adc, err := ReadADC(pin)
	if err != nil {
		return -1, err
	}
	vsense = adc / factor

	rth = (float64(r1) / ((vcc / vsense) - 1))
	temp = math.Log(float64(rth))
	temp = 1 / (a + (b * temp) + (c * temp * temp * temp))
	temp = temp - 273.15

	return temp, nil
}
