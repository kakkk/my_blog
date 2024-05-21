import { request } from './request';

export async function getCommentList(page: number, limit: number) {
  return request({
    url: `admin/comment/list?page=${page}&limit=${limit}`,
    method: 'GET',
  });
}

export async function updateCommentStatus(id: number, status: number) {
  return request({
    url: `/admin/comment/${id}/status`,
    method: 'PUT',
    data: {
      status,
    },
  });
}

export async function deleteComment(id: number) {
  return request({
    url: `/admin/comment/${id}`,
    method: 'DELETE',
  });
}
