#include <iostream>
#include <fstream>
#include <vector>
using namespace std;

int main() {
    auto paths = {
        "e1ce9d74-b40a-4000-88d6-5664aee31f5a/ca6f8ede-c6e5-4b8e-ac5f-fea0ec4afc65/f39f15ff-b05a-45f4-a904-bf577aca92a9",
        "1d82ad4a-e8b4-4842-825a-c01a73402cf3/3d2fcf29-deae-4440-8c7d-34de29bd212e/b50154ad-7a0e-478b-bca3-3a867c22cf8a",
        "35eb15ec-7d75-4daf-8d20-aadcd2cd125a",
        "e1ce9d74-b40a-4000-88d6-5664aee31f5a/ca6f8ede-c6e5-4b8e-ac5f-fea0ec4afc65/b254b8b0-156d-4554-930b-7cb11bd5a0d4"
    };

    auto bufs = vector<vector<uint8_t>>();
    for (auto &p : paths) {
        ifstream input(p, ios::binary);

        if (input.fail()) {
            cout << "open file failed, run program from directory where is it located" << endl;
            return 1;
        }

        vector<uint8_t> buf(istreambuf_iterator<char>(input), {});
        bufs.push_back(move(buf));
    }

    auto p = vector<uint8_t>();
    for (auto &b : bufs) {
        p.reserve(p.size() + distance(b.begin(), b.end()));
        p.insert(p.end(), b.begin(), b.end());
    }

    for (int i = 0; i < 37; i++) {
        cout << char(p[i] ^ p[i + 37]); 
    }

    cout << endl;

    return 0;
}