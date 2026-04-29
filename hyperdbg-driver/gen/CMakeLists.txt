cmake_minimum_required(VERSION 3.15)
project(HyperDbgDriver C ASM_MASM)

set(CMAKE_C_STANDARD 11)
set(CMAKE_C_STANDARD_REQUIRED ON)

if(CMAKE_BUILD_TYPE STREQUAL "")
    set(CMAKE_BUILD_TYPE Debug)
endif()

# WDK/KMDF settings
set(CMAKE_MODULE_PATH "${CMAKE_CURRENT_LIST_DIR}")
find_package(WDK REQUIRED)

# Include directories
include_directories(
    ${CMAKE_CURRENT_SOURCE_DIR}
)

# Find all generated C and ASM files
file(GLOB_RECURSE SOURCES
    "*.c"
)

file(GLOB_RECURSE ASM_SOURCES
    "*.asm"
)

set(ALL_SOURCES ${SOURCES} ${ASM_SOURCES})

# Exclude external dependencies
list(FILTER SOURCES EXCLUDE REGEX ".*zydis.*")
list(FILTER SOURCES EXCLUDE REGEX ".*zycore.*")

# Build the driver using wdk_add_driver
wdk_add_driver(HyperDbgDriver ${ALL_SOURCES})
