# SQLEyes 核心引擎

该项目为SQLEyes的引擎部分，该部分逻辑是报文数据的截取转发逻辑。包含了一个插件的引入机制，其他功能均由插件实现

# 引擎的基本功能

- 引擎初始化会加载配置文件，引入的插件会自动注册到引擎中
- 配置文件中配置会被解析并覆盖插件的默认配置

### 插件列表:
- **[plugin-postgresql](https://github.com/sqleyes/plugin-postgresql)**
- **[plugin-mysql](https://github.com/sqleyes/plugin-mysql)**