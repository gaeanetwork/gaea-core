#ifndef _SEALED_DATA_DEFINES_H_
#define _SEALED_DATA_DEFINES_H_

#include "sgx_error.h"

#define PLATFORM_SERVICE_DOWNGRADED  0xF001

#define REPLAY_DETECTED              0xF002
#define MAX_RELEASE_REACHED          0xF003

/* equal to sgx_calc_sealed_data_size(0,sizeof(replay_protected_pay_load))) */ 
#define SEALED_REPLAY_PROTECTED_PAY_LOAD_SIZE 620
#define REPLAY_PROTECTED_PAY_LOAD_MAX_RELEASE_VERSION 5

#define TIMESOURCE_CHANGED           0xF004
#define TIMESTAMP_UNEXPECTED         0xF005
#define LEASE_EXPIRED                0xF006

/* equal tosgx_calc_sealed_data_size(0,sizeof(time_based_pay_load))) */ 
#define TIME_BASED_PAY_LOAD_SIZE          624
#define TIME_BASED_LEASE_DURATION_SECOND  3

#endif