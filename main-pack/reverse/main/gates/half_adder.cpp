#include <memory>
#include "half_adder.hpp"

void HalfAdder::map_in1(std::shared_ptr<Gate<int>> in) {
    if (frozen) {
        return;
    }

    in1 = in;
    rebuild();
}

void HalfAdder::map_in2(std::shared_ptr<Gate<int>> in) {
    if (frozen) {
        return;
    }

    in2 = in;
    rebuild();
}

void HalfAdder::map_in(
    std::shared_ptr<Gate<int>> in1, 
    std::shared_ptr<Gate<int>> in2
) {
    if (frozen) {
        return;
    }

    this->in1 = in1;
    this->in2 = in2;
    rebuild();
}

void HalfAdder::rebuild() {
    auto u1 = std::make_shared<Nand>(Nand());
    auto u2 = std::make_shared<Nand>(Nand());
    auto u3 = std::make_shared<Nand>(Nand());

    u1->map_in1(in1);
    u1->map_in2(in2);
    u2->map_in1(u1);
    u2->map_in2(in1);
    u3->map_in1(u1);
    u3->map_in2(in2);

    sum = std::make_shared<Nand>(Nand());
    sum->map_in1(u2);
    sum->map_in2(u3);
    sum->freeze();

    carry_out = std::make_shared<Nand>(Nand());
    carry_out->map_in1(u1);
    carry_out->map_in2(u1);
    carry_out->freeze();
}

void HalfAdder::freeze() {
    frozen = true;
}
