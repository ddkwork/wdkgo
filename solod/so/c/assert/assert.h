#include <assert.h>

#ifdef NDEBUG
#define assert_Enabled false
#else
#define assert_Enabled true
#endif

#define assert_Assert(cond) assert(cond)
#define assert_Assertf(cond, msg) assert((cond) && msg)
