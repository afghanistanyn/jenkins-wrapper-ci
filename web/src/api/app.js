import service from '@/utils/request'

// @Tags App
// @Summary 创建App
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.App true "创建App"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /app/createApp [post]
export const createApp = (data) => {
  return service({
    url: '/app/createApp',
    method: 'post',
    data
  })
}

// @Tags App
// @Summary 删除App
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.App true "删除App"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /app/deleteApp [delete]
export const deleteApp = (data) => {
  return service({
    url: '/app/deleteApp',
    method: 'delete',
    data
  })
}

// @Tags App
// @Summary 删除App
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除App"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /app/deleteApp [delete]
export const deleteAppByIds = (data) => {
  return service({
    url: '/app/deleteAppByIds',
    method: 'delete',
    data
  })
}

// @Tags App
// @Summary 更新App
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.App true "更新App"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /app/updateApp [put]
export const updateApp = (data) => {
  return service({
    url: '/app/updateApp',
    method: 'put',
    data
  })
}

// @Tags App
// @Summary 用id查询App
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.App true "用id查询App"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /app/findApp [get]
export const findApp = (params) => {
  return service({
    url: '/app/findApp',
    method: 'get',
    params
  })
}

// @Tags App
// @Summary 分页获取App列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取App列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /app/getAppList [get]
export const getAppList = (params) => {
  return service({
    url: '/app/getAppList',
    method: 'get',
    params
  })
}
