package keylogger

// func TestReconnection(t *testing.T) {
// 	fd, err := os.CreateTemp("", "*")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer fd.Close()
// 	info, err := fd.Stat()
// 	fmt.Printf("%#v\n", info)
// 	fmt.Println(fd.Name())
// 	// try to create new keylogger with file descriptor which has the permission
// 	k, err := NewKeylogger(fd.Name())
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer k.Close()
// 	// run goroutine to receive keypress
// 	closedSig := make(chan int)
// 	go func() {
// 		for i := range k.Read() {
// 			fmt.Println(i)
// 		}
// 		fmt.Println("Out of loop")
// 		closedSig <- 1
// 	}()
// 	// test
// 	err = os.Remove(fd.Name())
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// block until decive disconnected or timeout
// 	select {
// 	case _ = <-closedSig:
// 		break
// 	case <-time.After(3 * time.Second):
// 		t.Fatal("Test listener timed out")
// 	}
// }
