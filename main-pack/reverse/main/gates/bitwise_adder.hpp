#include <memory>
#include <vector>
#include "gate.hpp"
#include "adder.hpp"

class BitwiseAdder {
    std::vector<Adder> outputs;

public:
    BitwiseAdder(
        std::vector<std::shared_ptr<Gate<int>>> inputs1, 
        std::vector<std::shared_ptr<Gate<int>>> inputs2,
        std::shared_ptr<Gate<int>> carry_in
    );

    std::vector<std::shared_ptr<Nand>> sum();
    std::shared_ptr<Nand> carry_out();
};