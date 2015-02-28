package main

/*
# include <malloc.h>
static void *gl_alloc (void *ptr, size_t osize, size_t nsize) {
	(void)osize;
	if (nsize == 0) {
		free(ptr);
		return NULL;
	}
	else
		return realloc(ptr, nsize);
}
*/
import "C"

import "../lua"
import "unsafe"
import "fmt"

//var refHolder = map[unsafe.Pointer][]byte{}
var heap []byte
var refHolder [][]byte

//a terrible allocator!
//meant to be illustrative of the mechanics,
//not usable as an actual implementation
func AllocatorF(ptr unsafe.Pointer, osize uint, nsize uint) unsafe.Pointer {
	//fmt.Printf("Allocation %v from (%v, %v)\n", nsize, osize, len(heap))
	if nsize == 0 {
		heap = make([]byte, 0)
		ptr = unsafe.Pointer(nil)
	} else if osize != nsize {
		var tmp []byte
		tmp, heap = heap, make([]byte, int(nsize))
		copy(heap, tmp)
		heap = make([]byte, int(nsize))

		ptr = unsafe.Pointer(&heap[0])
		//fmt.Println(ptrToSlice(ptr, int(nsize)))
	}
	//fmt.Println("in allocf");
	return ptr
}

func AllocatorC(ptr unsafe.Pointer, osize uint, nsize uint) unsafe.Pointer {
	if nsize == 0 {
		ptr = unsafe.Pointer(nil)
	} else if osize != nsize {
		ptr = unsafe.Pointer(C.realloc(ptr, C.size_t(nsize)))
		//fmt.Println(ptrToSlice(ptr, int(nsize)))
	}
	return ptr
}

func AllocatorOld(ptr unsafe.Pointer, osize uint, nsize uint) unsafe.Pointer {
	if nsize == 0 {
		//TODO: remove from reference holder
	} else if osize != nsize {
		//TODO: remove old ptr from list if its in there
		slice := make([]byte, nsize)
		ptr = unsafe.Pointer(&(slice[0]))
		//TODO: add slice to holder
		l := len(refHolder)
		refHolder = refHolder[0 : l+1]
		refHolder[l] = slice
	}
	//fmt.Println("in allocf");
	return ptr
}

func ptrToSlice(ptr unsafe.Pointer, length int) []byte {
	addr := uintptr(ptr)
	var sl = struct {
		addr uintptr
		len  int
		cap  int
	}{addr, length, length}
	b := *(*[]byte)(unsafe.Pointer(&sl))
	return b
}

func adder(L *lua.State) int {
	a := L.ToInteger(1)
	b := L.ToInteger(2)
	L.PushInteger(int64(a + b))
	return 1
}

func main() {

	heap = make([]byte, 0, 50000)
	refHolder = make([][]byte, 0, 500)

	fmt.Println("Creating state with allocator")
	L := lua.NewStateAlloc(AllocatorOld)
	defer L.Close()
	fmt.Println("Opening libs")
	L.OpenLibs()
	fmt.Println("Updating allocator")

	L.Register("adder", adder)
	for i := 0; i < 10; i++ {
		L.GetField(lua.LUA_GLOBALSINDEX, "print")
		L.PushString("Hello World!")
		L.Call(1, 0)
	}

	L.DoString("print(adder(4, 4))")
}
