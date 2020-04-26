# Changelog

**2.3.1**
- Added documentation on exported functions
- Unexported several functions that shouldn't have been necessary for package use

**2.2.2**
- Added missing 0s from solution vector

**2.2.1**
- b vector column in tableau always positive
- fixed infinite loop due to ratio test accepting `-0` (value in tableau is negative; b entry is 0)

**2.1.3**
- fixed bug in ratio test where a ratio of 0 is rejected.

**2.1.2**
- more small float errors.

**2.1.1**
- removed useless return in `SetObjectiveFunction`

**2.0.2**
- Fixed bug in which feasible LPs were labelled infeasible due to small float

**2.0.1**
- Allow non-equality constraints
- Changed signature on adding constraints to reflect mathematical representations.
- NewLP returns pointer to LP instead of the actual LP
- Renamed `GetGoal` to `GetObjective`
- Renamed function `ObjX` to `ObjectiveX`

**1.4.1**
- Added min LPs

**1.3.1**
- Changed the way test LPs are read
- Solving negates objective function in a later stage.

**1.2.2**
- Changes to get ready for minimization; minor bug fix

**1.2.1**
- Basic feasible solutions can be extracted from LPs

**1.1.1**
- Renamed AddConstraint to AddConstraintEq to allow for other constraint types when support for non-SEF LPs is added.

**1.0.1**
- Project made
