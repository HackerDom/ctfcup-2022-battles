#include <vector>
#include <iostream>


struct User {
    uint64_t id;
    std::string name;
};

std::vector<User> users;


void add_user() {
    User user;

    user.id = users.size() + 1;

    std::cout << "[?] Input name: ";
    std::cin >> user.name;

    users.push_back(user);

    std::cout << "[+] User added" << std::endl;

    return;
}

void delete_user() {
    size_t index;

    std::cout << "[?] Input index: ";
    std::cin >> index;

    users.erase(users.begin() + index);

    std::cout << "[+] User deleted" << std::endl;

    return;
}

void print_users() {
    std::cout << "[*] Users:" << std::endl;

    for (auto it = users.begin(); it != users.end(); it++) {
        std::cout << it->name << " (" << it->id << ") " << std::endl;
    }

    return;
}

void menu() {
    std::cout << std::endl;
    std::cout << "1. Add user" << std::endl;
    std::cout << "2. Delete user" << std::endl;
    std::cout << "3. Print users" << std::endl;
    std::cout << "4. Exit" << std::endl;
    std::cout << "> ";

    int choice;
    std::cin >> choice;

    switch (choice) {
        case 1:
            return add_user();

        case 2:
            return delete_user();

        case 3:
            return print_users();

        case 4:
            std::exit(0);

        default:
            std::cout << "[-] Invalid option" << std::endl;
            return;
    }
}

int main(int argc, char **argv) {
    std::cout << "[!] Hello" << std::endl;

    while (true) {
        menu();
    }
}
