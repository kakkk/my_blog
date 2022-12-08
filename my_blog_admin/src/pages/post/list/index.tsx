import React, { useEffect, useState } from 'react';
import {
  Breadcrumb,
  Button,
  Card,
  Input,
  Message,
  Popconfirm,
  Select,
  Table,
  Tag,
} from '@arco-design/web-react';
import { ColumnProps } from '@arco-design/web-react/es/Table/interface';
import { IconLink, IconUser } from '@arco-design/web-react/icon';
import { PaginationProps } from '@arco-design/web-react/es/Pagination/pagination';
import history from '../../../history';
import styles from './style/index.module.less';
import { formatDate } from '../../../utils/formatDate';
import getUrlParams from '../../../utils/getUrlParams';
import { deletePost, getPostList } from '../../../api/post';

export default function PostList(props) {
  const options = ['未发布', '已发布'];
  const [colData, setColData] = useState(
    new Array(0) as {
      id: number;
      title: string;
      categories: string;
      editor: string;
      comment: number;
      publish: boolean;
      updateAt: number;
      publishAt: number;
    }[]
  );
  const [searchTitle, setSearchTitle] = useState('');
  const [pagination, setPagination] = useState({
    sizeCanChange: true,
    showTotal: true,
    pageSize: 10,
    current: 1,
    total: 0,
    pageSizeChangeResetCurrent: true,
  } as PaginationProps);
  const [loading, setLoading] = useState(false);
  const makeCategory = (item) => {
    if (item.length === 0) {
      return '无';
    }
    let str: string = item[0].name;
    for (let i = 1; i < item.length; i++) {
      str = `${str}、${item[i].name}`;
    }
    return str;
  };
  const fetchData = async () => {
    setLoading(true);
    const urlParams = getUrlParams();
    const page = urlParams.page ? Number(urlParams.page) : 0;
    const size = urlParams.size ? Number(urlParams.size) : 10;
    const title = urlParams.title ? urlParams.title : '';
    setSearchTitle(title);
    try {
      const res: any = await getPostList(title, page, size);
      if (res.code === 0) {
        const newColData: {
          id: number;
          title: string;
          categories: string;
          editor: string;
          comment: number;
          publish: boolean;
          updateAt: number;
          publishAt: number;
        }[] = new Array(0);
        for (const datum of res.data) {
          const categoryStr = makeCategory(datum.categories);
          const d = {
            id: datum.id,
            categories: categoryStr,
            comment: datum.comment_count,
            editor: datum.editor_nickname,
            publish: datum.publish,
            publishAt: datum.publish_at,
            title: datum.title,
            updateAt: datum.update_at,
          };
          newColData.push(d);
        }
        setColData(newColData);
        setPagination({
          sizeCanChange: pagination.sizeCanChange,
          showTotal: pagination.showTotal,
          pageSize: pagination.pageSize,
          current: res.pagination.page,
          total: res.pagination.total,
          pageSizeChangeResetCurrent: pagination.pageSizeChangeResetCurrent,
        });
      } else {
        Message.error(res.msg);
      }
    } finally {
      setLoading(false);
    }
  };
  const onDelete = async (id: number) => {
    try {
      const res: any = await deletePost(id);
      if (res.code === 0) {
        fetchData();
        Message.success('删除成功');
      } else {
        Message.error(res.msg);
      }
    } finally {
    }
  };

  const columns: ColumnProps[] = [
    {
      title: '标题',
      dataIndex: 'title',
      width: 300,
      render: (_, record) => {
        return (
          <div>
            {record.title}
            <Button
              type="text"
              icon={<IconLink />}
              onClick={() => {
                window.open(`http://localhost:8080/api/v1/post/${record.id}`);
              }}
            />
          </div>
        );
      },
    },
    {
      title: '作者',
      dataIndex: 'editor',
      align: 'center',
      render: (item) => {
        return <Tag icon={<IconUser />}>{item}</Tag>;
      },
    },
    {
      title: '分类',
      dataIndex: 'categories',
      align: 'center',
    },
    {
      title: '评论',
      dataIndex: 'comment',
      align: 'center',
    },
    {
      title: '状态',
      dataIndex: 'publish',
      align: 'center',
      render: (item) => {
        return (
          <Select
            options={options}
            defaultValue={item === true ? '已发布' : '未发布'}
            style={{ width: 100 }}
          />
        );
      },
    },
    {
      title: '更新时间',
      dataIndex: 'updateAt',
      align: 'center',
      render: (item) => {
        return <div>{formatDate(item * 1000, 'yyyy-MM-dd HH:ss:mm')}</div>;
      },
    },
    {
      title: '操作',
      width: 200,
      align: 'center',
      dataIndex: 'operations',
      render: (_, record) => (
        <div className={styles.operations}>
          <Button
            type="text"
            size="small"
            onClick={() => {
              history.push(`/post/edit?id=${record.id}`);
            }}
          >
            修改
          </Button>
          <Popconfirm
            title="确认删除？"
            onOk={() => {
              onDelete(record.id);
            }}
          >
            <Button type="text" size="small">
              删除
            </Button>
          </Popconfirm>
        </div>
      ),
    },
  ];
  useEffect(() => {
    fetchData();
  }, [props.location]);

  const onChangeTable = (pagination) => {
    setPagination(pagination);
    history.push(
      `/post/list?page=${pagination.current}&size=${pagination.pageSize}${
        searchTitle === '' ? `&${searchTitle}` : ''
      }`
    );
  };
  return (
    <div className={styles.container}>
      <Breadcrumb style={{ marginBottom: 20 }}>
        <Breadcrumb.Item>文章</Breadcrumb.Item>
        <Breadcrumb.Item>
          <a
            onClick={() => {
              history.push(`/post/list`);
            }}
            href="#"
          >
            文章管理
          </a>
        </Breadcrumb.Item>
      </Breadcrumb>
      <Card bordered={false}>
        <div className={styles.toolbar}>
          <Button type="primary">添加文章</Button>
          <Input.Search style={{ width: 300 }} searchButton placeholder="搜索" />
        </div>
        <Table
          rowKey="id"
          columns={columns}
          border={false}
          data={colData}
          pagination={pagination}
          loading={loading}
          onChange={onChangeTable}
        />
      </Card>
    </div>
  );
}
