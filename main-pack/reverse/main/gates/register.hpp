#include <vector>
#include <memory>
#include "gate.hpp"

class Register : public Gate<std::vector<int>>
{
    std::vector<int> state;
    std::shared_ptr<Gate<std::vector<int>>> in;
    std::shared_ptr<Gate<int>> control;

public:
    Register(
        std::vector<int> state, 
        std::shared_ptr<Gate<std::vector<int>>> in, 
        std::shared_ptr<Gate<int>> control
    );

    virtual std::vector<int> out();
};
