#include <stdlib.h>

#define stdlib_ExitSuccess EXIT_SUCCESS
#define stdlib_ExitFailure EXIT_FAILURE

#define stdlib_Exit(status) exit(status)

#define stdlib_Malloc(size) malloc(size)
#define stdlib_Calloc(count, size) calloc(count, size)
#define stdlib_Realloc(ptr, size) realloc(ptr, size)
#define stdlib_Free(ptr) free(ptr)

#define stdlib_Atoi(str) atoi(str)
#define stdlib_Atof(str) atof(str)

#define stdlib_Getenv(name) (uint8_t*)getenv(name)
