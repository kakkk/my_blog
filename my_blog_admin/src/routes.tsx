import React from 'react';
import {
  IconDashboard,
  IconEdit,
  IconFile,
  IconList,
  IconSettings,
  IconTag,
  IconUserGroup,
} from '@arco-design/web-react/icon';

import { CommentOutlined } from '@ant-design/icons';

export const defaultRoute = 'dashboard';

export const routes = [
  {
    name: '概览',
    key: 'dashboard',
    icon: <IconDashboard />,
    componentPath: 'dashboard',
  },
  {
    name: '文章',
    key: 'post',
    icon: <IconEdit />,
    children: [
      {
        name: '新建文章',
        key: 'post/edit',
        componentPath: 'post/edit',
      },
      {
        name: '全部文章',
        key: 'post/list',
        componentPath: 'post/list',
      },
    ],
  },
  {
    name: '页面',
    key: 'page',
    icon: <IconFile />,
    children: [
      {
        name: '新建页面',
        key: 'page/edit',
        componentPath: 'page/edit',
      },
      {
        name: '全部页面',
        key: 'page/list',
        componentPath: 'page/list',
      },
    ],
  },
  {
    name: '分类管理',
    key: 'category',
    icon: <IconList />,
    componentPath: 'category',
  },
  {
    name: '标签管理',
    key: 'tag',
    icon: <IconTag />,
    componentPath: 'tag',
  },
  {
    name: '评论管理',
    key: 'comment',
    componentPath: 'comment',
    icon: <CommentOutlined className="arco-icon" />,
  },
  {
    name: '用户管理',
    key: 'user',
    icon: <IconUserGroup />,
    componentPath: 'user',
  },
  {
    name: '设置',
    key: 'setting',
    icon: <IconSettings />,
    children: [
      {
        name: '系统设置',
        key: 'setting/system',
        componentPath: 'setting/system',
      },
      {
        name: '用户设置',
        key: 'setting/user',
        componentPath: 'setting/user',
      },
    ],
  },
  {
    key: '403',
    componentPath: 'exception/403',
  },
];
