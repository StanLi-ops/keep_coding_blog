import request from '../utils/request'

// 用户相关接口
export const login = (data) => {
  return request({
    url: '/login',
    method: 'POST',
    data
  })
}

export const register = (data) => {
  return request({
    url: '/register',
    method: 'POST',
    data
  })
}

// 文章相关接口
export const getPosts = (params) => {
  return request({
    url: '/posts',
    method: 'GET',
    params
  })
}

export const getPostById = (id) => {
  return request({
    url: `/posts/${id}`,
    method: 'GET'
  })
}

export const createPost = (data) => {
  return request({
    url: '/posts',
    method: 'POST',
    data
  })
}

export const updatePost = (id, data) => {
  return request({
    url: `/posts/${id}`,
    method: 'PUT',
    data
  })
}

export const deletePost = (id) => {
  return request({
    url: `/posts/${id}`,
    method: 'DELETE'
  })
}

// 评论相关接口
export const getComments = (params) => {
  return request({
    url: '/comments',
    method: 'GET',
    params
  })
}

export const createComment = (data) => {
  return request({
    url: '/comments',
    method: 'POST',
    data
  })
}

export const updateComment = (id, data) => {
  return request({
    url: `/comments/${id}`,
    method: 'PUT',
    data
  })
}

export const deleteComment = (id) => {
  return request({
    url: `/comments/${id}`,
    method: 'DELETE'
  })
}

// 标签相关接口
export const getTags = () => {
  return request({
    url: '/tags',
    method: 'GET'
  })
} 