import calculator
import gleeunit
import gleeunit/should

pub fn main() {
  gleeunit.main()
}

pub fn add_positive_numbers_test() {
  calculator.add(2, 3)
  |> should.equal(5)
}

pub fn add_negative_numbers_test() {
  calculator.add(-1, -1)
  |> should.equal(-2)
}

pub fn subtract_numbers_test() {
  calculator.subtract(5, 3)
  |> should.equal(2)
}

pub fn multiply_numbers_test() {
  calculator.multiply(4, 3)
  |> should.equal(12)
}

pub fn divide_numbers_test() {
  calculator.divide(10.0, 2.0)
  |> should.equal(Ok(5.0))
}

pub fn divide_by_zero_test() {
  calculator.divide(10.0, 0.0)
  |> should.equal(Error("Cannot divide by zero"))
}
