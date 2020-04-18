#include <stdlib.h>
#include "lfs.h"

// LittleFS uses function pointers to callback functions in order to allow for
// a user-provided abstraction of a block device.  Go does not allow for passing
// Go function pointers to C, so instead we will use global Go functions with
// exported symbols in our C callback functions.  A pointer to the Go struct
// representing each LittleFS instance is saved in the lfs_config.context field,
// which is then used by those Go functions to dispatch callback invocations.
extern int go_lfs_block_device_read(void*, lfs_block_t, lfs_off_t, void*, lfs_size_t);
extern int go_lfs_block_device_prog(void*, lfs_block_t, lfs_off_t, const void*, lfs_size_t);
extern int go_lfs_block_device_erase(void*, lfs_block_t);
extern int go_lfs_block_device_sync(void*);

// These are the global C callbacks. Pointers to these functions are passed to
// the LittleFS library as the block device callbacks, and they in turn call
// the associated global Go callbacks which handle dispatching the callback
// invocations to the correct instance of the LFS struct in Go.
int go_lfs_c_cb_read(const struct lfs_config *c, lfs_block_t block, lfs_off_t off, void *buffer, lfs_size_t size);
int go_lfs_c_cb_prog(const struct lfs_config *c, lfs_block_t block, lfs_off_t off, const void *buffer, lfs_size_t size);
int go_lfs_c_cb_erase(const struct lfs_config *c, lfs_block_t block);
int go_lfs_c_cb_sync(const struct lfs_config *c);

// Helper functions used to allocate new LFS objects, needed because TinyGo
// does not support sizeof() yet
struct lfs* go_lfs_new_lfs(void);
struct lfs_config* go_lfs_new_lfs_config(void);
lfs_dir_t* go_lfs_new_lfs_dir(void);
lfs_file_t* go_lfs_new_lfs_file(void);

// Helper function to set the function pointers to the global callbacks on a
// provided LFS config struct
struct lfs_config* go_lfs_set_callbacks(struct lfs_config *cfg);
