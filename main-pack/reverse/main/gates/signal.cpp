#include "signal.hpp"

Signal::Signal(int signal):signal(signal) {}

int Signal::out() {
    return signal;
}

void Signal::set(int signal) {
    this->signal = signal;
}
