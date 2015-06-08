SOURCES=multiboot.o screen.go.o screen.gox kernel.go.o

GOFLAGS= -nostdlib -nostdinc -fno-stack-protector -fno-split-stack -static -m32 -g -I.
GO=gccgo
ASFLAGS= -felf
NASM= nasm $(ASFLAGS)
OBJCOPY=objcopy

LDFLAGS=-T link.ld -m elf_i386
 

all: $(SOURCES) link

clean: 
	rm *.o *.gox kernel 

link:
	ld $(LDFLAGS) -o kernel $(SOURCES)

qemu:
	qemu-system-i386 -kernel ./kernel

%.gox: %.go.o
		$(OBJCOPY) -j .go_export $< $@

%.go.o: %.go
	$(GO)	$(GOFLAGS) -o $@ -c $<

%.o: %.s
	$(NASM) $<
