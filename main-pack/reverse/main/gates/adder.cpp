#include <memory>
#include "gate.hpp"
#include "adder.hpp"

void Adder::map_in1(std::shared_ptr<Gate<int>> in) {
    if (frozen) {
        return;
    }

    in1 = in;
    rebuild();
}

void Adder::map_in2(std::shared_ptr<Gate<int>> in) {
    if (frozen) {
        return;
    }

    in2 = in;
    rebuild();
}

void Adder::map_carry_in(std::shared_ptr<Gate<int>> in) {
    if (frozen) {
        return;
    }

    carry_in = in;
    rebuild();
}

void Adder::map_in(
    std::shared_ptr<Gate<int>> in1, 
    std::shared_ptr<Gate<int>> in2,
    std::shared_ptr<Gate<int>> carry_in
) {
    if (frozen) {
        return;
    }

    this->in1 = in1;
    this->in2 = in2;
    this->carry_in = carry_in;
    rebuild();
}

std::shared_ptr<Nand> Adder::sum() {
    return _sum.sum;
}

std::shared_ptr<Nand> Adder::carry_out() {
    return _carry_out.result;
}

void Adder::rebuild() {
    auto half_adder = std::make_shared<HalfAdder>(HalfAdder());
    half_adder->map_in(in1, in2);
    
    _sum = HalfAdder();
    _sum.map_in(half_adder->sum, carry_in);
    _sum.freeze();

    _carry_out = Or();
    _carry_out.map_in(_sum.carry_out, half_adder->carry_out);
    _carry_out.freeze();
}

void Adder::freeze() {
    frozen = true;
}
