package memory

//extern end
var end uint32

var (
	placement_address uint32
	isInit            bool = false
)

//extern __unsafe_get_addr
func pointer2uint32(pointer *uint32) uint32

func _Init() {
	placement_address = pointer2uint32(&end)
	InitPaging()
}

//HEAP

func kmalloc(size uint32, align int, phys *uint32) uint32 {
	if !isInit {
		_Init()
	}
	if align == 1 && (placement_address&0xFFFFF000) != uint32(0) {
		placement_address &= 0xFFFFF000
		placement_address += 0x1000
	}
	if phys != nil {
		*phys = placement_address
	}
	res := placement_address
	placement_address += size
	return res
}

func Kmalloc(size uint32) uint32 {
	return kmalloc(size, 0, nil)
}

func KmallocPhys(size uint32, phys *uint32) uint32 {
	return kmalloc(size, 0, phys)
}

func KmallocPhysPage(size uint32, phys *uint32) uint32 {
	return kmalloc(size, 1, phys)
}

func KmallocPage(size uint32) uint32 {
	return kmalloc(size, 1, nil)
}

//PAGING

type page struct {
	Present   bool
	RW        bool
	User      bool
	Accessed  bool
	Writed    bool
	Unused    uint32
	FrameAddr uint32
}

type pageTable [1024]page

type pageDir struct {
	pageTables [1024]*pageTable
	tablesPhys [1024]uint32
	physAddr   uint32
}

const (
	sizeOfpageDir   = 12296
	sizeOfpageTable = 16384
)

var (
	nframes   uint32
	frames    []uint32
	kernelDir *pageDir
	curDir    *pageDir
)

func frameCalc(frameAddr uint32) (idx, offset uint32) {
	frame := frameAddr / 0x1000
	idx = frame / 32
	offset = frame - idx
	return idx, offset
}

func setFrame(frameAddr uint32) {
	idx, offset := frameCalc(frameAddr)
	frames[idx] = (0x1 << offset)
}

func clearFrame(frameAddr uint32) {
	idx, offset := frameCalc(frameAddr)
	frames[idx] &= 0xFFFFFFFF - (0x1 << offset)
}

func testFrame(frameAddr uint32) uint32 {
	idx, offset := frameCalc(frameAddr)
	return (frames[idx] & (0x1 << offset))
}

func firstFrame() uint32 {
	for i := uint32(0); i < uint32(nframes/32); i++ {
		if frames[i] != 0xFFFFFFFF {
			for j := uint32(0); j < 32; j++ {
				toTest := uint32(0x1 << j)
				if (frames[i] & toTest) == 0 {
					return i*4*8 + j
				}
			}
		}
	}
	return 0
}

func AllocFrame(Page *page, isKernel bool, isRW bool) {
	if Page.FrameAddr != 0 {
		return
	} else {
		frame := firstFrame()
		if frame == 0 {
			//TODO Add  panic
			return
		}
		setFrame(frame * 0x1000)
		Page.Present = true
		Page.RW = isRW
		Page.User = !isKernel
		Page.FrameAddr = frame
	}
}

func FreeFrame(Page *page) {
	if Page.FrameAddr == 0 {
		return
	} else {
		clearFrame(Page.FrameAddr)
		Page.FrameAddr = 0
	}
}

//extern __unsafe_get_addr
func pointer2uint32slice(pointer uint32) *[]uint32

//extern __unsafe_get_addr
func pointer2pageDir(pointer uint32) *pageDir

//extern __unsafe_get_addr
func pointer2pageTable(pointer uint32) *pageTable

func getPage(addr uint32, isMake bool, PageDir *pageDir) *page {
	addr /= 0x1000
	tableIdx := uint32(addr / 1024)
	if PageDir.pageTables[tableIdx] != nil {
		return &(PageDir.pageTables[tableIdx][addr-tableIdx])
	} else if isMake {
		var phys *uint32
		table_p := KmallocPhysPage(sizeOfpageTable, phys)
		PageDir.pageTables[tableIdx] = pointer2pageTable(table_p)
		PageDir.tablesPhys[tableIdx] = *phys | uint32(0x07)
		return &(PageDir.pageTables[tableIdx][addr-tableIdx])
	}
	return nil
}

func InitPaging() {
	const MemEndPage = 0x1000000
	nframes = MemEndPage / 0x1000
	frames_p := Kmalloc(nframes / 32)
	frames = *(pointer2uint32slice(frames_p))
	kernelDir_p := Kmalloc(sizeOfpageDir)
	kernelDir = pointer2pageDir(kernelDir_p)
	curDir = kernelDir
	for i := uint32(0); i < placement_address; i += 0x1000 {
		AllocFrame(getPage(i, true, kernelDir), true, true)
	}

	//TODO register interrupt

	switchPageDir(kernelDir)
}

//extern __asm_mov_to_cr3
func movToCr3(PageDir *pageDir)

//extern __asm_mov_from_cr0
func movFromCr0() uint32

//extern __asm_mov_to_cr0
func movToCr0(cr0 uint32)

func switchPageDir(PageDir *pageDir) {
	curDir = PageDir
	movToCr3(curDir)
	cr0 := movFromCr0()
	cr0 |= 0x80000000
	movToCr0(cr0)
}
