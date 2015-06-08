MBOOT_PAGE_ALIGN    equ 1<<0
MBOOT_MEM_INFO      equ 1<<1
MBOOT_HEADER_MAGIC  equ 0x1BADB002
MBOOT_HEADER_FLAGS  equ MBOOT_PAGE_ALIGN | MBOOT_MEM_INFO
MBOOT_CHECKSUM      equ -(MBOOT_HEADER_MAGIC + MBOOT_HEADER_FLAGS)

[BITS 32]

[GLOBAL mboot]
[EXTERN code]
[EXTERN bss]
[EXTERN end]

mboot:
  dd    MBOOT_HEADER_MAGIC
  dd    MBOOT_HEADER_FLAGS
  dd    MBOOT_CHECKSUM
  dd    mboot
  dd    code
  dd    bss
  dd    end
  dd    start

[GLOBAL start]
extern go.kernel.Load

global __go_runtime_error ;gccgo compability
global __go_register_gc_roots ;gccgo compability
global __unsafe_get_addr ;convert uint32 to pointer

__unsafe_get_addr:
  push ebp
  mov ebp, esp
  mov eax, [ebp+8]
  mov esp, ebp
  pop ebp
  ret

start:
  push  ebx
  cli
  call  go.kernel.Load
  jmp   $

__go_register_gc_roots:
__go_runtime_error:
  ret
