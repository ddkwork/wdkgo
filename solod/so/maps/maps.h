#include <time.h>

#if defined(__APPLE__)
#include <stdlib.h>
#elif defined(__linux__)
#include <sys/random.h>
#endif

// seed returns a random 64-bit seed for hash randomization.
static inline uint64_t maps_seed(void) {
    uint64_t seed = 0;
#if defined(__APPLE__)
    arc4random_buf(&seed, sizeof(seed));
#elif defined(__linux__)
    if (getrandom(&seed, sizeof(seed), 0) != sizeof(seed)) {
        // Fallback to time-based seed.
        struct timespec ts;
        clock_gettime(CLOCK_MONOTONIC, &ts);
        seed ^= (uint64_t)ts.tv_nsec ^ (uint64_t)ts.tv_sec;
    }
#else
    seed = (uint64_t)time(NULL) ^ (uintptr_t)&seed;
#endif
    return seed;
}

// keyHash hashes a key, dispatching to string or inline hash.
#define maps_keyHash(K, key_ptr, seed) _Generic((K){0}, \
    so_String: maps_hashString(key_ptr, seed),          \
    default: maps_hash(key_ptr, sizeof(K), seed))

// equal compares two typed key pointers for equality.
#define maps_keyEqual(K, a, b)                                       \
    _Generic((K){0},                                                 \
        so_String: so_string_eq(*(so_String*)(a), *(so_String*)(b)), \
        default: memcmp((a), (b), sizeof(K)) == 0)
