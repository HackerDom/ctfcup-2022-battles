#include <memory>
#include <vector>
#include <stdint.h>
#include "signal.hpp"
#include "data.hpp"

class Rom {
    std::vector<uint8_t> buf;
    int pos = 0;

    int get_bit(std::uint8_t b, int index) {
        if ((b & (1 << index)) > 0) {
            return 1;
        }

        return 0;
    }
public:
    std::shared_ptr<Signal> opcode_signal;
    std::shared_ptr<Signal> reg_signal1;
    std::shared_ptr<Signal> reg_signal2;
    std::shared_ptr<Data> opcode_output;

    Rom(
        std::vector<uint8_t> buf,
        std::shared_ptr<Signal> opcode_signal,
        std::shared_ptr<Signal> reg_signal1,
        std::shared_ptr<Signal> reg_signal2,
        std::shared_ptr<Data> opcode_output
    );

    int next();
};
