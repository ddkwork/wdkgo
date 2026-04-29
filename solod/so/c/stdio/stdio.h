#include <stdio.h>

#define stdio_EOF EOF

#define stdio_SeekSet SEEK_SET
#define stdio_SeekCur SEEK_CUR
#define stdio_SeekEnd SEEK_END

#define stdio_File FILE

#define stdio_Stdin stdin
#define stdio_Stdout stdout
#define stdio_Stderr stderr

#define stdio_Fopen(path, mode) fopen(path, mode)
#define stdio_Fclose(stream) fclose(stream)
#define stdio_Fflush(stream) fflush(stream)

#define stdio_Fseek(stream, offset, whence) fseek(stream, offset, whence)
#define stdio_Ftell(stream) ftell(stream)

#define stdio_Fgetc(stream) fgetc(stream)
#define stdio_Fputc(ch, stream) fputc(ch, stream)

#define stdio_Fgets(s, n, stream) (so_byte*)fgets((char*)s, n, stream)
#define stdio_Fputs(s, stream) fputs(s, stream)

#define stdio_Fread(ptr, size, count, stream) \
    ((so_int)fread(ptr, (size_t)size, (size_t)count, stream))
#define stdio_Fwrite(ptr, size, count, stream) \
    ((so_int)fwrite(ptr, (size_t)size, (size_t)count, stream))

#define stdio_Feof(stream) feof(stream)
#define stdio_Ferror(stream) ferror(stream)

#define stdio_Printf(format, ...) printf(format, ##__VA_ARGS__)
#define stdio_Fprintf(stream, format, ...) fprintf(stream, format, ##__VA_ARGS__)
#define stdio_Snprintf(buf, size, format, ...) snprintf((char*)buf, size, format, ##__VA_ARGS__)

#define stdio_Scanf(format, ...) scanf(format, ##__VA_ARGS__)
#define stdio_Fscanf(stream, format, ...) fscanf(stream, format, ##__VA_ARGS__)
#define stdio_Sscanf(s, format, ...) sscanf(s, format, ##__VA_ARGS__)
