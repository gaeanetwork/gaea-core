#include <stdio.h>
#include <sgx_urts.h>
#include "sgx_error_code_message.h"

void print_error_message(sgx_status_t ret)
{
    size_t idx = 0;
    size_t ttl = sizeof sgx_errlist/sizeof sgx_errlist[0];

    for (idx = 0; idx < ttl; idx++) {
        if(ret == sgx_errlist[idx].err) {
            if(NULL != sgx_errlist[idx].hint)
                printf("Info: %s\n", sgx_errlist[idx].hint);
            printf("Error: %s\n", sgx_errlist[idx].message);
            break;
        }
    }

    if (idx == ttl) {
        printf("Error: Unexpected error occurred.\n");
    }
}
