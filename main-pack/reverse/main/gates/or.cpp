#include <memory>
#include "or.hpp"

void Or::map_in1(std::shared_ptr<Gate<int>> in) {
    if (frozen) {
        return;
    }

    in1 = in;
    rebuild();
}

void Or::map_in2(std::shared_ptr<Gate<int>> in) {
    if (frozen) {
        return;
    }

    in2 = in;
    rebuild();
}

void Or::map_in(std::shared_ptr<Gate<int>> in1, std::shared_ptr<Gate<int>> in2) {
    if (frozen) {
        return;
    }

    this->in1 = in1;
    this->in2 = in2;

    rebuild();
}

void Or::rebuild() {
    auto u1 = std::make_shared<Nand>(Nand());
    auto u2 = std::make_shared<Nand>(Nand());

    u1->map_in(in1, in1);
    u2->map_in(in2, in2);

    result = std::make_shared<Nand>(Nand());
    result->map_in(u1, u2);
    result->freeze();
}

void Or::freeze() {
    frozen = true;
}
