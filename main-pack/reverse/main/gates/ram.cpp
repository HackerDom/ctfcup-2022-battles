#include "ram.hpp"
#include "data.hpp"

Ram::Ram(
    std::vector<uint8_t> buf,
    std::shared_ptr<Data> data,
    std::shared_ptr<Gate<std::vector<int>>> addr,
    std::shared_ptr<Gate<int>> control
):buf(buf),data(data),addr(addr),control(control){}

void Ram::load() {
    if (control->out()) {
        int address = line_to_number(addr->out());
        buf[address] = line_to_number(data->out());
    }
}

std::vector<int> Ram::out() {
    auto address = line_to_number(addr->out());
    return number_to_line(buf[address]);
}

int Ram::line_to_number(std::vector<int> line) {
    int n = 0;
    int p = 1;
    for (auto i : line) {
        n += i * p;
        p *= 2;
    }

    return n;
}

std::vector<int> Ram::number_to_line(int number) {
    auto line = std::vector<int>();

    for (auto i = 0; i < 8; i++) {
        int bit = (number & (1 << i)) >> i;
        line.push_back(bit);
    }

    return line;
}

std::vector<uint8_t> Ram::internal_buf() {
    return buf;
}
