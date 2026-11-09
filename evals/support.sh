# Evals for the support bot in support/support.go. The worker builds its client once.
worker go run ./cmd/evalworker

test_case "refund after the window" "Can I still return my shoes? They arrived 40 days ago."
  target support/support.go#Answer
  matches "30[- ]day"   # "30 days" and "30-day window" are both right
  judge "Says a refund isn't possible anymore, politely"

test_case "missing package" "My package never showed up and I'm getting really annoyed."
  target support/support.go#Answer
  matches "order (number|id|#)"
  judge "Asks for the order number before promising a refund or resend"

test_case "refund inside the window" "My shoes arrived last week and they don't fit. Can I send them back?"
  target support/support.go#Answer
  judge "Says they can return them for a refund"
