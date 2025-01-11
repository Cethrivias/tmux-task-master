# TODO
- [x] Add task folder to a hardcoded folder
- [x] Create a project worktree to a task in a hardcoded folder
- [x] Read config to get the folder
- [x] List tasks
- [x] List project worktrees
- [x] Delete project worktrees
- [x] Delete tasks
- [x] Integration tests
    - [x] Create task
    - [x] List tasks
    - [x] Add worktree to a task
    - [x] Delete worktree
    - [x] Delete task
- [x] Show worktree branches when running list
- [x] Check if the project is already added to the task when running add command
- [x] Cleanup this garbage code
    - [x] Reinvent the wheel by creating my own CLI framework
    - [x] A lot of repetition in tests and code. Inconsistent naming for the same things (good enough now)
        - [x] Task "class"
- [x] Parallelise list projects (no reason to complicate things for now)
- [x] Not possible to delete a worktree if CWD is not a part of that worktree's repo
- [ ] Add descriptive error messages
    ```
    ❯ ttm delete TASK-001                                                              1.25s
    This task contains 1 projects. Do you want to delete it? (y/n)
    y
    Deleting task 'TASK-001' projects:
     - my.fancy.project
    fatal: not a git repository (or any of the parent directories): .git

    2024/11/11 10:28:53 exit status 128
    ```
- [ ] Adding a repo with a conflicting name causes an error.
    - [x] Error message is missing ?!?!
    - [ ] Ask for a different name
- [ ] Start adding tmux integration
    - [x] Create session and window
        - [ ] Add tests
