import { request } from './request';

export async function addPost(
  title: string,
  content: string,
  publish: boolean,
  categoryIds: number[],
  tags: string[]
) {
  return request({
    url: '/admin/post',
    method: 'POST',
    data: {
      title,
      content,
      category_list: categoryIds,
      tags,
      status: publish ? 2 : 1,
    },
  });
}

export async function getPost(id: number) {
  return request({
    url: `admin/post/${id}`,
    method: 'GET',
  });
}

export async function updatePost(
  id: number,
  title: string,
  content: string,
  categoryIds: number[],
  tags: string[]
) {
  return request({
    url: `/admin/post/${id}`,
    method: 'PUT',
    data: {
      title,
      content,
      category_list: categoryIds,
      tags,
    },
  });
}

export async function updatePostStatus(id: number, status: number) {
  return request({
    url: `/admin/post/${id}/status`,
    method: 'PUT',
    data: {
      status,
    },
  });
}

export async function getPostList(
  title: string,
  categories: string[],
  tags: string[],
  page: number,
  limit: number
) {
  return request({
    url: `admin/post/list`,
    method: 'POST',
    data: {
      keyword: title,
      page,
      limit,
      categories,
      tags,
    },
  });
}

export async function deletePost(id: number) {
  return request({
    url: `/admin/post/${id}`,
    method: 'DELETE',
  });
}
