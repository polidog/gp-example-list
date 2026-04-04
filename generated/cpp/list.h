// Generated from: docs/gp/08-dsl-examples/example-01-minimal.yaml
// Configuration: Copy + Monomorphic, No LengthCounter, No Tracing
// Language: C++

#ifndef LIST_H
#define LIST_H

#include <cstddef>
#include <functional>

// R2: Node構造体 — CopyMonomorphicNode<T>
// Ownership=Copy: 値埋め込み (T data)
// Morphology=Monomorphic: 型Tのみ許容
template<typename T>
struct Node {
    T data;
    Node<T>* next;
};

// R1: List構造体
// LengthCounter=無効: lengthフィールドなし
// Tracing=無効: TraceLogなし
template<typename T>
class List {
public:
    List();
    ~List();

    // R3: insert — 先頭への挿入
    // Ownership=Copy: copy(element)で値を埋め込み
    void insert(const T& element);

    // R4: remove — 先頭要素の削除
    // Ownership=Copy: ノード解放のみ（値は自動破棄）
    bool remove();

    // R6: length — O(n)走査（LengthCounter無効）
    std::size_t length() const;

    // find — 走査ベースの検索
    Node<T>* find(const T& element) const;

    // traverse — イテレーション
    void traverse(std::function<void(const T&)> fn) const;

    // isEmpty
    bool isEmpty() const;

private:
    Node<T>* head_;
};

// --- テンプレート実装 ---

template<typename T>
List<T>::List() : head_(nullptr) {}

// R5: destroy — Ownership=Copy: 各ノードをdelete（値は自動破棄）
template<typename T>
List<T>::~List() {
    Node<T>* current = head_;
    while (current != nullptr) {
        Node<T>* next = current->next;
        delete current;
        current = next;
    }
}

// R3: insert
template<typename T>
void List<T>::insert(const T& element) {
    Node<T>* node = new Node<T>{element, head_};
    head_ = node;
}

// R4: remove
template<typename T>
bool List<T>::remove() {
    if (head_ == nullptr) {
        return false;
    }
    Node<T>* old_head = head_;
    head_ = head_->next;
    delete old_head;
    return true;
}

// R6: length — LengthCounter無効のためO(n)走査
template<typename T>
std::size_t List<T>::length() const {
    std::size_t count = 0;
    Node<T>* current = head_;
    while (current != nullptr) {
        count++;
        current = current->next;
    }
    return count;
}

template<typename T>
Node<T>* List<T>::find(const T& element) const {
    Node<T>* current = head_;
    while (current != nullptr) {
        if (current->data == element) {
            return current;
        }
        current = current->next;
    }
    return nullptr;
}

template<typename T>
void List<T>::traverse(std::function<void(const T&)> fn) const {
    Node<T>* current = head_;
    while (current != nullptr) {
        fn(current->data);
        current = current->next;
    }
}

template<typename T>
bool List<T>::isEmpty() const {
    return head_ == nullptr;
}

#endif // LIST_H
