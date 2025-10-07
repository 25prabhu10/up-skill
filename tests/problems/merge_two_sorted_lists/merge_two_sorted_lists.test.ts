import { describe, expect, test } from 'vitest'

import { mergeTwoLists } from '@/problems/merge_two_sorted_lists/merge_two_sorted_lists'
import { ListNode } from '@/utils/classes'

describe('Merge Two Sorted Lists', () => {
  test('should return the head of the merged linked list', () => {
    const head1 = new ListNode(1, new ListNode(2, new ListNode(4)))
    const head2 = new ListNode(1, new ListNode(3, new ListNode(4)))

    let resultHead = mergeTwoLists(head1, head2)
    const result: number[] = []
    const expected = [1, 1, 2, 3, 4, 4]

    while (resultHead) {
      result.push(resultHead.val)
      resultHead = resultHead.next
    }
    expect(result).toStrictEqual(expected)
  })

  test('should return the head of the merged linked list with single items', () => {
    const head1 = new ListNode(1, null)
    const head2 = new ListNode(2, null)

    let resultHead = mergeTwoLists(head1, head2)
    const result: number[] = []
    const expected = [1, 2]

    while (resultHead) {
      result.push(resultHead.val)
      resultHead = resultHead.next
    }
    expect(result).toStrictEqual(expected)
  })

  test('should return the head of the merged empty linked list', () => {
    let resultHead = mergeTwoLists(null, null)
    const result: number[] = []
    const expected: number[] = []

    while (resultHead) {
      result.push(resultHead.val)
      resultHead = resultHead.next
    }
    expect(result).toStrictEqual(expected)
  })

  test('should return the head of the merged empty linked list and non-empty linked list', () => {
    let resultHead = mergeTwoLists(null, new ListNode(0))
    const result: number[] = []
    const expected: number[] = [0]

    while (resultHead) {
      result.push(resultHead.val)
      resultHead = resultHead.next
    }
    expect(result).toStrictEqual(expected)
  })
})
