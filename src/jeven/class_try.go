package jeven

import (
	"fmt"
	"time"
)

type Car struct {
	Name string    //public
	Date time.Time //private
	desc string
}

// ToString 实例的方法
func (car *Car) ToString() {
	fmt.Printf("Name=%v, date=%v, desc=%v", car.Name, car.Date, car.desc)
}

func (car *Car) GetDesc() string {
	return car.desc
}
