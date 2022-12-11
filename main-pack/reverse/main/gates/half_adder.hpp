#pragma once
#include <memory>
#include "gate.hpp"
#include "nand.hpp"

class HalfAdder {
    std::shared_ptr<Gate<int>> in1;
    std::shared_ptr<Gate<int>> in2;
    bool frozen = false;

    void rebuild();
public:
    std::shared_ptr<Nand> sum;
    std::shared_ptr<Nand> carry_out;
    
    void map_in1(std::shared_ptr<Gate<int>> in);
    void map_in2(std::shared_ptr<Gate<int>> in);
    void map_in(std::shared_ptr<Gate<int>> in1, std::shared_ptr<Gate<int>> in2);
    void freeze();
};