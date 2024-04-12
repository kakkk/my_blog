import { request } from './request';

export async function addPage(title: string, content: string, slug: string) {
  return request({
    url: '/admin/page',
    method: 'POST',
    data: {
      title,
      content,
      slug,
    },
  });
}

export async function getPage(id: number) {
  return request({
    url: `admin/page/${id}`,
    method: 'GET',
  });
}

export async function editPage(id: number, title: string, content: string, slug: string) {
  return request({
    url: `admin/page/${id}`,
    method: 'PUT',
    data: {
      title,
      content,
      slug,
    },
  });
}

export async function getPageList() {
  return request({
    url: 'admin/page/list',
    method: 'GET',
  });
}

export async function deletePage(id: number) {
  return request({
    url: `admin/page/${id}`,
    method: 'DELETE',
  });
}
