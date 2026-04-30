cmake_minimum_required(VERSION 3.28)

include("C:/Program Files/CMake/bin/ewdk.cmake")

project(HyperDbgUnified C ASM_MASM)
set_property(GLOBAL PROPERTY USE_FOLDERS ON)

set(WDK_NTDDI_VERSION "0x0A000008")

set(CMAKE_BUILD_TYPE Debug CACHE STRING "Build type" FORCE)
set(CMAKE_C_FLAGS_DEBUG "/Od /Zi /RTC1")
set(CMAKE_C_FLAGS "/Od /Zi /RTC1")

include_directories(
    "."
    "HyperDbg/hyperdbg/include"
    "HyperDbg/hyperdbg/dependencies/zydis/include"
    "HyperDbg/hyperdbg/dependencies/zydis/src"
    "HyperDbg/hyperdbg/dependencies/zydis/dependencies/zycore/include"
    "HyperDbg/hyperdbg/dependencies"
)

add_definitions(
    -DZYCORE_STATIC_BUILD
    -DZYDIS_STATIC_BUILD
)

file(GLOB ZydisSourceFiles "HyperDbg/hyperdbg/dependencies/zydis/src/*.c")
file(GLOB ZycoreSourceFiles "HyperDbg/hyperdbg/dependencies/zydis/dependencies/zycore/src/*.c")
set(ZydisSourceFiles ${ZydisSourceFiles} ${ZycoreSourceFiles})

file(GLOB_RECURSE HyperkdSourceFiles "HyperDbg/hyperdbg/hyperkd/code/*.c")
file(GLOB HyperkdAsmFiles "HyperDbg/hyperdbg/hyperkd/code/assembly/*.asm")

file(GLOB_RECURSE HyperhvSourceFiles "HyperDbg/hyperdbg/hyperhv/code/*.c")
list(FILTER HyperhvSourceFiles EXCLUDE REGEX ".*code/common/UnloadDll\\.c$")
file(GLOB HyperhvAsmFiles "HyperDbg/hyperdbg/hyperhv/code/assembly/*.asm")

file(GLOB HyperlogSourceFiles "HyperDbg/hyperdbg/hyperlog/code/*.c")

file(GLOB KdserialSourceFiles "HyperDbg/hyperdbg/kdserial/*.c")

file(GLOB ScriptEvalSourceFiles "HyperDbg/hyperdbg/script-eval/code/*.c")

file(GLOB HyperevadeSourceFiles "HyperDbg/hyperdbg/hyperevade/code/*.c")

file(GLOB_RECURSE HypertraceSourceFiles "HyperDbg/hyperdbg/hypertrace/code/*.c")

file(GLOB IncludeOptimizationsSourceFiles "HyperDbg/hyperdbg/include/components/optimizations/code/*.c")
file(GLOB IncludeSpinlockSourceFiles "HyperDbg/hyperdbg/include/components/spinlock/code/*.c")
file(GLOB IncludePlatformSourceFiles "HyperDbg/hyperdbg/include/platform/kernel/code/*.c")

km_lib(zydis_kernel ${ZydisSourceFiles})
target_compile_definitions(zydis_kernel PRIVATE ZYAN_NO_LIBC ZYDIS_NO_LIBC)
target_compile_options(zydis_kernel PRIVATE /Od /Zi /RTC1)

km_lib(Hyperhv ${HyperhvSourceFiles} ${HyperhvAsmFiles} ${IncludeOptimizationsSourceFiles} ${IncludeSpinlockSourceFiles})
target_include_directories(Hyperhv PRIVATE
    "HyperDbg/hyperdbg/hyperhv"
    "HyperDbg/hyperdbg/hyperhv/header"
    "HyperDbg/hyperdbg/hyperhv/code"
    "HyperDbg/hyperdbg/include/components/optimizations/code"
    "HyperDbg/hyperdbg/include/components/spinlock/code"
)
target_compile_definitions(Hyperhv PRIVATE DHYPERDBG_VMM HYPERDBG_HYPEREVADE HYPERDBG_HYPER_LOG HYPERDBG_HYPERTRACE NTDDI_VERSION=0x0A000008)
target_compile_options(Hyperhv PRIVATE /Od /Zi /RTC1)
target_precompile_headers(Hyperhv PRIVATE "HyperDbg/hyperdbg/hyperhv/pch.h")

km_lib(Hyperlog ${HyperlogSourceFiles})
target_include_directories(Hyperlog PRIVATE
    "HyperDbg/hyperdbg/hyperlog"
    "HyperDbg/hyperdbg/hyperlog/header"
    "HyperDbg/hyperdbg/hyperlog/code"
)
target_compile_definitions(Hyperlog PRIVATE DHYPERDBG_HYPERLOG HYPERDBG_HYPER_LOG HYPERDBG_HYPERTRACE HYPERDBG_HYPEREVADE HYPERDBG_VMM NTDDI_VERSION=0x0A000008)
target_compile_options(Hyperlog PRIVATE /Od /Zi /RTC1)
target_precompile_headers(Hyperlog PRIVATE "HyperDbg/hyperdbg/hyperlog/header/pch.h")

km_lib(Kdserial ${KdserialSourceFiles})
target_include_directories(Kdserial PRIVATE
    "HyperDbg/hyperdbg/kdserial"
)
target_compile_options(Kdserial PRIVATE /Od /Zi /RTC1)

km_lib(Hyperevade ${HyperevadeSourceFiles})
target_include_directories(Hyperevade PRIVATE
    "HyperDbg/hyperdbg/hyperevade"
    "HyperDbg/hyperdbg/hyperevade/header"
    "HyperDbg/hyperdbg/hyperevade/code"
)
target_compile_definitions(Hyperevade PRIVATE DHYPERDBG_HYPEREVADE HYPERDBG_HYPEREVADE HYPERDBG_HYPER_LOG HYPERDBG_HYPERTRACE HYPERDBG_VMM NTDDI_VERSION=0x0A000008)
target_compile_options(Hyperevade PRIVATE /Od /Zi /RTC1)
target_precompile_headers(Hyperevade PRIVATE "HyperDbg/hyperdbg/hyperevade/header/pch.h")

km_lib(Hypertrace ${HypertraceSourceFiles})
target_include_directories(Hypertrace PRIVATE
    "HyperDbg/hyperdbg/hypertrace"
    "HyperDbg/hyperdbg/hypertrace/header"
    "HyperDbg/hyperdbg/hypertrace/code"
)
target_compile_definitions(Hypertrace PRIVATE DHYPERDBG_HYPERTRACE HYPERDBG_HYPERTRACE HYPERDBG_HYPER_LOG HYPERDBG_HYPEREVADE HYPERDBG_VMM NTDDI_VERSION=0x0A000008)
target_compile_options(Hypertrace PRIVATE /Od /Zi /RTC1)
target_precompile_headers(Hypertrace PRIVATE "HyperDbg/hyperdbg/hypertrace/header/pch.h")

km_lib(Platform ${IncludePlatformSourceFiles})
target_include_directories(Platform PRIVATE
    "HyperDbg/hyperdbg/include/platform/kernel/code"
    "HyperDbg/hyperdbg/include/platform/kernel/header"
    "HyperDbg/hyperdbg/include"
)
target_compile_definitions(Platform PRIVATE NTDDI_VERSION=0x0A000008)
target_compile_options(Platform PRIVATE /Od /Zi /RTC1 /FIntddk.h /FIntstrsafe.h /FISDK/headers/Constants.h)

if(NOT EXISTS "${CMAKE_CURRENT_SOURCE_DIR}/hyperdbg-driver/gen")
km_sys(hyperkd
    KMDF 1.15
    ${HyperkdSourceFiles}
    ${HyperkdAsmFiles}
    ${IncludeOptimizationsSourceFiles}
    ${IncludeSpinlockSourceFiles}
    ${ScriptEvalSourceFiles}
)
target_include_directories(hyperkd PRIVATE
    "HyperDbg/hyperdbg/hyperkd"
    "HyperDbg/hyperdbg/hyperkd/header"
    "HyperDbg/hyperdbg/hyperkd/code"
    "HyperDbg/hyperdbg/include/components/optimizations/code"
    "HyperDbg/hyperdbg/include/components/spinlock/code"
    "HyperDbg/hyperdbg/script-eval"
)
target_compile_definitions(hyperkd PRIVATE DHYPERDBG_DEBUGGER HYPERDBG_HYPEREVADE HYPERDBG_HYPER_LOG HYPERDBG_HYPERTRACE HYPERDBG_VMM NTDDI_VERSION=0x0A000008)
target_precompile_headers(hyperkd PRIVATE "HyperDbg/hyperdbg/hyperkd/header/pch.h")
target_compile_options(hyperkd PRIVATE /Od /Zi /RTC1)
target_link_options(hyperkd PRIVATE /DEBUG /FORCE:MULTIPLE)

target_link_libraries(hyperkd
    zydis_kernel
    Hyperhv
    Hyperlog
    Kdserial
    Hyperevade
    Hypertrace
    Platform
)
endif()

if(EXISTS "${CMAKE_CURRENT_SOURCE_DIR}/hyperdbg-driver/gen")
    file(GLOB_RECURSE SolodDriverSourceFiles "hyperdbg-driver/gen/*.c")

    set(SolodHyperkdSourceFiles ${HyperkdSourceFiles})
    list(FILTER SolodHyperkdSourceFiles EXCLUDE REGEX ".*code/driver/Driver\\.c$")
    list(APPEND SolodHyperkdSourceFiles
        "HyperDbg/hyperdbg/hyperkd/code/common/Common.c"
        "HyperDbg/hyperdbg/hyperkd/code/debugger/broadcast/DpcRoutines.c"
    )

    km_sys(solod_hyperkd
        KMDF 1.15
        ${SolodDriverSourceFiles}
        ${SolodHyperkdSourceFiles}
        ${HyperkdAsmFiles}
        ${IncludeOptimizationsSourceFiles}
        ${IncludeSpinlockSourceFiles}
        ${ScriptEvalSourceFiles}
    )
    target_include_directories(solod_hyperkd PRIVATE
        "HyperDbg/hyperdbg/hyperkd"
        "HyperDbg/hyperdbg/hyperkd/header"
        "HyperDbg/hyperdbg/hyperkd/code"
        "HyperDbg/hyperdbg/include/components/optimizations/code"
        "HyperDbg/hyperdbg/include/components/spinlock/code"
        "HyperDbg/hyperdbg/script-eval"
        "hyperdbg-driver/gen"
        "solod/internal/compiler/builtin"
    )
    target_compile_definitions(solod_hyperkd PRIVATE DHYPERDBG_DEBUGGER HYPERDBG_HYPEREVADE HYPERDBG_HYPER_LOG HYPERDBG_HYPERTRACE HYPERDBG_VMM NTDDI_VERSION=0x0A000008)
    target_compile_options(solod_hyperkd PRIVATE /Od /Zi /RTC1)
    target_link_options(solod_hyperkd PRIVATE /DEBUG /FORCE:MULTIPLE)

    target_link_libraries(solod_hyperkd
        zydis_kernel
        Hyperhv
        Hyperlog
        Kdserial
        Hyperevade
        Hypertrace
        Platform
    )
endif()
