#include <memory>
#include <vector>
#include "gate.hpp"
#include "nand.hpp"

class Decoder {
    std::shared_ptr<Gate<int>> in1;
    std::shared_ptr<Gate<int>> in2;

public:
    std::vector<std::shared_ptr<Nand>> outs;

    Decoder(std::shared_ptr<Gate<int>> in1, std::shared_ptr<Gate<int>> in2);
};
