import { describe, expect, test } from 'vitest'

import { hasCycle, ListNode } from '@/linked-list-cycle/linked-list-cycle'

describe('Linked List Cycle', () => {
  test('should return true if list is cyclic', () => {
    const head1 = new ListNode(3)
    const node12 = new ListNode(2)
    const node13 = new ListNode(0)
    const node14 = new ListNode(-4)

    head1.next = node12
    node12.next = node13
    node13.next = node14
    node14.next = node12
    expect(hasCycle(head1)).toBeTruthy()

    const head2 = new ListNode(1)
    const node22 = new ListNode(2)

    head2.next = node22
    node22.next = head2
    expect(hasCycle(head2)).toBeTruthy()

    const head3 = new ListNode(1)
    expect(hasCycle(head3)).toBeFalsy()
  })
})
