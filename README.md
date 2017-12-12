<h1 align="center">cr 📂  </h1>

<h5 align="center">The concurrent runner</h5>

<br/>


------


```yaml
Config:
  Timeout: '30s'

Jobs:
  # Jobs without `Run` field can be 
  # used as as synchronization barrier
  # as it counts as an entry in the execution
  # graph.
  - Name: 'job1'

  - Name: 'job2'
    Timeout: '10s'
    DependsOn:
    - 'job1'
```

### Templating

```
env     "arg"           -       retrieves the environment
                                variable `arg` from the environment
```

