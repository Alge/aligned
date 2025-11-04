"""Tests for the calculator module."""

import pytest
from calculator import add, subtract, multiply, divide


class TestAddition:
    """Test suite for addition operations."""

    def test_add_positive_numbers(self):
        """Test adding two positive numbers."""
        assert add(2, 3) == 5

    def test_add_negative_numbers(self):
        """Test adding two negative numbers."""
        assert add(-1, -1) == -2


class TestSubtraction:
    """Test suite for subtraction operations."""

    def test_subtract_numbers(self):
        """Test subtracting numbers."""
        assert subtract(5, 3) == 2


class TestMultiplication:
    """Test suite for multiplication operations."""

    def test_multiply_numbers(self):
        """Test multiplying numbers."""
        assert multiply(4, 3) == 12


class TestDivision:
    """Test suite for division operations."""

    def test_divide_numbers(self):
        """Test dividing numbers."""
        assert divide(10, 2) == 5

    def test_divide_by_zero(self):
        """Test that dividing by zero raises an error."""
        with pytest.raises(ValueError, match="Cannot divide by zero"):
            divide(10, 0)
