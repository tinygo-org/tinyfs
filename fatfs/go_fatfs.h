#include <stdlib.h>
#include "diskio.h"
#include "ff.h"

//DSTATUS go_fatfs_disk_status(void* ptr);
//DSTATUS go_fatfs_disk_initialize(void* ptr);
extern DRESULT go_fatfs_disk_read(void* drv, void* buff, DWORD sector, UINT count);
extern DRESULT go_fatfs_disk_write(void* drv, void* buff, DWORD sector, UINT count);
extern DRESULT go_fatfs_disk_ioctl(void* drv, BYTE cmd, DWORD* param);

extern DWORD go_fatfs_get_fattime();

// Helper functions used to allocate new FatFs objects, needed because TinyGo
// does not support sizeof() yet
FATFS* go_fatfs_new_fatfs(void);
FIL* go_fatfs_new_fil(void);
FF_DIR* go_fatfs_new_ff_dir(void);

//struct lfs_config* go_lfs_new_lfs_config(void);
//lfs_dir_t* go_lfs_new_lfs_dir(void);
