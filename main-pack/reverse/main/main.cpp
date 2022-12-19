#include <iostream>
#include <memory>
#include <cmath>
#include <fstream>
#include "gates/half_adder.hpp"
#include "gates/signal.hpp"
#include "gates/decoder.hpp"
#include "gates/bitwise_adder.hpp"
#include "gates/register.hpp"
#include "gates/data.hpp"
#include "gates/rom.hpp"
#include "gates/ram.hpp"

using namespace std;

vector<uint8_t> read_from(const char* path) {
    ifstream input(path);
    if (input.fail()) {
        cout << "read file \"" << path << "\" failed" << endl;
        exit(1);
    }

    auto input_buf = vector<uint8_t>(istreambuf_iterator<char>(input), {});
    if (input_buf.size() > 100) {
        cout << "Wrong answer!" << endl;
        exit(1);
    }
 
    return input_buf; 
}

int main() {
    auto opcode_signal = make_shared<Signal>(0);
    auto reg_signal1 = make_shared<Signal>(0);
    auto reg_signal2 = make_shared<Signal>(0);
    auto opcode_output = make_shared<Data>(vector<int>{0, 0, 0, 0, 0, 0, 0, 0});
    
    auto rom = Rom(
        read_from("./data/rom"), opcode_signal, reg_signal1, reg_signal2, opcode_output
    );

    auto adder_output = make_shared<Data>(vector<int>{0, 0, 0, 0, 0, 0, 0, 0});

    auto a_inputs = vector<shared_ptr<Gate<int>>>(); 
    a_inputs.reserve(8);
    auto b_inputs = vector<shared_ptr<Gate<int>>>(); 
    b_inputs.reserve(8);

    auto zero = make_shared<Signal>(0);
    auto reg_decoder = Decoder(reg_signal1, reg_signal2);
    
    auto reg_ands = vector<shared_ptr<Nand>>();
    for (int i = 0; i < 3; i++) {
        auto u = make_shared<Nand>();
        u->map_in(opcode_signal, reg_decoder.outs[i]);
        auto not_u = make_shared<Nand>();
        not_u->map_in(u, u);
        reg_ands.push_back(not_u);
    }

    auto not_opcode_signal = make_shared<Nand>();
    not_opcode_signal->map_in(opcode_signal, opcode_signal);
    auto addr = make_shared<Register>(
        vector<int>{0, 0, 0, 0, 0, 0, 0, 0}, opcode_output, not_opcode_signal
    );
    auto ram = make_shared<Ram>(
        read_from("./data/ram"), adder_output, addr, reg_ands[2]
    );
    auto a = Register(
        vector<int>{0, 0, 0, 0, 0, 0, 0, 0}, ram, reg_ands[0]
    );
    auto b = Register(
        vector<int>{0, 0, 0, 0, 0, 0, 0, 0}, ram, reg_ands[1]
    );

    while (rom.next() + 1) {
        auto a_out = a.out();
        for (auto i = 0ull; i < a.out().size(); i++) {
            a_inputs.push_back(make_shared<Signal>(a_out[i]));
        }
        auto b_out = b.out();
        for (auto i = 0ull; i < b.out().size(); i++) {
            b_inputs.push_back(make_shared<Signal>(b_out[i]));
        }
        auto adder = BitwiseAdder(a_inputs, b_inputs, zero);
        auto adder_out = vector<int>();
        for (auto n : adder.sum()) {
            adder_out.push_back(n->out());
        }
        a_inputs.clear();
        b_inputs.clear();
        adder_output->set(adder_out);

        addr->out();
        ram->load();
    }

    auto ram_buf = ram->internal_buf();
    int j = 0;
    for (int i = 0; i < 10; i += 2) {
        if (ram_buf[i] + ram_buf[i + 1] != ram_buf[j + 10]) {
            cout << "Wrong answer!" << endl;
            return -1;
        }
        j++;
    }

    cout << "Congratulations! Your flag: " << getenv("FLAG") << endl;
    return 0; 
}