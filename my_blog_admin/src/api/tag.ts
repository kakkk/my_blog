import { request } from './request';

export async function getTagList(keyword: string, page: number, limit: number) {
  return request({
    url: `admin/tags?keyword=${keyword}&page=${page}&limit=${limit}`,
    method: 'GET',
  });
}

export async function createTag(name: string) {
  return request({
    url: `admin/tag`,
    method: 'POST',
    data: {
      name,
    },
  });
}

export async function updateTag(id: number, name: string) {
  return request({
    url: `admin/tag/${id}`,
    method: 'PUT',
    data: {
      name,
    },
  });
}

export async function deleteTag(id: number) {
  return request({
    url: `admin/tag/${id}`,
    method: 'DELETE',
  });
}
