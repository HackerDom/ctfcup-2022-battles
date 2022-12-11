#include <memory>
#include <vector>
#include "decoder.hpp"

Decoder::Decoder(std::shared_ptr<Gate<int>> in1, std::shared_ptr<Gate<int>> in2) {
    auto not_in1 = std::make_shared<Nand>();
    not_in1->map_in(in1, in1);
    auto not_in2 = std::make_shared<Nand>();
    not_in2->map_in(in2, in2);

    auto u1 = std::make_shared<Nand>(); 
    u1->map_in(not_in1, in2);
    auto u2 = std::make_shared<Nand>(); 
    u2->map_in(in1, in2);
    auto u3 = std::make_shared<Nand>(); 
    u3->map_in(not_in1, not_in2);
    auto u4 = std::make_shared<Nand>(); 
    u4->map_in(in1, not_in2);

    for (auto u : std::vector<std::shared_ptr<Nand>>{u1, u2, u3, u4}) {
        auto out = std::make_shared<Nand>();
        out->map_in(u, u);
        out->freeze();
        outs.push_back(out);
    }
}
