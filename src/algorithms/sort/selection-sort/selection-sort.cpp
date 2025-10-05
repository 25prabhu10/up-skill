#include <cstdlib>
#include <iostream>

#define MAX 100

using namespace std;

void print_array(int arr[], int size) {
  for (int i = 0; i < size; i++) {
    cout << arr[i] << " ";
  }
  cout << endl;
}

void swap(int arr[], int x, int y) {
  int temp = arr[x];
  arr[x] = arr[y];
  arr[y] = temp;
}

int find_smallest(int arr[], int start, int end) {
  int index_of_smallest = start;

  for (int j = start; j <= end; j++) {
    if (arr[index_of_smallest] > arr[j]) {
      index_of_smallest = j;
    }
  }

  return index_of_smallest;
}

void selection_sort(int arr[], int n) {
  for (int i = 0; i < n; i++) {
    int swap_index = find_smallest(arr, i, n - 1);
    swap(arr, i, swap_index);
    print_array(arr, n);
  }
}

int main() {
  int arr[MAX] = {5, 2, 9, 1, 5};
  // int n = sizeof(arr) / sizeof(int);
  int n = 5;

  // cout << "Enter number of elements in the array:" << endl;
  // cin >> n;
  //
  // for (int i = 0; i < n; i++) {
  //   arr[i] = rand();
  // }

  print_array(arr, n);

  selection_sort(arr, n);

  print_array(arr, n);

  return 0;
}
