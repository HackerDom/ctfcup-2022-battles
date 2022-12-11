#pragma once
#include <memory>
#include "gate.hpp"
#include "half_adder.hpp"
#include "or.hpp"

class Adder {
    std::shared_ptr<Gate<int>> in1;
    std::shared_ptr<Gate<int>> in2;
    std::shared_ptr<Gate<int>> carry_in;
    bool frozen = false;
    HalfAdder _sum;
    Or _carry_out;

    void rebuild();
public:
    std::shared_ptr<Nand> sum();
    std::shared_ptr<Nand> carry_out();
    void map_in1(std::shared_ptr<Gate<int>> in);
    void map_in2(std::shared_ptr<Gate<int>> in);
    void map_carry_in(std::shared_ptr<Gate<int>> in);
    void map_in(
        std::shared_ptr<Gate<int>> in1, 
        std::shared_ptr<Gate<int>> in2, 
        std::shared_ptr<Gate<int>> carry_in
    );
    void freeze();
};