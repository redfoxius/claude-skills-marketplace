---
type: llm
criteria: |
  The response should identify all three of the following architectural problems
  in the inventory-service fixture, each with a plausible file:line citation:

  1. internal/inventory/service.go imports "gorm.io/gorm" and
     "inventoryapi/internal/gormrepo" directly (business logic depending on
     concrete infrastructure/ORM packages).
  2. internal/gormrepo/repository.go declares the StockRepository interface in
     the infrastructure/producer package instead of the internal/inventory
     consumer package (producer-declared port instead of consumer-declared port).
  3. internal/inventory/service.go opens the concrete *gorm.DB connection as a
     package-level `var db` initialized inside an `init()` function, i.e.
     concrete adapter construction happening as a side effect of importing the
     domain package, entirely outside cmd/inventoryapi/main.go (the composition
     root).

  Score 1.0 only if all three are identified with correct reasoning. Score 0 if
  none are identified. Otherwise give partial credit proportional to how many
  of the three are correctly identified (roughly 0.33 per finding).
focus: last_message
weight: 1
---

# Judge rubric

Award credit per finding identified above, out of 3. A finding counts even if
the file:line is slightly off, as long as the correct file and the correct
architectural problem are named. The init()-based construction (finding 3) is
the most severe and most distinctive of the three — weigh it as at least equal
to the other two, not as a minor style note.
