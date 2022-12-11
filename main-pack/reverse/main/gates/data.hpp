#pragma once
#include <vector>
#include "gate.hpp"

class Data : public Gate<std::vector<int>> {
    std::vector<int> data;
public:
    Data(std::vector<int> data);

    void set(std::vector<int> data);
    virtual std::vector<int> out();
};