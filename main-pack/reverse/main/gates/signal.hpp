#pragma once
#include "gate.hpp"

class Signal : public Gate<int> {
    int signal;

public:
    Signal(int signal);

    void set(int signal);
    virtual int out();
};