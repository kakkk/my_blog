import React, { useEffect } from 'react';
import {
  Avatar,
  Breadcrumb,
  Button,
  Card,
  Input,
  Message,
  Modal,
  Popconfirm,
  Select,
  Table,
  Tag as ArcoTag,
  Typography,
} from '@arco-design/web-react';
import { ColumnProps } from '@arco-design/web-react/es/Table/interface';
import { IconUser } from '@arco-design/web-react/icon';
import { useSelector, useDispatch } from 'react-redux';
import styles from './style/index.module.less';
import history from '../../history';
import { ReducerState } from '../../redux';
import { UPDATE_LOADING, UPDATE_LIST, TOGGLE_VISIBLE } from './redux/actionTypes';
import getUrlParams from '../../utils/getUrlParams';
import { getUserList, resetUserPasswordById } from '../../api/user';
import EditUser from './edit';

export default function User(props) {
  const dispatch = useDispatch();
  const userState = useSelector((state: ReducerState) => state.user);
  const [editTitle, setEditTitle] = React.useState('添加用户');
  const [passwordContent, setPasswordContent] = React.useState('');
  const [passwordVisible, setPasswordVisible] = React.useState(false);
  const [userInfo, setUserInfo] = React.useState({
    id: 0,
    userName: '',
    email: '',
    role: '',
    nickname: '',
  });
  const { data, loading, pagination } = userState;
  const fetchData = async () => {
    dispatch({ type: UPDATE_LOADING, payload: { loading: true } });
    const urlParams = getUrlParams();
    const page = urlParams.page ? Number(urlParams.page) : 0;
    const size = urlParams.size ? Number(urlParams.size) : 10;
    const username = urlParams.username ? urlParams.username : '';
    const nickname = urlParams.nickname ? urlParams.nickname : '';
    const email = urlParams.email ? urlParams.email : '';
    try {
      const res: any = await getUserList(username, nickname, email, page, size);
      if (res.code === 0) {
        dispatch({
          type: UPDATE_LIST,
          payload: {
            data: res.data != null ? res.data : [],
            pagination: {
              pageSize: res.pagination.limit,
              current: res.pagination.page,
              total: res.pagination.total,
            },
          },
        });
      } else {
        Message.error(res.msg);
      }
    } finally {
      dispatch({ type: UPDATE_LOADING, payload: { loading: false } });
    }
  };

  useEffect(() => {
    fetchData();
  }, [props.location]);

  let searchSelectValue = 'username';

  const onSearch = (value: string) => {
    switch (searchSelectValue) {
      case 'username':
        history.push(`/user?page=1&size=10${value !== '' ? `&username=${value}` : ''}`);
        break;
      case 'nickname':
        history.push(`/user?page=1&size=10${value !== '' ? `&nickname=${value}` : ''}`);
        break;
      case 'email':
        history.push(`/user?page=1&size=10${value !== '' ? `&email=${value}` : ''}`);
        break;
      default:
        history.push(`/user?page=1&size=10`);
        break;
    }
  };

  const onReset = async (id: number) => {
    try {
      const res: any = await resetUserPasswordById(id);
      if (res.code === 0) {
        setPasswordContent(res.data.password);
        setPasswordVisible(true);
      } else {
        Message.error(res.msg);
      }
    } finally {
    }
  };

  const columns: ColumnProps[] = [
    {
      title: '用户',
      render: (_, record) => {
        return (
          <div>
            <Avatar>
              <img alt="avatar" src={record.avatar} />
            </Avatar>
            <a style={{ marginLeft: '10px' }}>{record.nickname}</a>
          </div>
        );
      },
    },
    {
      title: '用户名',
      dataIndex: 'username',
      align: 'center',
    },
    {
      title: '身份',
      render: (_, record) => {
        return (
          <ArcoTag size="medium" icon={<IconUser />}>
            {record.role === 'ADMINISTRATOR'
              ? '管理员'
              : record.role === 'EDITOR'
              ? '作者'
              : '外星人'}
          </ArcoTag>
        );
      },
      align: 'center',
    },
    {
      title: '文章',
      dataIndex: 'post_count',
      align: 'center',
    },
    {
      title: '页面',
      dataIndex: 'page_count',
      align: 'center',
    },
    {
      title: '操作',
      align: 'center',
      width: 250,
      render: (_, record) => (
        <div className={styles.operations}>
          <Button
            type="text"
            size="small"
            onClick={() => {
              setEditTitle('修改用户');
              dispatch({ type: TOGGLE_VISIBLE, payload: { visible: true } });
              setUserInfo({
                id: record.id,
                userName: record.username,
                email: record.email,
                role: record.role,
                nickname: record.nickname,
              });
            }}
          >
            修改
          </Button>
          <Popconfirm
            title="重置密码？"
            onOk={async () => {
              await onReset(record.id);
            }}
          >
            <Button type="text" size="small">
              重置密码
            </Button>
          </Popconfirm>
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
              history.push(`/user`);
            }}
            href="#"
          >
            用户管理
          </a>
        </Breadcrumb.Item>
      </Breadcrumb>
      <Card bordered={false}>
        <div className={styles.toolbar}>
          <div>
            <Button
              type="primary"
              onClick={() => {
                setEditTitle('添加用户');
                setUserInfo({ id: 0, userName: '', email: '', role: '', nickname: '' });
                dispatch({ type: TOGGLE_VISIBLE, payload: { visible: true } });
              }}
            >
              添加用户
            </Button>
          </div>
          <EditUser title={editTitle} userInfo={userInfo} />
          <div>
            <Input.Search
              style={{ width: 320 }}
              searchButton
              placeholder="搜索"
              addBefore={
                <Select
                  size="small"
                  style={{ width: 90 }}
                  defaultValue="username"
                  bordered={false}
                  onChange={(value) => {
                    searchSelectValue = value;
                  }}
                >
                  <Select.Option value="username">用户名</Select.Option>
                  <Select.Option value="nickname">昵称</Select.Option>
                  <Select.Option value="email">邮箱</Select.Option>
                </Select>
              }
              onSearch={onSearch}
            />
          </div>
        </div>
        <Table
          rowKey="id"
          loading={loading}
          pagination={pagination}
          columns={columns}
          data={data}
          border={false}
        />
      </Card>
      <Modal
        title={<div style={{ textAlign: 'left' }}>重置密码</div>}
        visible={passwordVisible}
        onOk={() => {
          setPasswordVisible(false);
        }}
        onCancel={() => {
          setPasswordVisible(false);
        }}
        hideCancel
      >
        <Typography>
          <Typography.Paragraph>该密码只会出现一次，请妥善保存:</Typography.Paragraph>
          <Typography.Paragraph copyable style={{ fontSize: 'large' }}>
            {passwordContent}
          </Typography.Paragraph>
        </Typography>
      </Modal>
    </div>
  );
}
