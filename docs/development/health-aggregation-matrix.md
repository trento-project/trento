# Node health status aggregation matrix

The following table explains how we determine the aggregated, color-coded, health status of hosts.

| Heartbeat | Config Checks                | Aggregate Result
|-----------|------------------------------|------------------
| Ack       | Ran and passed               | Green
| Timeout   | \                            | Red
| Ack       | Ran and passed with warnings | Yellow
| Ack       | Ran and didn't pass          | Red
| Ack       | Didn't run at all            | Gray *
| Ack       | Failed to run                | Yellow *

Color codes legend:

- Red: Critical situation, user should immediately act
- Yellow: User be aware and know what they're doing
- Green: User can go to sleep fine and dandy
- Gray: Not enough information is available

> *: maybe allow user to pick the severity level for such case
