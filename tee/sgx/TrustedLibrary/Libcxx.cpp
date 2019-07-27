#include <cstdlib>
#include <string>
#include <map>
#include <algorithm>

#include "../Enclave.h"
#include "Enclave_t.h"

using namespace std;

/*
 * ecall_exception:
 *   throw/catch C++ exception inside the enclave.
 */

void ecall_exception(void)
{
    std::string foo = "foo";
    try {
        throw std::runtime_error(foo);
    }
    catch (std::runtime_error const& e) {
        assert( foo == e.what() );
        std::runtime_error clone("");
        clone = e;
        assert(foo == clone.what() );
    }
    catch (...) {
        assert( false );
    }
}


/*
 * ecall_map:
 *   Utilize STL <map> in the enclave.
 */
void ecall_map(void)
{
    typedef map<char, int, less<char> > map_t;
    typedef map_t::value_type map_value;
    map_t m;

    m.insert(map_value('a', 1));
    m.insert(map_value('b', 2));
    m.insert(map_value('c', 3));
    m.insert(map_value('d', 4));

    assert(m['a'] == 1);
    assert(m['b'] == 2);
    assert(m['c'] == 3);
    assert(m['d'] == 4);

    assert(m.find('e') == m.end());

    return;
}
