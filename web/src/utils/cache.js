/**
 * 闭包实现缓存
 * 属性: key -value 将数据缓存在一个对象中
 * 通过setCache进行设置缓存，getCache 读取缓存
 * @returns
 */
const configCache = function () {
  const _cache = new Map()
  return {
    setCache(k, v) {
      _cache[k] = v
      return this
    },
    getCache(k) {
      return _cache[k]
    },
    removeCache(k) {
      delete _cache(k)
      return this
    }
  }
}

const newConfigCache = configCache()

export default newConfigCache
