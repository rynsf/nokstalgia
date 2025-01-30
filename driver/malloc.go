package driver

import (
	"log"
)

type tBlock *blockMeta
type blockMeta struct {
	size uint32
	ptr  uint32
	prev tBlock
	next tBlock
	free bool
}

var globalHead tBlock

func firstFit(s uint32) tBlock {
	current := globalHead
	for current != nil {
		if current.free && current.size >= s {
			return current
		}
		current = current.next
	}
	return nil
}

func split(b tBlock, s uint32) {
	new := &blockMeta{
		size: b.size - s,
		ptr:  b.ptr + s,
		prev: b,
		next: b.next,
		free: true,
	}
	b.size = s
	b.next = new
	if new.next != nil {
		new.next.prev = new
	}
}

func Malloc(s uint32) uint32 {
	s = (((s + 3) >> 2) << 2) // word align
	b := firstFit(s)
	if b != nil {
		if b.size-s >= 4 { // is the block large enough to split
			split(b, s)
		}
		b.free = false
	} else {
		log.Println("malloc: can't allocate space, out of memory")
		return 0
	}
	return b.ptr
}

func findBlock(ptr uint32) tBlock {
	for current := globalHead; current != nil; current = current.next {
		if ptr == current.ptr && !current.free {
			return current
		}
	}
	return nil
}

func fusion(b tBlock) tBlock {
	if b.next != nil && b.next.free {
		b.size += b.next.size
		b.next = b.next.next
		if b.next != nil {
			b.next.prev = b
		}
	}
	return b
}

func Free(ptr uint32) {
	b := findBlock(ptr)
	if b != nil {
		b.free = true
		if b.prev != nil && b.prev.free {
			b = fusion(b.prev)
		}
		if b.next != nil {
			fusion(b)
		}
	}
}

func InitDynamicMem(base, size uint32) {
	globalHead = &blockMeta{
		size: size,
		ptr:  base,
		prev: nil,
		next: nil,
		free: true,
	}
}
