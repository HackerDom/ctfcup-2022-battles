#include <memory>
#include <vector>
#include "gate.hpp"
#include "bitwise_adder.hpp"
#include "adder.hpp"

BitwiseAdder::BitwiseAdder(
        std::vector<std::shared_ptr<Gate<int>>> inputs1, 
        std::vector<std::shared_ptr<Gate<int>>> inputs2,
        std::shared_ptr<Gate<int>> carry_in
) {
    if (inputs1.size() != inputs2.size()) {
        return;
    }
 
    outputs = std::vector<Adder>();
    std::shared_ptr<Gate<int>> c_in = carry_in;
    for (auto i = 0; i < inputs1.size(); i++) {
        auto a = Adder();
        a.map_in(inputs1[i], inputs2[i], c_in);
        outputs.push_back(a);
        c_in = a.carry_out();
    }
}

std::vector<std::shared_ptr<Nand>> BitwiseAdder::sum() {
    auto res = std::vector<std::shared_ptr<Nand>>();
    for (auto i : outputs) {
        res.push_back(i.sum());
    }

    return res;
}

std::shared_ptr<Nand> BitwiseAdder::carry_out() {
    return outputs[outputs.size() - 1].carry_out();
}
