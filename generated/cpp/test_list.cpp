// 生成コードのコンパイル検証用
#include "list.h"
#include <iostream>
#include <string>

int main() {
    // int型リストのテスト
    List<int> intList;
    intList.insert(3);
    intList.insert(2);
    intList.insert(1);

    std::cout << "length: " << intList.length() << std::endl;

    std::cout << "traverse: ";
    intList.traverse([](const int& v) { std::cout << v << " "; });
    std::cout << std::endl;

    auto* found = intList.find(2);
    std::cout << "find(2): " << (found ? "found" : "not found") << std::endl;

    intList.remove();
    std::cout << "after remove, length: " << intList.length() << std::endl;

    // string型リストのテスト
    List<std::string> strList;
    strList.insert("world");
    strList.insert("hello");
    strList.traverse([](const std::string& v) { std::cout << v << " "; });
    std::cout << std::endl;

    return 0;
}
