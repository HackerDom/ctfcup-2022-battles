#pragma once

template <typename GateOut>
class Gate {
public:
    virtual GateOut out() = 0;
};