const timeFn = (date) => {
  var endDate = new Date() // 结束时间
  var timeDifference = endDate.getTime() - new Date(date).getTime() // 时间差
  var days = Math.floor(timeDifference / (24 * 3600 * 1000))
  var leave1 = timeDifference % (24 * 3600 * 1000) //计算天数后剩余的毫秒数
  var hours = Math.floor(leave1 / (3600 * 1000))
  //计算相差分钟数
  var leave2 = leave1 % (3600 * 1000) //计算小时数后剩余的毫秒数
  var minutes = Math.floor(leave2 / (60 * 1000))
  //计算相差秒数
  var leave3 = leave2 % (60 * 1000) //计算分钟数后剩余的毫秒数
  var seconds = Math.round(leave3 / 1000)
  if (days === 0) {
    return `${hours}小时${minutes}分钟${seconds}秒`
  } else if (hours === 0) {
    return `${minutes}分钟${seconds}秒`
  } else if (minutes === 0) {
    return `${seconds}秒`
  } else {
    return `${days}天${hours}小时${minutes}分钟${seconds}秒`
  }
}

export default timeFn
