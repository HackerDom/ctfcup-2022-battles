#include "nand.hpp"
#include "gate.hpp"
#include <memory>
#include <iostream>

int Nand::out() {
    return (in1->out() & in2->out()) ^ 1;
}

void Nand::freeze() {
    frozen = true;
}

void Nand::map_in1(std::shared_ptr<Gate<int>> in) {
    if (frozen) {
        return;
    }

    in1 = in;
}

void Nand::map_in2(std::shared_ptr<Gate<int>> in) {
    if (frozen) {
        return;
    }

    in2 = in;
}

void Nand::map_in(std::shared_ptr<Gate<int>> in1, std::shared_ptr<Gate<int>> in2) {
    if (frozen) {
        return;
    }

    this->in1 = in1;
    this->in2 = in2;
}