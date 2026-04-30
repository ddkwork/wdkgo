cmake_minimum_required(VERSION 3.28)

include("C:/Program Files/CMake/bin/ewdk.cmake")

project(${PROJECT_NAME} C ASM_MASM)
set_property(GLOBAL PROPERTY USE_FOLDERS ON)

set(WDK_NTDDI_VERSION "0x0A000008")

set(CMAKE_BUILD_TYPE Debug CACHE STRING "Build type" FORCE)
set(CMAKE_C_FLAGS_DEBUG "/Od /Zi /RTC1")
set(CMAKE_C_FLAGS "/Od /Zi /RTC1")

include_directories(
    ${CMAKE_CURRENT_SOURCE_DIR}
)

add_definitions(
    -DZYCORE_STATIC_BUILD
    -DZYDIS_STATIC_BUILD
)

file(GLOB_RECURSE SOURCES "${CMAKE_CURRENT_SOURCE_DIR}/*.c")
file(GLOB_RECURSE ASM_SOURCES "${CMAKE_CURRENT_SOURCE_DIR}/*.asm")

list(FILTER SOURCES EXCLUDE REGEX ".*zydis.*")
list(FILTER SOURCES EXCLUDE REGEX ".*zycore.*")

set(ALL_SOURCES ${SOURCES} ${ASM_SOURCES})

km_sys(${PROJECT_NAME}
    KMDF 1.15
    ${ALL_SOURCES}
)
target_compile_definitions(${PROJECT_NAME} PRIVATE NTDDI_VERSION=0x0A000008)
target_compile_options(${PROJECT_NAME} PRIVATE /Od /Zi /RTC1)
target_link_options(${PROJECT_NAME} PRIVATE /DEBUG /FORCE:MULTIPLE)
