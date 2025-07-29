#pragma once

#include <mach/mach.h>
#include <mach-o/dyld.h>
#include <string.h>
#include <stdint.h>
#include <unistd.h>
#include <stdlib.h>


uintptr_t get_image_slide(const char *image_name) {
    uint32_t count = _dyld_image_count();
    for (uint32_t i = 0; i < count; i++) {
        const char *name = _dyld_get_image_name(i);
        if (strcmp(name, image_name) == 0) {
            return _dyld_get_image_vmaddr_slide(i);
        }
    }
    return 0; // Return 0 if the image is not found
}


int patch_text(uintptr_t addr, size_t size, const char *data, const char **error) {
    vm_address_t start = trunc_page((vm_address_t)addr);
    vm_address_t end   = round_page((vm_address_t)addr + size);
    vm_size_t region_size = end - start;
    kern_return_t kr;

    kr = vm_protect(mach_task_self(), start, region_size, FALSE,
                    VM_PROT_READ|VM_PROT_WRITE|VM_PROT_COPY);
    if (kr != KERN_SUCCESS) {
        if (error) *error = "vm_protect -> WRITE failed";
        return -1;
    }

    memcpy((void *)addr, data, size);

    kr = vm_protect(mach_task_self(), start, region_size, FALSE,
                    VM_PROT_READ|VM_PROT_EXECUTE);
    if (kr != KERN_SUCCESS) {
        if (error) *error = "vm_protect -> RESTORE failed";
        return -1;
    }

    if (error) *error = NULL;
    return 0;
}