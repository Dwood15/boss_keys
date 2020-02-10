# bosskeys

This package is an experiment in tree-representation of game progression. 

---
## Initial Goals:

1) Represent progression through a game (OoT in first case) as a series of progressions through locks and keys. 
2) Be able to traverse a vanilla progression chart via breadth-first and depth-first searches.

### Hints

- The most interesting aspects are under base_pools, in the json files. 
- There is main execution method at this time. Use the tests. 

---
## On Language choice

- Golang was only chosen because it's what I'm most familiar with and it's easier for me to use and test.
- In the future for this project, other possible language choices are: Nim, C/C++, Clojure, Haskell.

---

### License Choice:

- I chose LGPL because I'm building a generic tool. These algorithms and data structures should work for more than just one use-case.

- While inspired by randomizers, the codebase is currently my own... contributions greatly desired, as I'm probably in over my head with this.