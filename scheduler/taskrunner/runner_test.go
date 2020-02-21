package taskrunner

import (
	"log"
	"testing"
	"time"
)

func TestRunner_StartAll(t *testing.T) {
	d := func(d dataChan) error {
		for i:=0;i<30;i++{
			d <- i
			log.Printf("data send:%d",i)
		}
		return nil
	}

	e := func(d dataChan) error {
		loopfor:
			for{
				select {
				case f := <-d:
					log.Printf("data received:%d",f)
					//return errors.New("Interrupt")
				default:
					break loopfor
				}
			}
			return nil
	}

	r := NewRunner(30,false,d,e)
	go r.StartAll()
	time.Sleep(3 * time.Second)
}
