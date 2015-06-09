;gccgo compability
global __go_runtime_error 
global __go_register_gc_roots 
global __go_type_hash_error
global __go_type_equal_error

__go_runtime_error:
__go_register_gc_roots: 
__go_type_hash_error:
__go_type_equal_error:
  ret

global __unsafe_get_addr ;convert uint32 to pointer

__unsafe_get_addr:
  push ebp
  mov ebp, esp
  mov eax, [ebp+8]
  mov esp, ebp
  pop ebp
  ret

global __asm_mov_to_cr3

__asm_mov_to_cr3:
  push ebp
  mov ebp, esp
  mov eax, [ebp+8]
  mov eax, cr3
  pop ebp
  ret

global __asm_mov_from_cr0

__asm_mov_from_cr0:
  push ebp
  mov eax, cr0
  pop ebp
  push eax
  ret

global __asm_mov_to_cr0:

__asm_mov_to_cr0:
  push ebp
  mov ebp, esp
  mov eax, [ebp+8]
  mov eax, cr0
  pop ebp
  ret

extern go.runtime.New

global __go_new
global __go_new_nopointers

__go_new_nopointers:
__go_new:
  call go.runtime.New
  ret

extern go.runtime.StringPlus

global __go_string_plus
__go_string_plus:
  call go.runtime.StringPlus
  ret

extern go.runtime.TypeEqualIdentity

global __go_type_equal_identity
__go_type_equal_identity:
  call go.runtime.TypeEqualIdentity
  ret

extern go.runtime.TypeHashIdentity

global __go_type_hash_identity
__go_type_hash_identity:
  call go.runtime.TypeHashIdentity
  ret

extern go.screen.PrintString 

global __go_print_string
__go_print_string:
  call go.screen.PrintString
  ret

extern go.screen.PrintNl

global __go_print_nl
__go_print_nl:
  call go.screen.PrintNl
  ret

extern go.runtime.ByteArrayToString

global __go_byte_array_to_string
__go_byte_array_to_string:
  call go.runtime.ByteArrayToString
  ret
