#pragma once

#include <ntifs.h>
#include <ntstrsafe.h>
typedef UCHAR uint8_t;
typedef CHAR int8_t;
typedef USHORT uint16_t;
typedef SHORT int16_t;
typedef ULONG uint32_t;
typedef LONG int32_t;
typedef ULONGLONG uint64_t;
typedef LONGLONG int64_t;
#include <malloc.h>
#include <intrin.h>

#pragma intrinsic(__readcr0)
#pragma intrinsic(__readcr3)
#pragma intrinsic(__readcr4)
#pragma intrinsic(__writecr3)
#pragma intrinsic(__writecr4)
#pragma intrinsic(__vmx_on)
#pragma intrinsic(__vmx_off)
#pragma intrinsic(__vmx_vmclear)
#pragma intrinsic(__vmx_vmlaunch)
#pragma intrinsic(__vmx_vmread)
#pragma intrinsic(__vmx_vmresume)
#pragma intrinsic(__vmx_vmwrite)
#pragma intrinsic(__vmx_vmptrld)
#pragma intrinsic(__readgsqword)

#ifndef __cplusplus
typedef unsigned char bool;
#define true 1
#define false 0
#endif

static inline uint8_t AsmVmxVmread(uint64_t Field, uint64_t* FieldValue) {
    return __vmx_vmread(Field, FieldValue);
}
static inline uint8_t AsmVmxVmwrite(uint64_t Field, uint64_t FieldValue) {
    return __vmx_vmwrite(Field, FieldValue);
}
static inline uint8_t AsmVmxVmread32(uint64_t Field, uint32_t* FieldValue) {
    uint64_t val;
    uint8_t ret = __vmx_vmread(Field, &val);
    *FieldValue = (uint32_t)val;
    return ret;
}
static inline uint8_t AsmVmxVmwrite32(uint64_t Field, uint32_t FieldValue) {
    return __vmx_vmwrite(Field, (uint64_t)FieldValue);
}
static inline uint8_t AsmVmxVmlaunch(void) {
    return __vmx_vmlaunch();
}
static inline uint8_t AsmVmxVmresume(void) {
    return __vmx_vmresume();
}
static inline uint8_t AsmVmxVmxClear(uint64_t PhysicalAddr) {
    return __vmx_vmclear(PhysicalAddr);
}
static inline uint8_t AsmVmxVmxPtrld(uint64_t PhysicalAddr) {
    return __vmx_vmptrld(PhysicalAddr);
}
static inline void AsmVmxVmxOff(void) {
    __vmx_off();
}

#ifndef ExAllocatePool2
extern PVOID ExAllocatePool2(POOL_FLAGS Flags, SIZE_T NumberOfBytes, ULONG Tag);
#endif

extern char* PsGetProcessImageFileName(void* Process);

#ifndef so_MaxAllocaSize
#define so_MaxAllocaSize (64 << 10)
#endif

typedef uint8_t so_byte;
typedef int32_t so_rune;
typedef int64_t so_int;
typedef uint64_t so_uint;

typedef struct {
    const char* ptr;
    size_t len;
} so_String;

#define so_str(s) ((so_String){s, sizeof(s) - 1})

typedef struct {
    void* ptr;
    size_t len;
    size_t cap;
} so_Slice;

extern so_byte so_Nil[];

static inline const char* so_cstr_raw(so_String s, char* buf) {
    if (s.len > 0) memcpy(buf, s.ptr, s.len);
    buf[s.len] = '\0';
    return buf;
}

static inline bool so_string_eq(so_String s1, so_String s2) {
    return s1.len == s2.len && (s1.len == 0 || memcmp(s1.ptr, s2.ptr, s1.len) == 0);
}

static inline so_String so_string_add(so_String a, so_String b) {
    size_t n = a.len + b.len;
    char* p = (char*)_alloca(n);
    if (a.len > 0) memcpy(p, a.ptr, a.len);
    if (b.len > 0) memcpy(p + a.len, b.ptr, b.len);
    so_String r;
    r.ptr = p;
    r.len = n;
    return r;
}

#define so_len(s) ((so_int)((s).len))
#define so_cap(s) ((so_int)((s).cap))

#define so_at(T, s, i) (((T*)(s).ptr)[i])
#define so_at_ptr(T, s, i) ((T*)(s).ptr + (i))

#define so_min(a, b) ((a) < (b) ? (a) : (b))
#define so_max(a, b) ((a) > (b) ? (a) : (b))

#define so_alloca(size) (_alloca(size))

so_rune so_utf8_decode(so_String s, so_int i, so_int* w);
size_t so_utf8_encode(so_rune r, char* buf);

struct so_Error_ {
    const char* msg;
};
typedef struct so_Error_* so_Error;

#define errors_New(s) (&(struct so_Error_){s})

static inline so_Error fmt_Errorf(so_String format, so_Slice args) {
    return &(struct so_Error_){format.ptr};
}

#define SO_GLOBAL_ERROR(name, msg) static struct so_Error_ _err_##name = {msg}; so_Error name = &_err_##name

static inline so_String fmt_Sprintf(so_String format, so_Slice args) {
    return format;
}

#define so_panic(msg)                                     \
    do {                                                  \
        DbgPrint("panic: %s\n  %s:%d (func %s)\n",       \
                msg, __FILE__, __LINE__, __func__);       \
        __debugbreak();                                   \
    } while (0)

static inline so_Slice so_make_slice_impl(size_t elem_size, so_int len, so_int cap) {
    size_t n = (size_t)elem_size * (size_t)cap;
    void* p = n ? _alloca(n) : so_Nil;
    if (n) memset(p, 0, n);
    so_Slice s;
    s.ptr = p;
    s.len = (size_t)len;
    s.cap = (size_t)cap;
    return s;
}
#define so_make_slice(T, len, cap) so_make_slice_impl(sizeof(T), (len), (cap))

static inline so_Slice so_slice_impl(so_Slice s, size_t elem_size, size_t from, size_t to) {
    if (to > s.cap || from > to)
        so_panic("slice bounds out of range");
    so_Slice r;
    r.ptr = (char*)s.ptr + from * elem_size;
    r.len = to - from;
    r.cap = s.cap - from;
    return r;
}
#define so_slice(T, s, from, to) so_slice_impl((s), sizeof(T), (size_t)(from), (size_t)(to))

static inline so_Slice so_append_impl(so_Slice s, const void* val, size_t elem_size, size_t n) {
    if (s.len + n > s.cap) so_panic("append: out of capacity");
    memcpy((char*)s.ptr + s.len * elem_size, val, n * elem_size);
    so_Slice r;
    r.ptr = s.ptr;
    r.len = s.len + n;
    r.cap = s.cap;
    return r;
}
#define so_append(T, s, ...) so_append_impl((s), (T[]){__VA_ARGS__}, sizeof(T), sizeof((T[]){__VA_ARGS__})/sizeof(T))

static inline so_Slice so_extend_impl(so_Slice dst, so_Slice src, size_t elem_size) {
    if (dst.len + src.len > dst.cap)
        so_panic("extend: out of capacity");
    if (src.len > 0) memcpy((char*)dst.ptr + dst.len * elem_size, src.ptr, src.len * elem_size);
    so_Slice r;
    r.ptr = dst.ptr;
    r.len = dst.len + src.len;
    r.cap = dst.cap;
    return r;
}
#define so_extend(T, dst, src) so_extend_impl((dst), (src), sizeof(T))

static inline so_int so_copy_impl(so_Slice dst, so_Slice src, size_t elem_size) {
    size_t n = dst.len < src.len ? dst.len : src.len;
    if (n > 0) memmove(dst.ptr, src.ptr, n * elem_size);
    return (so_int)n;
}
#define so_copy(T, dst, src) so_copy_impl((dst), (src), sizeof(T))

static inline so_int so_copy_string(so_Slice dst, so_String src) {
    size_t n = dst.len < src.len ? dst.len : src.len;
    if (n > 0) memmove(dst.ptr, src.ptr, n);
    return (so_int)n;
}

#define so_array_slice(T, arr, from, to, max) \
    ((so_Slice){(T*)(arr) + (from), (to) - (from), (max) - (from)})

#define so_string_bytes(s) \
    ((so_Slice){(void*)((s).ptr), (s).len, (s).len})

#define so_bytes_string(bs) \
    ((so_String){(const char*)((bs).ptr), (bs).len})

static inline const char* so_cstr(so_String s) {
    char* buf = (char*)_alloca(s.len + 1);
    if (s.len > 0) memcpy(buf, s.ptr, s.len);
    buf[s.len] = '\0';
    return buf;
}

static inline so_String so_string_slice(so_String s, so_int from, so_int to) {
    if (to > s.len || from > to) so_panic("slice bounds out of range");
    so_String r;
    r.ptr = s.ptr + from;
    r.len = to - from;
    return r;
}

#define unsafe_Sizeof(x) sizeof(x)

static inline void* unsafe_Add(void* ptr, size_t offset) {
    return (char*)ptr + offset;
}

static inline so_String unsafe_String(void* ptr, size_t len) {
    if (ptr == NULL) return (so_String){(char*)&so_Nil, 0};
    so_String r;
    r.ptr = (char*)ptr;
    r.len = len;
    return r;
}

static inline so_byte* unsafe_StringData(so_String s) {
    if (s.len == 0) return NULL;
    return (so_byte*)s.ptr;
}

static inline so_Slice unsafe_Slice(void* ptr, size_t len) {
    if (ptr == NULL) return (so_Slice){&so_Nil, 0, 0};
    so_Slice r;
    r.ptr = ptr;
    r.len = len;
    r.cap = len;
    return r;
}

static inline void* unsafe_SliceData(so_Slice s) {
    if (s.cap == 0) return NULL;
    return s.ptr;
}

#define so_clear(T, s) do { \
    so_Slice _so_clr_s = (s); \
    memset(_so_clr_s.ptr, 0, _so_clr_s.len * sizeof(T)); \
} while (0)

#define so_decay(s) ((s).cap ? (s).ptr : NULL)

#define slices_Contains(SliceT, ElemT, s, v) so_slices_Contains_##ElemT((s), (v))

static inline bool so_slices_Contains_uint32_t(so_Slice s, uint32_t v) {
    uint32_t* p = (uint32_t*)s.ptr;
    for (size_t i = 0; i < s.len; i++) {
        if (p[i] == v) return true;
    }
    return false;
}

typedef struct {
    void* keys;
    void* vals;
    uint8_t* used;
    size_t len;
    size_t cap;
} so_Map;

static inline uint64_t so_key_hash_default(const void* ptr, size_t n, uint64_t seed) {
    const uint8_t* p = (const uint8_t*)ptr;
    uint64_t h = seed;
    for (size_t i = 0; i < n; i++) {
        h ^= p[i];
        h *= 0x100000001b3ULL;
    }
    return h;
}

static inline uint64_t so_key_hash_string(const void* ptr, size_t n, uint64_t seed) {
    (void)n;
    const so_String* s = (const so_String*)ptr;
    return so_key_hash_default(s->ptr, s->len, seed);
}

static inline bool so_key_eq_default(const void* a, const void* b, size_t n) {
    return memcmp(a, b, n) == 0;
}

static inline bool so_key_eq_string(const void* a, const void* b, size_t n) {
    (void)n;
    return so_string_eq(*(const so_String*)a, *(const so_String*)b);
}

static inline size_t so_map_nextpow2(size_t n) {
    if (n == 0) return 1;
    n--;
    n |= n >> 1;
    n |= n >> 2;
    n |= n >> 4;
    n |= n >> 8;
    n |= n >> 16;
    n |= n >> 32;
    return n + 1;
}

static inline size_t so_map_cap(size_t n) {
    if (n == 0) return 0;
    return so_map_nextpow2(n + n / 3 + 1);
}

static inline so_Map* so_make_map_impl(size_t ksz, size_t vsz, size_t n) {
    if (n == 0) n = 1;
    size_t cap = so_map_cap(n);
    size_t ktotal = ksz * cap;
    size_t vtotal = vsz * cap;
    size_t utotal = sizeof(uint8_t) * cap;
    void* kp = _alloca(ktotal);
    void* vp = _alloca(vtotal);
    uint8_t* up = (uint8_t*)_alloca(utotal);
    if (kp) memset(kp, 0, ktotal);
    if (vp) memset(vp, 0, vtotal);
    if (up) memset(up, 0, utotal);
    so_Map* mp = (so_Map*)_alloca(sizeof(so_Map));
    mp->keys = kp;
    mp->vals = vp;
    mp->used = up;
    mp->len = 0;
    mp->cap = cap;
    return mp;
}
#define so_make_map(K, V, n) so_make_map_impl(sizeof(K), sizeof(V), (n))

static inline void so_map_set_impl(so_Map* m, const void* key, const void* val,
    size_t ksz, size_t vsz,
    uint64_t(*hash_fn)(const void*, size_t, uint64_t),
    bool(*eq_fn)(const void*, const void*, size_t)) {
    uint64_t h = hash_fn(key, ksz, (uintptr_t)m);
    size_t mask = m->cap - 1;
    size_t step = (size_t)(h >> 32) | 1;
    size_t idx = (size_t)h & mask;
    for (size_t p = 0; ; p++) {
        if (p >= m->cap) so_panic("map: out of capacity");
        if (!m->used[idx]) {
            memcpy((char*)m->keys + idx * ksz, key, ksz);
            memcpy((char*)m->vals + idx * vsz, val, vsz);
            m->used[idx] = 1;
            m->len++;
            break;
        }
        if (eq_fn((const char*)m->keys + idx * ksz, key, ksz)) {
            memcpy((char*)m->vals + idx * vsz, val, vsz);
            break;
        }
        idx = (idx + step) & mask;
    }
}

static inline void* so_map_get_impl(const so_Map* m, const void* key, void* val,
    size_t ksz, size_t vsz,
    uint64_t(*hash_fn)(const void*, size_t, uint64_t),
    bool(*eq_fn)(const void*, const void*, size_t)) {
    memset(val, 0, vsz);
    if (m->cap == 0) return val;
    uint64_t h = hash_fn(key, ksz, (uintptr_t)m);
    size_t mask = m->cap - 1;
    size_t step = (size_t)(h >> 32) | 1;
    size_t idx = (size_t)h & mask;
    for (size_t p = 0; p < m->cap; p++) {
        if (!m->used[idx]) break;
        if (eq_fn((const char*)m->keys + idx * ksz, key, ksz)) {
            memcpy(val, (const char*)m->vals + idx * vsz, vsz);
            break;
        }
        idx = (idx + step) & mask;
    }
    return val;
}

#define so_map_set(K, V, m, key, val) do { \
    K _so_map_k = (key); \
    V _so_map_v = (val); \
    so_map_set_impl((m), &_so_map_k, &_so_map_v, sizeof(K), sizeof(V), \
        so_key_hash_default, so_key_eq_default); \
} while (0)

#define so_map_set_strkey(K, V, m, key, val) do { \
    K _so_map_k = (key); \
    V _so_map_v = (val); \
    so_map_set_impl((m), &_so_map_k, &_so_map_v, sizeof(K), sizeof(V), \
        so_key_hash_string, so_key_eq_string); \
} while (0)

#define so_map_get(K, V, m, key) \
    (*(V*)so_map_get_impl((m), &(K){(key)}, (V[]){0}, sizeof(K), sizeof(V), \
        so_key_hash_default, so_key_eq_default))

#define so_map_get_strkey(K, V, m, key) \
    (*(V*)so_map_get_impl((m), &(K){(key)}, (V[]){0}, sizeof(K), sizeof(V), \
        so_key_hash_string, so_key_eq_string))

#define so_map_lit(K, V, n, keys, vals) \
    so_map_lit_impl(sizeof(K), sizeof(V), (n), (keys), (vals))

static inline so_Map* so_map_lit_impl(size_t ksz, size_t vsz, size_t n,
    const void* keys, const void* vals) {
    so_Map* m = so_make_map_impl(ksz, vsz, n);
    for (size_t i = 0; i < n; i++) {
        so_map_set_impl(m,
            (const char*)keys + i * ksz,
            (const char*)vals + i * vsz,
            ksz, vsz, so_key_hash_default, so_key_eq_default);
    }
    return m;
}

static inline bool so_map_has_impl(const so_Map* m, const void* key,
    size_t ksz,
    uint64_t(*hash_fn)(const void*, size_t, uint64_t),
    bool(*eq_fn)(const void*, const void*, size_t)) {
    if (m->cap == 0) return false;
    uint64_t h = hash_fn(key, ksz, (uintptr_t)m);
    size_t mask = m->cap - 1;
    size_t step = (size_t)(h >> 32) | 1;
    size_t idx = (size_t)h & mask;
    for (size_t p = 0; p < m->cap; p++) {
        if (!m->used[idx]) return false;
        if (eq_fn((const char*)m->keys + idx * ksz, key, ksz)) return true;
        idx = (idx + step) & mask;
    }
    return false;
}

#define so_map_has(K, m, key) \
    so_map_has_impl((m), &(K){(key)}, sizeof(K), \
        so_key_hash_default, so_key_eq_default)

static inline void so_map_remove_impl(so_Map* m, const void* key,
    size_t ksz, size_t vsz,
    uint64_t(*hash_fn)(const void*, size_t, uint64_t),
    bool(*eq_fn)(const void*, const void*, size_t)) {
    if (m->cap == 0) return;
    uint64_t h = hash_fn(key, ksz, (uintptr_t)m);
    size_t mask = m->cap - 1;
    size_t step = (size_t)(h >> 32) | 1;
    size_t idx = (size_t)h & mask;
    for (size_t p = 0; p < m->cap; p++) {
        if (!m->used[idx]) return;
        if (eq_fn((const char*)m->keys + idx * ksz, key, ksz)) {
            m->used[idx] = 0;
            m->len--;
            return;
        }
        idx = (idx + step) & mask;
    }
}

#define so_map_remove(K, m, key) do { \
    so_map_remove_impl((m), &(K){(key)}, sizeof(K), 0, \
        so_key_hash_default, so_key_eq_default); \
} while (0)

typedef struct {
    bool val;
    so_Error err;
} so_R_bool_err;
typedef struct {
    double val;
    so_Error err;
} so_R_f64_err;
typedef struct {
    float val;
    so_Error err;
} so_R_f32_err;
typedef struct {
    int32_t val;
    so_Error err;
} so_R_i32_err;
typedef struct {
    int64_t val;
    so_Error err;
} so_R_i64_err;
typedef struct {
    so_byte val;
    so_Error err;
} so_R_byte_err;
typedef struct {
    so_int val;
    so_Error err;
} so_R_int_err;
typedef struct {
    so_rune val;
    so_Error err;
} so_R_rune_err;
typedef struct {
    so_Slice val;
    so_Error err;
} so_R_slice_err;
typedef struct {
    so_String val;
    so_Error err;
} so_R_str_err;
typedef struct {
    so_uint val;
    so_Error err;
} so_R_uint_err;
typedef struct {
    uint32_t val;
    so_Error err;
} so_R_u32_err;
typedef struct {
    uint64_t val;
    so_Error err;
} so_R_u64_err;
typedef struct {
    void* val;
    so_Error err;
} so_R_ptr_err;

typedef struct {
    bool val;
    bool val2;
} so_R_bool_bool;
typedef struct {
    bool val;
    so_int val2;
} so_R_bool_int;
typedef struct {
    double val;
    bool val2;
} so_R_f64_bool;
typedef struct {
    double val;
    double val2;
} so_R_f64_f64;
typedef struct {
    double val;
    so_int val2;
} so_R_f64_int;
typedef struct {
    float val;
    bool val2;
} so_R_f32_bool;
typedef struct {
    int64_t val;
    int32_t val2;
} so_R_i64_i32;
typedef struct {
    so_int val;
    bool val2;
} so_R_int_bool;
typedef struct {
    so_int val;
    so_int val2;
} so_R_int_int;
typedef struct {
    so_int val;
    uint64_t val2;
} so_R_int_u64;
typedef struct {
    so_rune val;
    bool val2;
} so_R_rune_bool;
typedef struct {
    so_rune val;
    so_int val2;
} so_R_rune_int;
typedef struct {
    so_String val;
    bool val2;
} so_R_str_bool;
typedef struct {
    so_String val;
    so_String val2;
} so_R_str_str;
typedef struct {
    so_uint val;
    so_uint val2;
} so_R_uint_uint;
typedef struct {
    uint32_t val;
    bool val2;
} so_R_u32_bool;
typedef struct {
    uint32_t val;
    so_int val2;
} so_R_u32_int;
typedef struct {
    uint32_t val;
    uint32_t val2;
} so_R_u32_u32;
typedef struct {
    uint64_t val;
    bool val2;
} so_R_u64_bool;
typedef struct {
    uint64_t val;
    so_int val2;
} so_R_u64_int;
typedef struct {
    uint64_t val;
    uint64_t val2;
} so_R_u64_u64;
typedef struct {
    so_Slice val;
    so_int val2;
} so_R_slice_int;
typedef struct {
    void* val;
    bool val2;
} so_R_ptr_bool;

static inline bool atomic_CompareAndSwapUint32(uint32_t* addr, uint32_t oldVal, uint32_t newVal) {
    return (uint32_t)InterlockedCompareExchange((volatile LONG*)addr, (LONG)newVal, (LONG)oldVal) == oldVal;
}

static inline void atomic_StoreUint32(uint32_t* addr, uint32_t val) {
    InterlockedExchange((volatile LONG*)addr, (LONG)val);
}

static inline uint32_t atomic_LoadUint32(uint32_t* addr) {
    return (uint32_t)InterlockedCompareExchange((volatile LONG*)addr, 0, 0);
}

static inline uint32_t atomic_AddUint32(uint32_t* addr, uint32_t delta) {
    return (uint32_t)InterlockedAdd((volatile LONG*)addr, (LONG)delta);
}
