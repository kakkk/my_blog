import React from 'react';
import {
  Breadcrumb,
  Button,
  Card,
  Form,
  Grid,
  Input,
  Popconfirm,
  Table,
} from '@arco-design/web-react';
import { IconLink, IconSearch } from '@arco-design/web-react/icon';
import styles from './style/index.module.less';
import history from '../../history';

export default function Comment() {
  const { Row, Col } = Grid;
  const columns: any[] = [
    {
      title: '用户',
      dataIndex: 'user',
      width: 100,
    },
    {
      title: '内容',
      dataIndex: 'content',
      align: 'center',
      width: 350,
    },
    {
      title: '文章',
      dataIndex: 'post',
      align: 'center',
      render: (item) => {
        return (
          <div>
            {item}
            <Button type="text" icon={<IconLink />} />
          </div>
        );
      },
    },
    {
      title: '时间',
      dataIndex: 'createdAt',
      align: 'center',
    },
    {
      title: '操作',
      width: 200,
      align: 'center',
      dataIndex: 'operations',
      render: () => (
        <div className={styles.operations}>
          <Button type="text" size="small">
            回复
          </Button>
          <Popconfirm title="确认删除？">
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
        <div className={styles.toolbar}>
          <Form labelAlign="left" labelCol={{ span: 4 }} wrapperCol={{ span: 20 }}>
            <Row>
              <Col span={8} style={{ paddingRight: '20px' }}>
                <Form.Item label="昵称">
                  <Input />
                </Form.Item>
              </Col>
              <Col span={8} style={{ paddingRight: '20px' }}>
                <Form.Item label="状态">
                  <Input />
                </Form.Item>
              </Col>
              <Col span={8} style={{ paddingRight: '20px' }}>
                <Form.Item label="内容">
                  <Input />
                </Form.Item>
              </Col>
            </Row>
          </Form>
          <div
            style={{
              display: 'flex',
              flexDirection: 'column',
              justifyContent: 'space-between',
              paddingLeft: '20px',
              marginBottom: '20px',
              borderLeft: '1px solid var(--color-border-2)',
              boxSizing: 'border-box',
            }}
          >
            <Button type="primary" icon={<IconSearch />}>
              查询
            </Button>
          </div>
        </div>
        <Table rowKey="id" columns={columns} border={false} />
      </Card>
    </div>
  );
}
