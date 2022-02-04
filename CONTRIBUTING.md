# Contribution guidelines

## Code formatting

- For go code, simply use go conventions (can be formatted via `go fmt`)
- Avoid _when possible_ lines of 80+ chars
- Every file should end with an empty line

## Git

- Branch names must follow this pattern: `<type>/name-of-related-issue`
- Commit messages must be 50 chars max (add a description if necessary)
- Commit messages must start with a verb in the present simple, such as `Add ...`, `Implement ...`, `Fix ...`
- Commit messages may be prefixed with a type, such as `feat:`, `fix:`, ...

## Pull requests & Merge policy

- Do not create a Pull request if it is not directly related to an existing issue. Create an issue first if a change is needed.
- Merge must be performed using the method `squash`
- A Pull request must have been tested and reviewed before merging:

  - all tests must succeed
  - all comments must have been addressed
  - **2** reviewers at least must have approved the changes

- PR Status:

  - draft -> WIP (no reviews)
  - open -> ready for review

- No push on someone else's PR

- Reviewers should prefix their comments with one of these keywords:
  - [suggestion] -> suggestion, could be ignored
  - [nitpick] -> as a suggestion, but even less important
  - [change needed] -> not mergable before the code is changed (or the reviewer changes their mind)
  - [question] -> question, does not necessary imply that a modification is needed
