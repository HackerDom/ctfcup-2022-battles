#pragma once
#include "gate.hpp"
#include <memory>
#include "signal.hpp"

class Nand : public Gate<int> {
    std::shared_ptr<Gate<int>> in1 = std::make_shared<Signal>(Signal(0));
    std::shared_ptr<Gate<int>> in2 = std::make_shared<Signal>(Signal(0));
    bool frozen = false; // remap inputs isn't allowed

public:
    virtual int out();
    void freeze();
    void map_in1(std::shared_ptr<Gate<int>> in);
    void map_in2(std::shared_ptr<Gate<int>> in);
    void map_in(std::shared_ptr<Gate<int>> in1, std::shared_ptr<Gate<int>> in2);
};