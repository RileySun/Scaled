# Three-Tier Architecture
A horizontal scaling architecture design.
Decouples handler, business, and data logic into three seperate modules that can be scaled independently of one another.
(I find this architecture design to be useful, but not apropriate for every use case)

**Advantages**:
* Decoupled modules perfect for scaling horizontally
* Message Queuing for module communication (To-Do)
* Data store caching