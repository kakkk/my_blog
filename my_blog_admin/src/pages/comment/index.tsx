import React, { useEffect, useState } from 'react';
import {
  Breadcrumb,
  Button,
  Card,
  Message,
  Popconfirm,
  Select,
  Table,
} from '@arco-design/web-react';
import { IconLink } from '@arco-design/web-react/icon';
import { PaginationProps } from '@arco-design/web-react/es/Pagination/pagination';
import styles from './style/index.module.less';
import history from '../../history';
import getUrlParams from '../../utils/getUrlParams';
import { deleteComment, getCommentList, updateCommentStatus } from '../../api/comment';

export default function Comment(props) {
  const [pagination, setPagination] = useState({
    sizeCanChange: true,
    showTotal: true,
    pageSize: 10,
    current: 1,
    total: 0,
    pageSizeChangeResetCurrent: true,
  } as PaginationProps);
  const onChangeTable = (pagination) => {
    setPagination(pagination);
    history.push(`/comment?page=${pagination.current}&size=${pagination.pageSize}`);
  };
  const options = ['待审', '通过', '封禁'];
  const [colData, setColData] = useState(
    new Array(0) as {
      id: string;
      nickname: string;
      avatar: string;
      website: string;
      article: {
        id: number;
        title: string;
      };
      content: string;
      status: number;
      createAt: number;
    }[]
  );
  const [loading, setLoading] = useState(false);

  const fetchData = async () => {
    setLoading(true);
    const urlParams = getUrlParams();
    const page = urlParams.page ? Number(urlParams.page) : 0;
    const size = urlParams.size ? Number(urlParams.size) : 10;
    try {
      const res: any = await getCommentList(page, size);
      if (res.code === 0) {
        const newColData: {
          id: string;
          nickname: string;
          avatar: string;
          website: string;
          article: {
            id: number;
            title: string;
          };
          content: string;
          status: number;
          createAt: number;
        }[] = new Array(0);
        if (res.data.comment) {
          for (const datum of res.data.comment) {
            const d = {
              id: datum.id,
              nickname: datum.nickname,
              avatar: datum.avatar,
              website: datum.website,
              article: {
                id: datum.article.id,
                title: datum.article.title,
              },
              content: datum.content,
              status: datum.status,
              createAt: datum.create_at,
            };
            newColData.push(d);
          }
        }

        setColData(newColData);
        setPagination({
          sizeCanChange: pagination.sizeCanChange,
          showTotal: pagination.showTotal,
          pageSize: pagination.pageSize,
          current: res.data.pagination.page,
          total: res.data.pagination.total,
          pageSizeChangeResetCurrent: pagination.pageSizeChangeResetCurrent,
        });
      } else {
        Message.error(res.msg);
      }
    } finally {
      setLoading(false);
    }
  };

  const onChangeStatus = async (id: number, value: string) => {
    let status: number;
    if (value === '封禁') {
      status = 3;
    } else if (value === '通过') {
      status = 2;
    } else {
      status = 1;
    }
    try {
      const res: any = await updateCommentStatus(id, status);
      if (res.code === 0) {
        fetchData();
        Message.success('操作成功');
      } else {
        Message.error(res.msg);
      }
    } finally {
    }
  };

  const onDelete = async (id: number) => {
    try {
      const res: any = await deleteComment(id);
      if (res.code === 0) {
        fetchData();
        Message.success('删除成功');
      } else {
        Message.error(res.msg);
      }
    } finally {
    }
  };

  const columns: any[] = [
    {
      title: '昵称',
      dataIndex: 'nickname',
      width: 100,
    },
    {
      title: '评论',
      dataIndex: 'content',
      align: 'center',
      width: 350,
    },
    {
      title: '文章',
      dataIndex: 'article',
      align: 'center',
      render: (item) => {
        return (
          <div>
            {item.title}
            <Button
              type="text"
              icon={<IconLink />}
              onClick={() => {
                window.open(`http://localhost:8888/archives/${item.id}`);
              }}
            />
          </div>
        );
      },
    },
    {
      title: '状态',
      dataIndex: 'status',
      align: 'center',
      render: (_, record) => {
        return (
          <Select
            options={options}
            defaultValue={record.status === 2 ? '通过' : record.status === 1 ? '待审' : '封禁'}
            onChange={(value) => {
              onChangeStatus(record.id, value);
            }}
            style={{ width: 80 }}
          />
        );
      },
    },
    {
      title: '操作',
      width: 200,
      align: 'center',
      dataIndex: 'operations',
      render: (_, record) => (
        <div className={styles.operations}>
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

  return (
    <div className={styles.container}>
      <Breadcrumb style={{ marginBottom: 20 }}>
        <Breadcrumb.Item>
          <a
            onClick={() => {
              history.push(`/comment`);
            }}
            href="#"
          >
            评论管理
          </a>
        </Breadcrumb.Item>
      </Breadcrumb>
      <Card bordered={false}>
        <Table
          rowKey="id"
          columns={columns}
          border={false}
          pagination={pagination}
          loading={loading}
          data={colData}
          onChange={onChangeTable}
        />
      </Card>
    </div>
  );
}
