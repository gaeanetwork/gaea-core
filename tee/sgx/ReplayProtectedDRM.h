#include "stdlib.h"
#include "sgx.h"
#include "sgx_urts.h"
#include "sgx_sealed_data.h"


class ReplayProtectedDRM
{
public:
    ReplayProtectedDRM();
    ~ReplayProtectedDRM(void);
    
    uint32_t Init(uint8_t*  stored_sealed_activity_log);
    uint32_t Init();
    uint32_t PerformFunction();
    uint32_t PerformFunction(uint8_t* stored_sealed_activity_log);
    uint32_t UpdateSecret();
    uint32_t UpdateSecret(uint8_t* stored_sealed_activity_log);

    uint32_t DeleteSecret();
    uint32_t DeleteSecret(uint8_t* stored_sealed_activity_log);

    uint32_t GetActivityLog(uint8_t* stored_sealed_activity_log);

    static const uint32_t sealed_activity_log_length = SEALED_REPLAY_PROTECTED_PAY_LOAD_SIZE;

private:
    uint8_t  sealed_activity_log[sealed_activity_log_length];
    sgx_enclave_id_t enclave_id;
};
