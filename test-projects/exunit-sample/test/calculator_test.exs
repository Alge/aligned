defmodule CalculatorTest do
  use ExUnit.Case
  doctest Calculator

  test "add_positive_numbers" do
    assert Calculator.add(2, 3) == 5
  end

  test "add_negative_numbers" do
    assert Calculator.add(-1, -1) == -2
  end

  test "subtract_numbers" do
    assert Calculator.subtract(5, 3) == 2
  end

  test "multiply_numbers" do
    assert Calculator.multiply(4, 3) == 12
  end

  test "divide_numbers" do
    assert Calculator.divide(10, 2) == {:ok, 5.0}
  end

  test "divide_by_zero" do
    assert Calculator.divide(10, 0) == {:error, "Cannot divide by zero"}
  end
end
