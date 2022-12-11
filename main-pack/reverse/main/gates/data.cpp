#include <vector>
#include "data.hpp"

Data::Data(std::vector<int> data):data(data){}

std::vector<int> Data::out() {
    return data;
}

void Data::set(std::vector<int> data) {
    this->data = data;
}
