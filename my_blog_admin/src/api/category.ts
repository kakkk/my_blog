import { request } from './request';

export async function getCategoryList() {
  return request({
    url: `admin/category/list`,
    method: 'GET',
  });
}

export async function createCategory(data) {
  const reqData = {
    name: data.name,
    slug: data.slug,
  };
  console.error(reqData);
  return request({
    url: `admin/category`,
    method: 'POST',
    data: reqData,
  });
}

export async function deleteCategory(id: number) {
  return request({
    url: `admin/category/${id}`,
    method: 'DELETE',
  });
}

export async function updateCategoryById(id: number, data: any) {
  return request({
    url: `admin/category/${id}`,
    method: 'PUT',
    data: {
      name: data.name,
      slug: data.slug,
    },
  });
}
