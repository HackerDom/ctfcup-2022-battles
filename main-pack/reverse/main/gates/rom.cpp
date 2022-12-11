#include "rom.hpp"

Rom::Rom(
        std::vector<uint8_t> buf,
        std::shared_ptr<Signal> opcode_signal,
        std::shared_ptr<Signal> reg_signal1,
        std::shared_ptr<Signal> reg_signal2,
        std::shared_ptr<Data> opcode_output
):buf(buf),opcode_signal(opcode_signal),reg_signal1(reg_signal1),reg_signal2(reg_signal2),
opcode_output(opcode_output){}

int Rom::next() {
    if ((long unsigned int)pos > buf.size() - 1) {
        return -1;
    }

    auto control = buf[pos];
    opcode_signal->set(get_bit(control, 0));
    reg_signal1->set(get_bit(control, 1));
    reg_signal2->set(get_bit(control, 2));

    auto out = buf[++pos];
    auto output = std::vector<int>(); 
    for (auto i = 0; i < 8; i++) {
        int bit = (out & (1 << i)) >> i;
        output.push_back(bit);
    }
    opcode_output->set(output);

    pos++;
    return 0;
}
