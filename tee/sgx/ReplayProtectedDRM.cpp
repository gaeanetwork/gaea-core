#include <iostream>
#include "sgx.h"
#include "sgx_urts.h"
#include "sgx_uae_service.h"
#include "DRM_enclave_u.h"
#include "ReplayProtectedDRM.h"
#include "sgx_tseal.h"

using namespace std;

#define ENCLAVE_NAME    "DRM_enclave.signed.so"


ReplayProtectedDRM::ReplayProtectedDRM(): enclave_id(0)
{
    sgx_status_t sgx_ret = SGX_ERROR_UNEXPECTED;
    sgx_ret = sgx_create_enclave(ENCLAVE_NAME, SGX_DEBUG_FLAG, NULL, NULL, &enclave_id, NULL);
    if (sgx_ret)
    {
        cerr << "cannot create enclave, error code = 0x" << hex << sgx_ret << endl;
    }
}

ReplayProtectedDRM::~ReplayProtectedDRM(void)
{
    if(enclave_id)
    {
        sgx_destroy_enclave(enclave_id);
    }
}

uint32_t ReplayProtectedDRM:: Init(uint8_t*  stored_sealed_activity_log)
{
    sgx_status_t sgx_ret = SGX_ERROR_UNEXPECTED;
    sgx_ps_cap_t ps_cap;
    memset(&ps_cap, 0, sizeof(sgx_ps_cap_t));
    sgx_ret = sgx_get_ps_cap(&ps_cap);
    if (sgx_ret)
    {
        cerr << "cannot get platform service capability, error code = 0x" << hex << sgx_ret <<endl;
        return sgx_ret;
    }

    if (!SGX_IS_MONOTONIC_COUNTER_AVAILABLE(ps_cap))
    {
        cerr << "monotonic counter is not supported" << endl;
        return SGX_ERROR_SERVICE_UNAVAILABLE;
    }

    uint32_t enclave_ret = 0;
    sgx_ret = create_sealed_policy(enclave_id, &enclave_ret, (uint8_t *)stored_sealed_activity_log, sealed_activity_log_length);
    if (sgx_ret)
    {
        cerr << "call create_sealed_policy fail, error code = 0x" << hex << sgx_ret << endl;
        return sgx_ret;
    }

    if (enclave_ret)
    {
        cerr << "cannot create_sealed_policy, function return fail, error code =" << "0x" << hex<< enclave_ret << endl;
        return enclave_ret;
    }
    return 0;
}

uint32_t ReplayProtectedDRM:: Init()
{
    return Init(sealed_activity_log);
}

uint32_t ReplayProtectedDRM:: PerformFunction(uint8_t* stored_sealed_activity_log)
{
    sgx_status_t sgx_ret = SGX_ERROR_UNEXPECTED;
    uint32_t enclave_ret = 0;

    sgx_ret = perform_sealed_policy(enclave_id, &enclave_ret, (uint8_t *)stored_sealed_activity_log, sealed_activity_log_length);
    if (sgx_ret)
    {
        cerr << "call perform_sealed_policy fail, error code = 0x" << hex << sgx_ret << endl;
        return sgx_ret;
    }

    if (enclave_ret)
    {
        cerr << "cannot perform_sealed_policy, function return fail, error code =" << "0x"<< hex << enclave_ret << endl;
        return enclave_ret;
    }

    return 0;
}

uint32_t ReplayProtectedDRM:: PerformFunction()
{
    return PerformFunction(sealed_activity_log);
}

uint32_t ReplayProtectedDRM::UpdateSecret(uint8_t* stored_sealed_activity_log)
{
    sgx_status_t sgx_ret = SGX_ERROR_UNEXPECTED;
    uint32_t enclave_ret = 0;

    sgx_ret = update_sealed_policy(enclave_id, &enclave_ret, (uint8_t *)stored_sealed_activity_log, sealed_activity_log_length);
    if (sgx_ret)
    {
        cerr << "call update_sealed_policy fail, error code = 0x" << hex << sgx_ret << endl;
        return sgx_ret;
    }

    if (enclave_ret)
    {
        cerr << "cannot update_sealed_policy, function return fail, error code =" << "0x" << hex << enclave_ret << endl;
        return enclave_ret;
    }

    return 0;
}

uint32_t ReplayProtectedDRM::UpdateSecret()
{
    return UpdateSecret(sealed_activity_log);
}

uint32_t ReplayProtectedDRM::DeleteSecret(uint8_t* stored_sealed_activity_log)
{
    sgx_status_t sgx_ret = SGX_ERROR_UNEXPECTED;
    uint32_t enclave_ret = 0;

    sgx_ret = delete_sealed_policy(enclave_id, &enclave_ret, (uint8_t *)stored_sealed_activity_log, sealed_activity_log_length);
    if (sgx_ret)
    {
        cerr << "call delete_sealed_policy fail, error code = 0x" << hex<< sgx_ret << endl;
        return sgx_ret;
    }

    if (enclave_ret)
    {
        cerr << "cannot delete_sealed_policy, function return fail, error code =" << "0x" << hex << enclave_ret << endl;
        return enclave_ret;
    }
    return 0;
}

uint32_t ReplayProtectedDRM::DeleteSecret()
{
    return DeleteSecret(sealed_activity_log);
}
