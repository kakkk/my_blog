import { request } from './request';

export async function addPost(
  title: string,
  content: string,
  publish: boolean,
  categoryIds: number[],
  tags: string[],
  publishAt: number
) {
  return request({
    url: '/admin/post',
    method: 'POST',
    data: {
      title,
      content,
      publish,
      categories_id: categoryIds,
      tags,
      publish_at: publishAt,
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
  tags: string[],
  publishAt: number
) {
  return request({
    url: `/admin/post/${id}`,
    method: 'PUT',
    data: {
      title,
      content,
      categories_id: categoryIds,
      tags,
      publish_at: publishAt,
    },
  });
}

export async function updatePostPublish(id: number, publish: boolean) {
  return request({
    url: `/admin/post/${id}/publish`,
    method: 'PUT',
    data: {
      publish,
    },
  });
}

export async function getPostList(title: string, page: number, limit: number) {
  return request({
    url: `admin/posts?keyword=${title}&page=${page}&limit=${limit}`,
    method: 'GET',
  });
}

export async function deletePost(id: number) {
  return request({
    url: `/admin/post/${id}`,
    method: 'DELETE',
  });
}
