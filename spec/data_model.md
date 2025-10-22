# Data Model

The data model represents the in-memory structure of parsed specifications. The core types are `Specification` (a document) and `Section` (a hierarchical node representing a heading and its content).

## Section Properties

### Identify leaf sections

A section is a leaf if it has no children. Leaf sections are the only sections that require test references.

**Test:** `Alge/aligned/internal/spec.TestSection_IsLeaf`

### Check if section has test reference

Determine whether a section has an associated test by checking if the TestName field is populated.

**Test:** `Alge/aligned/internal/spec.TestSection_HasTest`

### Maintain parent-child relationships

Sections form a bidirectional tree structure where children reference their parent and parents reference their children. This enables traversal in both directions.

**Test:** `Alge/aligned/internal/spec.TestSection_ParentChildRelationships`

### Determine if section requires test reference

Only leaf sections that are not interfaces (or descendants of interfaces) require test references. Parent sections and interface-related sections do not.

**Test:** `Alge/aligned/internal/spec.TestSection_RequiresTest`

## Specification Queries

### Find all leaf sections in specification tree

Traverse the entire specification hierarchy and return all leaf sections, regardless of nesting depth.

**Test:** `Alge/aligned/internal/spec.TestSpecification_AllLeaves`

### Get required tests from specification

Extract all test names from leaf sections that have test references. This provides the list of tests that should exist according to the specification.

**Test:** `Alge/aligned/internal/spec.TestSpecification_RequiredTests`

## Interface System

### Detect interface markers

Identify sections marked with `[INTERFACE]` in their title. These sections define contracts that implementations must fulfill.

**Test:** `Alge/aligned/internal/spec.TestDetectInterfaceMarker`

### Skip test requirements for interface sections

Interface sections and all their children should not require test references, as they define structure rather than testable behavior.

**Test:** `Alge/aligned/internal/spec.TestInterfaceSkipsTestRequirements`

### Detect implementation markers

Identify sections marked with `[IMPLEMENTS: InterfaceName]` in their title. These sections declare that they implement a specific interface.

**Test:** `Alge/aligned/internal/spec.TestDetectImplementationMarker`

### Extract interface name from implementation marker

Parse the interface name from `[IMPLEMENTS: InterfaceName]` markers to enable validation against the referenced interface.

**Test:** `Alge/aligned/internal/spec.TestExtractInterfaceName`

### Validate implementation structure matches interface

Verify that implementations contain all required sections from the interface. Section matching is case-insensitive to allow for natural language variations.

**Test:** `Alge/aligned/internal/spec.TestValidateImplementationStructure`

### Allow additional sections in implementations

Implementations can include sections beyond what the interface requires. Only missing interface sections are validation errors.

**Test:** `Alge/aligned/internal/spec.TestAllowAdditionalSectionsInImplementation`
