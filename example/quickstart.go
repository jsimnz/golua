package main

import (
	"../lua"
	//"fmt"
	"unsafe"
)

var heap []byte

func adder(L *lua.State) int {
	a := L.ToInteger(1)
	b := L.ToInteger(2)
	L.PushInteger(int64(a + b))
	return 1
}

func AllocatorF(ptr unsafe.Pointer, osize uint, nsize uint) unsafe.Pointer {
	//fmt.Printf("Allocation %v from (%v, %v)\n", nsize, osize, len(heap))
	if nsize == 0 {
		heap = make([]byte, 0)
		ptr = unsafe.Pointer(nil)
	} else if osize != nsize {
		var tmp []byte
		tmp, heap = heap, make([]byte, int(nsize))
		copy(heap, tmp)
		_ = tmp // gc will clean it up
		ptr = unsafe.Pointer(&(heap[0]))
	}
	//fmt.Println("in allocf");
	return ptr
}

func main() {
	heap = make([]byte, 0, 50000)
	L := lua.NewStateAlloc(AllocatorF)
	defer L.Close()
	L.OpenLibs()

	L.GetField(lua.LUA_GLOBALSINDEX, "print")
	L.PushString("Hello World!")
	L.Call(1, 0)

	//fmt.Println(heap)
	L.Register("adder", adder)
	//fmt.Println(heap)
	L.DoString("print(adder(2, 2))")
	//fmt.Println(heap)
}
