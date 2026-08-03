package main

/*
#cgo LDFLAGS: -framework ServiceManagement -framework Foundation
void registerLoginItemObjC(void);
*/
import "C"

func registerLoginItem() {
	C.registerLoginItemObjC()
}
