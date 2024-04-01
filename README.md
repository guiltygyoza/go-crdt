# CRDT playground

This repo started as a fork of https://github.com/neurodrone/crdt.

### Infinite-phase set
- reference: https://arxiv.org/pdf/2304.01929.pdf
- source: [infphase_set.go](./infphase_set.go)
- test:
```
% go test -run -v "(TestInfPhase*)"
```
- todo: implement and test `merge()`
