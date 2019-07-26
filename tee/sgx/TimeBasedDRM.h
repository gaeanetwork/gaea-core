#include "stdlib.h"
#include "sgx.h"
#include "sgx_urts.h"
#include "sgx_sealed_data.h"

class TimeBasedDRM
{
public:
    TimeBasedDRM(void);
    ~TimeBasedDRM(void);

    uint32_t Init(uint8_t*  stored_time_based_policy);
    uint32_t Init();
    
    uint32_t PerformFunction();
    uint32_t PerformFunction(uint8_t* stored_time_based_policy);
    uint32_t GetTimeBasedPolicy(uint8_t* stored_time_based_policy);

    static const uint32_t time_based_policy_length = TIME_BASED_PAY_LOAD_SIZE;

private:
    uint8_t time_based_policy[time_based_policy_length];
    sgx_enclave_id_t enclave_id;
};
