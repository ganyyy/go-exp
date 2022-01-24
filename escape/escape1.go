package main

type Elem struct {
	t [65535]bool
}

func main() {
	var sendChan = make(chan *Elem, 1)
	go func() {
		sendChan <- &Elem{}
	}()

	var v = <-sendChan
	println(v.t[100])
}
