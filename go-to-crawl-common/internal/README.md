### internal层目录结构设计
- 把层次后缀带入到包名后缀
```
|_ logic
  |_ task: 
    |_ **1task
    |_ **2task
|_ service:
  |_ **1service
  |_ **2service
```

### 优势
- 在分层架构中，上层调用底层，import依赖的时候，包名冲突需要别名区分，操作更麻烦（尤其对于小项目逐渐变大项目，频繁变动目录结构的场景）
- A层调用B层，B调用C层，BC层都有同名包名P，A层同时调用BC层的时候，即要你选择导入哪一层，又要给其中一层取别名，略显麻烦