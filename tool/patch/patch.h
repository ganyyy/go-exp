#include <dlfcn.h>
#include <limits.h>
#include <stdlib.h>
#include <stdint.h>

static void* pluginOpen(const char *path, char **err)
{
    void *h = dlopen(path, RTLD_NOW | RTLD_GLOBAL);
    if (h == NULL)
    {
        *err = (char *)dlerror();
    }
    return h;
}


static void *pluginLookup(void *h, const char *name, char **err)
{
    void *r = dlsym(h, name);
    if (r == NULL)
    {
        *err = (char *)dlerror();
    }
    return r;
}


static int pluginClose(void *handle, char **err)
{
    int ret = dlclose(handle);
    if (ret != 0)
    {
        *err = (char *)dlerror();
    }
    return ret;
}