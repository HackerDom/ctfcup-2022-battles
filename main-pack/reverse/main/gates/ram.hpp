#include <vector>
#include <stdint.h>
#include <memory>
#include "gate.hpp"
#include "data.hpp"

class Ram : public Gate<std::vector<int>> {
    std::vector<uint8_t> buf;
    std::shared_ptr<Data> data;
    std::shared_ptr<Gate<std::vector<int>>> addr;
    std::shared_ptr<Gate<int>> control;

    int line_to_number(std::vector<int> line);
    std::vector<int> number_to_line(int number);
public:
    Ram(
        std::vector<uint8_t> buf,
        std::shared_ptr<Data> data,
        std::shared_ptr<Gate<std::vector<int>>> addr,
        std::shared_ptr<Gate<int>> control
    );

    void load();
    std::vector<uint8_t> internal_buf();
    virtual std::vector<int> out();
};
