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
import { ListNode } from "@/utils/classes"

export function mergeTwoLists(list1: ListNode | null, list2: ListNode | null): ListNode | null {
  if (list1 === null) {
    return list2
  } else if (list2 === null) {
    return list1
  }

  let currentHead: ListNode | null = null

  if (list1.val < list2.val) {
    currentHead = new ListNode(list1.val, null)
    list1 = list1.next
  } else {
    currentHead = new ListNode(list2.val, null)
    list2 = list2.next
  }

  let head: ListNode | null = currentHead

  while (list1 && list2) {
    if (list1.val < list2.val) {
      currentHead.next = new ListNode(list1.val, null)
      list1 = list1.next
    } else {
      currentHead.next = new ListNode(list2.val, null)
      list2 = list2.next
    }

    currentHead = currentHead.next
  }

  if (list1) {
    currentHead.next = list1
  } else if (list2) {
    currentHead.next = list2
  }

  return head
}

export function mergeTwoListsFastestJS(
  list1: ListNode | null,
  list2: ListNode | null
): ListNode | null {
  let resultSt = new ListNode()
  let curr = resultSt
  while (list1 && list2) {
    if (list1.val <= list2.val) {
      curr.next = list1
      list1 = list1.next
    } else {
      curr.next = list2
      list2 = list2.next
    }
    curr = curr.next
  }
  curr.next = list1 ?? list2
  return resultSt.next
}
