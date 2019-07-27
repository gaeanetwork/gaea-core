#include <stdio.h>
#include <sgx_urts.h>
#include "sgx_error_code_message.h"

void print_error_message(sgx_status_t ret)
{
    int found = 0;

    for (size_t idx = 0; idx < ttl; idx++) {
        if(ret == sgx_errlist[idx].err) {
            if(NULL != sgx_errlist[idx].hint) {
                printf("Info: %s\n", sgx_errlist[idx].hint);
            }
            printf("Error: %s\n", sgx_errlist[idx].message);
            found = 1;
            break;
        }
    }

    if (!found) {
        printf("Error: Unexpected error occurred.\n");
    }
}
