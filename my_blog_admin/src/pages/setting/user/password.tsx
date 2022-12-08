import React from 'react';
import { Button, Card, Form, Input, Message, Popconfirm } from '@arco-design/web-react';
import { useDispatch, useSelector } from 'react-redux';
import { updateUserPassword } from '../../../api/user';
import { ReducerState } from '../../../redux';
import { UPDATE_PASSWORD_RETYPE_OK } from './redux/actionType';

export default function Password() {
  const FormItem = Form.Item;
  const InputPassword = Input.Password;
  const [form] = Form.useForm();
  const userSettingState = useSelector((state: ReducerState) => state.userSetting);
  const dispatch = useDispatch();
  const { passwordRetypeOk } = userSettingState;
  const onUpdate = async () => {
    const data = form.getFields();
    if (data.newPassword !== data.retype) {
      Message.error('两次密码不一致！');
      return;
    }
    const close = Message.loading('loading...');
    try {
      const res: any = await updateUserPassword(data.oldPassword, data.newPassword);
      if (res.code === 0) {
        Message.success('修改成功！');
        form.resetFields();
      } else {
        Message.error(res.msg);
      }
    } finally {
    }
    close();
  };
  const onChange = () => {
    const data = form.getFields();
    if (data.newPassword !== data.retype) {
      dispatch({
        type: UPDATE_PASSWORD_RETYPE_OK,
        payload: { passwordRetypeOk: false },
      });
    } else {
      dispatch({
        type: UPDATE_PASSWORD_RETYPE_OK,
        payload: { passwordRetypeOk: true },
      });
    }
  };
  return (
    <Card title="修改密码" style={{ margin: '5px', padding: '0px 20px' }}>
      <Form
        style={{ marginTop: '6px' }}
        labelCol={{ flex: '80px' }}
        wrapperCol={{ flex: 'auto' }}
        form={form}
      >
        <FormItem label="旧密码" field="oldPassword">
          <InputPassword />
        </FormItem>
        <FormItem label="新密码" field="newPassword">
          <InputPassword onChange={onChange} />
        </FormItem>
        <FormItem label="重新输入" field="retype">
          <InputPassword onChange={onChange} error={!passwordRetypeOk} />
        </FormItem>
        <div style={{ textAlign: 'right' }}>
          <Popconfirm
            title="确认修改？"
            onOk={() => {
              onUpdate();
            }}
          >
            <Button type="primary">修改</Button>
          </Popconfirm>
        </div>
      </Form>
    </Card>
  );
}
