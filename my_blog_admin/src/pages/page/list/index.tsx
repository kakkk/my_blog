import React, { useEffect } from 'react';
import {
  Breadcrumb,
  Button,
  Card,
  Input,
  Message,
  Popconfirm,
  Table,
} from '@arco-design/web-react';
import { ColumnProps } from '@arco-design/web-react/es/Table/interface';
import { IconLink } from '@arco-design/web-react/icon';
import styles from './style/index.module.less';
import history from '../../../history';
import { deletePage, getPageList } from '../../../api/page';

export default function PageList() {
  const [colData, setColData] = React.useState([] as any[]);
  const fetchData = async () => {
    try {
      const res: any = await getPageList();
      if (res.code === 0) {
        if (res.data === null) {
          setColData([] as any[]);
        } else {
          setColData(res.data.page_list);
        }
      } else {
        Message.error(res.msg);
      }
    } finally {
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const onDelete = async (id: number) => {
    try {
      const res: any = await deletePage(id);
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
                window.open(`http://localhost:8888/pages/${record.slug}`);
              }}
            />
          </div>
        );
      },
    },
    {
      title: '缩略名',
      dataIndex: 'slug',
      align: 'center',
    },
    {
      title: '操作',
      width: 300,
      align: 'center',
      render: (_, record) => (
        <div className={styles.operations}>
          <Button
            type="text"
            size="small"
            onClick={() => {
              history.push(`/page/edit?id=${record.id}`);
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

  return (
    <div className={styles.container}>
      <Breadcrumb style={{ marginBottom: 20 }}>
        <Breadcrumb.Item>页面</Breadcrumb.Item>
        <Breadcrumb.Item>
          <a
            onClick={() => {
              history.push(`/page/list`);
            }}
            href="#"
          >
            页面管理
          </a>
        </Breadcrumb.Item>
      </Breadcrumb>
      <Card bordered={false}>
        <div className={styles.toolbar}>
          <Button
            type="primary"
            onClick={() => {
              history.push('/page/edit');
            }}
          >
            添加页面
          </Button>
          <Input.Search style={{ width: 300 }} searchButton placeholder="搜索" />
        </div>
        <Table rowKey="id" columns={columns} border={false} data={colData} pagination={false} />
      </Card>
    </div>
  );
}
