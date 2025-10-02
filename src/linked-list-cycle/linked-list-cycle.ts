/**
 * Definition for singly-linked list.
 * class ListNode {
 *     val: number
 *     next: ListNode | null
 *     constructor(val?: number, next?: ListNode | null) {
 *         this.val = (val===undefined ? 0 : val)
 *         this.next = (next===undefined ? null : next)
 *     }
 * }
 */
export class ListNode {
  val: number
  next: ListNode | null
  constructor(val?: number, next?: ListNode | null) {
    this.val = val === undefined ? 0 : val
    this.next = next === undefined ? null : next
  }
}

export function hasCycle(head: ListNode | null): boolean {
  let node = head
  const items = new Set<ListNode>()

  while (node?.next) {
    if (items.has(node.next)) {
      return true
    }
    items.add(node)
    node = node.next
  }

  return false
}

export function hasCycleFastest(head: ListNode | null): boolean {
  if (!head || !head.next) {
    return false
  } // no nodes or single node → no cycle

  let slow: ListNode | null | undefined = head
  let fast: ListNode | null | undefined = head.next

  while (slow !== fast) {
    if (!fast || !fast.next) {
      return false // reached end → no cycle
    }
    slow = slow?.next
    fast = fast.next.next
  }

  return true // slow and fast met → cycle exists
}

export function hasCycleFastestJS(head: ListNode | null): boolean {
  let slow: ListNode | null | undefined = head
  let fast: ListNode | null | undefined = head
  while (fast && fast.next) {
    slow = slow?.next
    fast = fast.next.next
    if (slow === fast) {
      return true
    }
  }
  return false
}
