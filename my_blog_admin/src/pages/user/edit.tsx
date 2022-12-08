import {
  Button,
  Form,
  Input,
  Message,
  Modal,
  Popconfirm,
  Select,
  Typography,
} from '@arco-design/web-react';
import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { ReducerState } from '../../redux';
import {
  ADD_LIST_ITEM,
  TOGGLE_PASSWORD_VISIBLE,
  TOGGLE_VISIBLE,
  UPDATE_LIST_ITEM,
} from './redux/actionTypes';
import { createUser, updateUserInfoById } from '../../api/user';

export default function EditUser(props) {
  const userState = useSelector((state: ReducerState) => state.user);
  const dispatch = useDispatch();
  const { title, userInfo } = props;
  const [usernameDisable, setUsernameDisable] = React.useState(false);
  const [passwordContent, setPasswordContent] = React.useState('');
  const [passwordVisible, setPasswordVisible] = React.useState(false);
  const { visible, confirmLoading } = userState;
  const [form] = Form.useForm();
  const FormItem = Form.Item;
  useEffect(() => {
    if (userInfo.id === 0) {
      form.clearFields();
    } else {
      form.setFieldValue('username', userInfo.userName);
      form.setFieldValue('email', userInfo.email);
      form.setFieldValue('role', userInfo.role);
      form.setFieldValue('nickname', userInfo.nickname);
    }
    if (title === '修改用户') {
      setUsernameDisable(true);
    } else {
      setUsernameDisable(false);
    }
  }, [userInfo]);
  const addUser = async () => {
    const username = form.getFieldValue('username');
    const nickname = form.getFieldValue('nickname');
    const email = form.getFieldValue('email');
    const role = form.getFieldValue('role');
    try {
      const res: any = await createUser({ username, nickname, email, role });
      if (res.code === 0) {
        dispatch({
          type: ADD_LIST_ITEM,
          payload: {
            item: {
              id: res.data.id,
              username: res.data.id,
              nickname: res.data.nickname,
              email: res.data.email,
              role: res.data.role,
              avatar: res.data.avatar,
              post_count: 0,
              page_count: 0,
            },
          },
        });
        setPasswordContent(res.data.password);
        setPasswordVisible(true);
      } else {
        Message.error(res.msg);
      }
    } finally {
      dispatch({ type: TOGGLE_VISIBLE, payload: { visible: false } });
    }
  };
  const updateUser = async () => {
    const nickname = form.getFieldValue('nickname');
    const email = form.getFieldValue('email');
    const role = form.getFieldValue('role');
    try {
      const res: any = await updateUserInfoById(userInfo.id, { nickname, email, role });
      if (res.code === 0) {
        dispatch({
          type: UPDATE_LIST_ITEM,
          payload: { userInfo: { id: userInfo.id, nickname, email, role } },
        });
        Message.success('修改成功');
      } else {
        Message.error(res.msg);
      }
    } finally {
      dispatch({ type: TOGGLE_VISIBLE, payload: { visible: false } });
    }
  };
  const onOK = async () => {
    if (title === '添加用户') {
      await addUser();
    } else {
      await updateUser();
    }
  };
  const formItemLayout = {
    labelCol: {
      span: 6,
    },
    wrapperCol: {
      span: 18,
    },
  };
  return (
    <Modal
      title={<div style={{ textAlign: 'left' }}>{title}</div>}
      visible={visible}
      onOk={() => {
        dispatch({ type: TOGGLE_VISIBLE, payload: { visible: false } });
        dispatch({ type: TOGGLE_PASSWORD_VISIBLE, payload: { passwordVisible: true } });
      }}
      confirmLoading={confirmLoading}
      onCancel={() => {
        dispatch({ type: TOGGLE_VISIBLE, payload: { visible: false } });
      }}
      footer={
        <>
          <Button
            onClick={() => {
              dispatch({ type: TOGGLE_VISIBLE, payload: { visible: false } });
            }}
          >
            取消
          </Button>
          <Popconfirm title={title === '修改用户' ? '确认修改？' : '确认添加'} onOk={onOK}>
            <Button type="primary">确定</Button>
          </Popconfirm>
        </>
      }
    >
      <Form {...formItemLayout} form={form}>
        <FormItem
          label="用户名："
          field="username"
          rules={[{ required: true, message: '请输入用户名' }]}
        >
          <Input placeholder="用户名" disabled={usernameDisable} />
        </FormItem>
        <FormItem
          label="昵称："
          field="nickname"
          rules={[{ required: true, message: '请输入昵称' }]}
        >
          <Input placeholder="昵称" />
        </FormItem>
        <FormItem label="邮箱：" field="email" rules={[{ required: true, message: '请输入邮箱' }]}>
          <Input placeholder="邮箱" />
        </FormItem>
        <FormItem label="角色：" field="role" rules={[{ required: true, message: '请选择角色' }]}>
          <Select placeholder="请选择...">
            <Select.Option value="EDITOR">作者</Select.Option>
            <Select.Option value="ADMINISTRATOR">管理员</Select.Option>
          </Select>
        </FormItem>
      </Form>
      <Modal
        title={<div style={{ textAlign: 'left' }}>初始密码</div>}
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
    </Modal>
  );
}
