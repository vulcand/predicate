Predicate
=========

Library for creating predicate mini-languages in Go. 
This helpful for creating mini configuration languages, e.g. for setting up various thresholds.

Using this library you can create interpreted languages like this one:

```
LatencyMs() > 50 || ErrorRate > 0.2
```

