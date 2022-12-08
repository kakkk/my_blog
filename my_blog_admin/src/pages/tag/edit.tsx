import React, { useEffect } from 'react';
import { Form, Input, Message, Modal } from '@arco-design/web-react';
import { useDispatch, useSelector } from 'react-redux';
import { ReducerState } from '../../redux';
import {
  ADD_LIST_ITEM,
  CLEAR_EDIT_MODAL_CONTENT,
  TOGGLE_CONFIRM_LOADING,
  TOGGLE_VISIBLE,
  UPDATE_LIST,
} from './redux/actionTypes';
import { createTag, updateTag as apiUpdateTag } from '../../api/tag';

const FormItem = Form.Item;

const formItemLayout = {
  labelCol: {
    span: 6,
  },
  wrapperCol: {
    span: 18,
  },
};

export default function EditTag(props) {
  const { title } = props;
  const tagState = useSelector((state: ReducerState) => state.tag);
  const dispatch = useDispatch();

  const [form] = Form.useForm();

  const { visible, editModalContent, confirmLoading, data: listData } = tagState;

  const onSetField = () => {
    form.setFieldValue('name', editModalContent.name);
  };

  useEffect(onSetField, [visible]);
  const onCancel = () => {
    dispatch({ type: TOGGLE_VISIBLE, payload: { visible: false } });
    dispatch({ type: CLEAR_EDIT_MODAL_CONTENT });
    form.resetFields();
  };

  const addTag = async (data) => {
    try {
      const res: any = await createTag(data);
      if (res.code === 0) {
        Message.success('添加成功！');
        dispatch({ type: ADD_LIST_ITEM, payload: { item: res.data } });
      } else {
        Message.error(res.msg);
      }
    } finally {
    }
  };

  const updateTag = async (id: number, name: string) => {
    try {
      const res: any = await apiUpdateTag(id, name);
      if (res.code === 0) {
        Message.success('修改成功！');
        const ret = listData;
        for (let i = 0; i < ret.length; i++) {
          if (ret[i].id === id) {
            ret[i] = {
              id,
              name,
              count: ret[i].count,
            };
            break;
          }
        }
        dispatch({ type: UPDATE_LIST, payload: { data: ret } });
      } else {
        Message.error(res.msg);
      }
    } finally {
    }
  };
  const onOk = async () => {
    await form.validate();
    const data = form.getFields();
    dispatch({ type: TOGGLE_CONFIRM_LOADING, payload: { confirmLoading: true } });
    if (editModalContent.id === 0) {
      await addTag(data.name);
    } else {
      await updateTag(editModalContent.id, data.name);
    }
    dispatch({ type: TOGGLE_CONFIRM_LOADING, payload: { confirmLoading: false } });
    onCancel();
  };

  return (
    <Modal
      title={<div style={{ textAlign: 'left' }}>{title}</div>}
      visible={visible}
      onOk={onOk}
      confirmLoading={confirmLoading}
      onCancel={onCancel}
    >
      <Form {...formItemLayout} form={form}>
        <FormItem
          label="标签名称："
          field="name"
          rules={[{ required: true, message: '请输入标签名称' }]}
        >
          <Input placeholder="请输入标签名称" />
        </FormItem>
      </Form>
    </Modal>
  );
}
