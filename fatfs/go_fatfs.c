#include "go_fatfs.h"

// implementation of disk interface layer defined in diskio.h

DRESULT disk_read(void *drv, BYTE* buff, DWORD sector, UINT count) {
    return go_fatfs_disk_read(drv, buff, sector, count);
}

DRESULT disk_write(void *drv, const BYTE* buff, DWORD sector, UINT count) {
    return go_fatfs_disk_write(drv, (BYTE*)buff, sector, count);
}

DRESULT disk_ioctl(void *drv, BYTE cmd, void* buff) {
    return go_fatfs_disk_ioctl(drv, cmd, buff);
}

DWORD get_fattime() {
    return go_fatfs_get_fattime();
}

// Helper functions for creating FatFs structs

FATFS* go_fatfs_new_fatfs(void) {
    return malloc(sizeof(FATFS));
}

FIL* go_fatfs_new_fil(void) {
    return malloc(sizeof(FIL));
}

FF_DIR* go_fatfs_new_ff_dir(void) {
    return malloc(sizeof(FF_DIR));
}

// if ffconf.h has FF_FS_READONLY set, certain functions aren't implemented,
// which prevents the Go code from linking properly.
#if FF_FS_READONLY == 1

FRESULT f_mkfs (FATFS *fs, BYTE opt, DWORD au, void* work, UINT len) {
    return 99;
}

FRESULT f_mkdir (FATFS *fs, const TCHAR* path) {
    return 99;
}

FRESULT f_unlink (FATFS *fs, const TCHAR* path) {
    return 99;
}

FRESULT f_write (FIL* fp, const void* buff, UINT btw, UINT* bw) {
    return 99;
}

FRESULT f_sync (FIL* fp) {
    return 99;
}

FRESULT f_rename (FATFS *fs, const TCHAR* path_old, const TCHAR* path_new) {
    return 99;
}

FRESULT f_getfree (FATFS *fs, DWORD* nclst) {
    return 99;
}

#endif