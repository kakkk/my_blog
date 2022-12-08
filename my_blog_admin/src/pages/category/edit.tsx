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
import { createCategory, updateCategoryById } from '../../api/category';
import getUrlParams from '../../utils/getUrlParams';

const FormItem = Form.Item;

const formItemLayout = {
  labelCol: {
    span: 6,
  },
  wrapperCol: {
    span: 18,
  },
};

export default function EditCategory(props) {
  const { title } = props;
  const categoryState = useSelector((state: ReducerState) => state.category);
  const dispatch = useDispatch();

  const [form] = Form.useForm();

  const { editModalContent, visible, confirmLoading, data: listData } = categoryState;

  const onSetField = () => {
    form.setFieldValue('name', editModalContent.name);
    form.setFieldValue('slug', editModalContent.slug);
    form.setFieldValue('description', editModalContent.description);
  };

  useEffect(onSetField, [visible]);
  const onCancel = () => {
    dispatch({ type: TOGGLE_VISIBLE, payload: { visible: false } });
    dispatch({ type: CLEAR_EDIT_MODAL_CONTENT });
    form.resetFields();
  };

  const addCategory = async (data) => {
    const urlParams = getUrlParams();
    try {
      const res: any = await createCategory(data, urlParams.parent);
      if (res.code === 0) {
        Message.success('添加成功！');
        dispatch({ type: ADD_LIST_ITEM, payload: { item: res.data } });
      } else {
        Message.error(res.msg);
      }
    } finally {
    }
  };

  const updateCategory = async (id: number, data) => {
    const urlParams = getUrlParams();
    try {
      const res: any = await updateCategoryById(id, {
        ...data,
        parent: urlParams.parent,
      });
      if (res.code === 0) {
        Message.success('修改成功！');
        const ret = listData;
        for (let i = 0; i < ret.length; i++) {
          if (ret[i].id === id) {
            ret[i] = {
              id,
              ...data,
              count: ret[i].count,
              children: ret[i].children,
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
      await addCategory(data);
    } else {
      await updateCategory(editModalContent.id, data);
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
          label="分类名称："
          field="name"
          rules={[{ required: true, message: '请输入分类名称' }]}
        >
          <Input placeholder="请输入分类名称" />
        </FormItem>
        <FormItem
          label="缩略名："
          field="slug"
          rules={[{ required: true, message: '请输入缩略名' }]}
        >
          <Input placeholder="请输入缩略名" />
        </FormItem>
        <FormItem label="描述：" field="description">
          <Input placeholder="描述" />
        </FormItem>
      </Form>
    </Modal>
  );
}
