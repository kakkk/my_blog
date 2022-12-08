import React, { useEffect } from 'react';
import { Button, Card, Form, Input, Message, Popconfirm } from '@arco-design/web-react';
import { useDispatch, useSelector } from 'react-redux';
import { ReducerState } from '../../../redux';
import { UPDATE_USER_INFO } from '../../../redux/actionTypes';
import { updateUserInfo as apiUpdateInfo } from '../../../api/user';

export default function BasicInfo() {
  const FormItem = Form.Item;
  const InputTextArea = Input.TextArea;
  const globalState = useSelector((state: ReducerState) => state.global);
  const dispatch = useDispatch();
  const { userInfo } = globalState;
  const [form] = Form.useForm();
  const onReset = () => {
    form.setFieldValue('nickname', userInfo.nickname);
    form.setFieldValue('email', userInfo.email);
    form.setFieldValue('description', userInfo.description);
  };
  const updateUserInfo = async (info: any) => {
    try {
      const newInfo = {
        ...info,
        id: userInfo.id,
        username: userInfo.username,
        role: userInfo.role,
        avatar: userInfo.avatar,
      };
      const res: any = await apiUpdateInfo(newInfo);
      if (res.code === 0) {
        Message.success('修改成功！');
        dispatch({ type: UPDATE_USER_INFO, payload: { userInfo: newInfo } });
      } else {
        Message.error(res.msg);
      }
    } finally {
    }
  };
  const onUpdate = async () => {
    const close = Message.loading('loading...');
    const data = form.getFields();
    await updateUserInfo(data);
    close();
  };
  useEffect(onReset, []);
  return (
    <Card title="基本设置" style={{ margin: '5px', padding: '0px 15px' }}>
      <Form
        style={{ marginTop: '6px' }}
        labelCol={{ flex: '80px' }}
        wrapperCol={{ flex: 'auto' }}
        form={form}
      >
        <FormItem label="昵称" field="nickname">
          <Input placeholder="昵称" />
        </FormItem>
        <FormItem label="邮箱" field="email">
          <Input placeholder="邮箱" />
        </FormItem>
        <FormItem label="自我介绍" field="description">
          <InputTextArea style={{ height: '100px' }} placeholder="一小段自我介绍" />
        </FormItem>
      </Form>
      <div style={{ textAlign: 'right' }}>
        <Button type="primary" style={{ marginRight: '10px' }} onClick={onReset}>
          重置
        </Button>
        <Popconfirm
          title="确认修改？"
          onOk={() => {
            onUpdate();
          }}
        >
          <Button type="primary">修改</Button>
        </Popconfirm>
      </div>
    </Card>
  );
}
