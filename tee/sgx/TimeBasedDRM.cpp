#include <iostream>

#include "TimeBasedDRM.h"
#include "sgx_urts.h"
#include "sgx_uae_service.h"
#include "DRM_enclave_u.h"

using namespace std;

#define ENCLAVE_NAME    "DRM_enclave.signed.so"

TimeBasedDRM::TimeBasedDRM(void): enclave_id(0)
{
    sgx_status_t sgx_ret = SGX_ERROR_UNEXPECTED;
    sgx_ret = sgx_create_enclave(ENCLAVE_NAME, SGX_DEBUG_FLAG, NULL, NULL, &enclave_id, NULL);
    if (sgx_ret)
    {
        cerr << "cannot create enclave, error code = 0x" << hex << sgx_ret << endl;
    }
}


TimeBasedDRM::~TimeBasedDRM(void)
{
    if(enclave_id)
    {
        sgx_destroy_enclave(enclave_id);
    }
}

uint32_t TimeBasedDRM:: Init(uint8_t*  stored_time_based_policy)
{
    sgx_status_t sgx_ret = SGX_ERROR_UNEXPECTED;
    sgx_ps_cap_t ps_cap;
    memset(&ps_cap, 0, sizeof(sgx_ps_cap_t));
    sgx_ret = sgx_get_ps_cap(&ps_cap);
    if (sgx_ret)
    {
        cerr << "cannot get platform service capability, error code = 0x" << hex << sgx_ret << endl;
        return sgx_ret;
    }

    if (!SGX_IS_TRUSTED_TIME_AVAILABLE(ps_cap))
    {
        cerr << "trusted time is not supported" << endl;
        return SGX_ERROR_SERVICE_UNAVAILABLE;
    }

    uint32_t enclave_ret = 0;
    sgx_ret = create_time_based_policy(enclave_id, &enclave_ret, (uint8_t *)stored_time_based_policy, time_based_policy_length);
    if (sgx_ret)
    {
        cerr << "call create_time_based_policy fail, error code = 0x"<< hex << sgx_ret << endl;
        return sgx_ret;
    }

    if (enclave_ret)
    {
        cerr << "cannot create_time_based_policy, function return fail, error code = 0x" << hex << enclave_ret << endl;
        return enclave_ret;
    }

    return 0;
}


uint32_t TimeBasedDRM:: Init()
{
    return Init(time_based_policy);
}



uint32_t TimeBasedDRM::PerformFunction(uint8_t* stored_time_based_policy)
{
    sgx_status_t sgx_ret = SGX_ERROR_UNEXPECTED;
    uint32_t enclave_ret = 0;
    sgx_ret = perform_time_based_policy(enclave_id, &enclave_ret,
        stored_time_based_policy, time_based_policy_length);
    if (sgx_ret)
    {
        cerr<<"call perform_time_based_policy fail, error code = 0x"<< hex<<
            sgx_ret <<endl;
        return sgx_ret;
    }
    if (enclave_ret)
    {
        cerr<<"cannot perform_time_based_policy, function return fail, error code = 0x"
            << hex<< enclave_ret <<endl;
        return enclave_ret;
    }
    return 0;
}

uint32_t TimeBasedDRM::PerformFunction()
{
    return PerformFunction(time_based_policy);
}