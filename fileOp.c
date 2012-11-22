#include "fileOp.h"
#include <stdio.h>
#include <fcntl.h>

#define USING_DYNAMIC_PLUGIN 1	//notice: if not using plugin, we need to change t.totalSize


#define _FILE_OFFSET_BITS 64

void ShowMessage(char* str)
{
	printf("%s\n", str);
	MessageBoxA(NULL, str, "", MB_OK);
}

typedef int (__cdecl *fnReadAt)(int, void*, int, __int64);
typedef int (__cdecl *fnWriteAt)(int, void*, int, __int64 );
typedef int (__cdecl *fnOpen)(const char*, int);
typedef int (__cdecl *fnClose)(int);

fnReadAt  readAtFn = NULL;
fnWriteAt writeAtFn = NULL;
fnOpen    openFn = NULL;
fnClose   closeFn = NULL;

int Init(){
#if USING_DYNAMIC_PLUGIN	
	HINSTANCE hinstLib = LoadLibrary("p2pdll.dll"); 
	if (hinstLib == NULL){
		ShowMessage("load dll failed");
		return -1;
	}

	readAtFn = (fnReadAt)GetProcAddress(hinstLib, "ReadAt");
	writeAtFn = (fnWriteAt)GetProcAddress(hinstLib, "WriteAt");
	openFn = (fnOpen)GetProcAddress(hinstLib, "Open");
	closeFn = (fnClose)GetProcAddress(hinstLib, "Close");
	if (readAtFn == NULL || writeAtFn == NULL || openFn == NULL || closeFn == NULL){
		ShowMessage("some of export function not found");
		return -1;
	}
#endif
	return 0;
}


int ReadAt(int fd, void* buf, int len, __int64 offset){
#if USING_DYNAMIC_PLUGIN	
	return readAtFn(fd, buf, len, offset);
#else
	int n = _lseeki64(fd, offset, SEEK_SET);
	if (n == -1)
	{
		printf("_lseeki64 errno %d\n", errno);
		ShowMessage("oops _lseeki64 failed");
	}
	n = _read(fd, buf, len);

	//printf("ReadAt %I64d\n", offset);

	if (n <= 0)
	{
		printf("errno %d\n", errno);
		ShowMessage("oops");
	}

	return n;	
#endif	
}

int WriteAt(int fd, void* buf, int len, __int64 offset){
#if USING_DYNAMIC_PLUGIN	
	return writeAtFn(fd, buf, len, offset);
#else
	int n = _lseeki64(fd, offset, SEEK_SET);
	if (n == -1)
	{
		printf("_lseeki64 errno %d\n", errno);
		ShowMessage("oops");
	}
	n = _write(fd, buf, len);

	//printf("WriteAt %I64d\n", offset);

	if (n <= 0)
	{
		printf("errno %d\n", errno);
		ShowMessage("oops");
	}	

	return n;
#endif	
}

int Open(const char* fname, int oflag){
#if USING_DYNAMIC_PLUGIN	
	return openFn(fname, oflag);
#else
	return _open(fname, oflag);
#endif	
}

void Close(int fd){
	ShowMessage("close file");
#if USING_DYNAMIC_PLUGIN		
	closeFn(fd);
#else
	_close(fd);
#endif	
}
