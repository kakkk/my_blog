import React, { useEffect } from 'react';
import { Table, Button, Breadcrumb, Card, Message, Popconfirm } from '@arco-design/web-react';
import { useSelector, useDispatch } from 'react-redux';
import { ColumnProps } from '@arco-design/web-react/es/Table/interface';
import { ReducerState } from '../../redux';
import styles from './style/index.module.less';
import {
  TOGGLE_VISIBLE,
  UPDATE_EDIT_MODAL_CONTENT,
  UPDATE_LIST,
  UPDATE_LOADING,
} from './redux/actionTypes';
import { deleteCategory, getCategoryById, getCategoryList } from '../../api/category';
import getUrlParams from '../../utils/getUrlParams';
import history from '../../history';
import EditCategory from './edit';

function Category(props) {
  const categoryState = useSelector((state: ReducerState) => state.category);

  const { data, loading } = categoryState;

  const dispatch = useDispatch();

  useEffect(() => {
    fetchData();
  }, [props.location]);
  const onDelete = async (row) => {
    try {
      const res: any = await deleteCategory(row.id);
      if (res.code === 0) {
        await fetchData();
        Message.success('删除成功！');
      } else {
        Message.error(res.msg);
      }
    } catch (error) {
    } finally {
    }
  };
  const onEdit = async (row) => {
    try {
      const res: any = await getCategoryById(row.id);
      if (res.code === 0) {
        dispatch({ type: UPDATE_EDIT_MODAL_CONTENT, payload: { content: res.data } });
        dispatch({ type: TOGGLE_VISIBLE, payload: { visible: true } });
      } else {
        Message.error(res.msg);
        dispatch({ type: TOGGLE_VISIBLE, payload: { visible: false } });
      }
    } catch (error) {
    } finally {
    }
  };

  async function fetchData() {
    const urlParams = getUrlParams();
    dispatch({ type: UPDATE_LOADING, payload: { loading: true } });
    try {
      let res: any;
      if (urlParams.parent) {
        res = await getCategoryList(urlParams.parent);
      } else {
        res = await getCategoryList(0);
      }
      if (res.code === 0) {
        if (res.data === null) {
          dispatch({ type: UPDATE_LIST, payload: { data: [] } });
        } else {
          dispatch({ type: UPDATE_LIST, payload: { data: res.data } });
        }
        dispatch({ type: UPDATE_LOADING, payload: { loading: false } });
      } else {
        Message.error(res.msg);
      }
    } catch (error) {
    } finally {
      dispatch({ type: UPDATE_LOADING, payload: { loading: false } });
    }
  }

  const onAdd = () => {
    dispatch({ type: TOGGLE_VISIBLE, payload: { visible: true } });
  };

  const columns: ColumnProps[] = [
    {
      title: '分类名称',
      dataIndex: 'name',
    },
    {
      title: '缩略名',
      dataIndex: 'slug',
    },
    {
      title: '子分类',
      dataIndex: 'children',
      width: 200,
      align: 'center',
      render: (_, record) => (
        <div>
          <Button
            type="text"
            size="small"
            onClick={() => {
              history.push(`/category?parent=${record.id}`);
            }}
          >
            {record.children === 0 ? '新增' : record.children}
          </Button>
        </div>
      ),
    },
    {
      title: '文章数量',
      dataIndex: 'count',
      width: 200,
      align: 'center',
    },
    {
      title: '操作',
      width: 300,
      align: 'center',
      dataIndex: 'operations',
      render: (_, record) => (
        <div className={styles.operations}>
          <Button
            type="text"
            size="small"
            onClick={() => {
              onEdit(record);
            }}
          >
            修改
          </Button>
          <Popconfirm
            title="确认删除？"
            onOk={() => {
              onDelete(record);
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
        <Breadcrumb.Item>
          <a
            onClick={() => {
              history.push(`/category`);
              fetchData();
            }}
            href="#"
          >
            分类管理
          </a>
        </Breadcrumb.Item>
        {getUrlParams().parent && <Breadcrumb.Item>子分类</Breadcrumb.Item>}
      </Breadcrumb>
      <Card bordered={false}>
        <div className={styles.toolbar}>
          <div>
            <Button type="primary" onClick={onAdd}>
              添加分类
            </Button>
          </div>
        </div>
        <Table
          rowKey="id"
          loading={loading}
          columns={columns}
          data={data}
          border={false}
          pagination={false}
        />
        <EditCategory title="添加分类" />
      </Card>
    </div>
  );
}

export default Category;
