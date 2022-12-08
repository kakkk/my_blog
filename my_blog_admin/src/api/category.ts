import { request } from './request';

export async function getCategoryList(parent: number) {
  return request({
    url: `admin/categories?parent=${parent}`,
    method: 'GET',
  });
}

export async function createCategory(data, parent) {
  const reqData = {
    name: data.name,
    slug: data.slug,
    description: data.description,
    parent_id: Number(parent),
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

export async function getCategoryById(id: number) {
  return request({
    url: `admin/category/${id}`,
    method: 'GET',
  });
}

export async function updateCategoryById(id: number, data: any) {
  return request({
    url: `admin/category/${id}`,
    method: 'PUT',
    data: {
      name: data.name,
      slug: data.slug,
      description: data.description,
      parent_id: Number(data.parent),
    },
  });
}
