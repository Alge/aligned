import { describe, it, expect } from 'vitest'

describe('Math operations', () => {
  describe('Addition', () => {
    it('adds 1 + 2 to equal 3', () => {
      expect(1 + 2).toBe(3)
    })

    it('adds negative numbers', () => {
      expect(-1 + -2).toBe(-3)
    })
  })

  describe('Subtraction', () => {
    it('subtracts 5 - 3 to equal 2', () => {
      expect(5 - 3).toBe(2)
    })
  })
})

describe('String operations', () => {
  it('concatenates strings', () => {
    expect('hello' + ' ' + 'world').toBe('hello world')
  })
})

