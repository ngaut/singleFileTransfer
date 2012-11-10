#include "fileOp.h"
#include <stdio.h>
#include <fcntl.h>


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
	HINSTANCE hinstLib = LoadLibrary("transfer.dll"); 
	if (hinstLib == NULL){
		ShowMessage("load dll failed");
	}

	readAtFn = (fnReadAt)GetProcAddress(hinstLib, "ReadAt");
	writeAtFn = (fnWriteAt)GetProcAddress(hinstLib, "WriteAt");
	openFn = (fnOpen)GetProcAddress(hinstLib, "Open");
	closeFn = (fnClose)GetProcAddress(hinstLib, "Close");

	return 0;
}


int ReadAt(int fd, void* buf, int len, __int64 offset){
/*
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
	*/
	return readAtFn(fd, buf, len, offset);
}

int WriteAt(int fd, void* buf, int len, __int64 offset){
/*
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
	*/
	return writeAtFn(fd, buf, len, offset);
}

int Open(const char* fname, int oflag){
/*
	return _open(fname, oflag);
	*/
	return openFn(fname, oflag);
}

void Close(int fd){
/*
	ShowMessage("close file");
	_close(fd);
	*/
	closeFn(fd);
}
