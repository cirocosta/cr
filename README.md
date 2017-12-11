<h1 align="center">cr 📂  </h1>

<h5 align="center">The concurrent runner</h5>

<br/>


------


```yaml
Config:
  Timeout: '30s'

Jobs:
  - Name: 'job1'

  - Name: 'job2'
    Timeout: '10s'
    DependsOn:
    - 'job1'
```

```golang


```


