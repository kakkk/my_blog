import { request } from './request';

export async function getUserInfo() {
  return request({
    url: 'admin/user/info',
    method: 'GET',
  });
}

export async function updateUserInfo(info: any) {
  return request({
    url: `admin/user/info`,
    method: 'PUT',
    data: {
      nickname: info.nickname,
      email: info.email,
      avatar: info.avatar,
      description: info.description,
    },
  });
}

export async function updateUserInfoById(id: number, info: any) {
  return request({
    url: `admin/user/${id}/info`,
    method: 'PUT',
    data: {
      nickname: info.nickname,
      email: info.email,
      role: info.role,
    },
  });
}

export async function createUser(user: any) {
  return request({
    url: `admin/user`,
    method: 'POST',
    data: {
      username: user.username,
      nickname: user.nickname,
      email: user.email,
      role: user.role,
      avatar: 'https://cdn.kakkk.net/img/kakkk.jpg',
    },
  });
}

export async function updateUserPassword(pwd: string, newPwd: string) {
  return request({
    url: `admin/user/password`,
    method: 'PUT',
    data: {
      password: pwd,
      new_password: newPwd,
    },
  });
}

export async function getUserList(
  username: string,
  nickname: string,
  email: string,
  page: number,
  limit: number
) {
  return request({
    url: `admin/users?username=${username}&nickname=${nickname}&email=${email}&page=${page}&limit=${limit}`,
    method: 'GET',
  });
}

export async function resetUserPasswordById(id: number) {
  return request({
    url: `admin/user/${id}/password`,
    method: 'PUT',
    data: {},
  });
}
