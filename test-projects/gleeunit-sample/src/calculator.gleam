pub fn add(a: Int, b: Int) -> Int {
  a + b
}

pub fn subtract(a: Int, b: Int) -> Int {
  a - b
}

pub fn multiply(a: Int, b: Int) -> Int {
  a * b
}

pub fn divide(a: Float, b: Float) -> Result(Float, String) {
  case b {
    0.0 -> Error("Cannot divide by zero")
    _ -> Ok(a /. b)
  }
}
