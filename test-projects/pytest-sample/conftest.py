"""Pytest configuration to ensure calculator module can be imported."""
import sys
from pathlib import Path

# Add parent directory to path so tests can import calculator
sys.path.insert(0, str(Path(__file__).parent))
