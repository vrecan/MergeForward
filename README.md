# MergeForward
Simple tool to merge configuration files taking values from the src and adding them to the destination file while still keeping the new keys.

Example YAML:
```
key: custom
key2: default
list:
  - !!OBJECT
  key: custom
  key2: default
  - !!OBJECT
  key: default
  key2: default  
```

destination conf:

```
key: default
key2: default
new: default
list:
  - !!OBJECT
  key: default
  key2: default
  key3: default
  - !!OBJECT
  key: default
  key2: default  
  - !!OBJECT
  key: default
  key2: default 
```

result:
```
key: custom
key2: default
new: default
list:
  - !!OBJECT
  key: custom
  key2: default
  key3: default
  - !!OBJECT
  key: default
  key2: default  
  - !!OBJECT
  key: default
  key2: default 
```
