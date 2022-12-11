#include "register.hpp"

Register::Register(
    std::vector<int> state, 
    std::shared_ptr<Gate<std::vector<int>>> in, 
    std::shared_ptr<Gate<int>> control
):state(state),in(in),control(control){}

std::vector<int> Register::out() {
    if (!control->out()) {
        return state;
    }

    state = in->out();
    return state;
}
