defmodule Calculator do
  @moduledoc """
  Simple calculator module for demonstration.
  """

  @doc """
  Adds two numbers.
  """
  def add(a, b), do: a + b

  @doc """
  Subtracts the second number from the first.
  """
  def subtract(a, b), do: a - b

  @doc """
  Multiplies two numbers.
  """
  def multiply(a, b), do: a * b

  @doc """
  Divides the first number by the second.
  Raises an error if dividing by zero.
  """
  def divide(_a, 0), do: {:error, "Cannot divide by zero"}
  def divide(a, b), do: {:ok, a / b}
end
