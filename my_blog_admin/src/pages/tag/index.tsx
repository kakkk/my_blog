import React, { useEffect } from 'react';
import {
  Table,
  Button,
  Input,
  Breadcrumb,
  Card,
  Tag as ArcoTag,
  Message,
  Popconfirm,
} from '@arco-design/web-react';
import { useSelector, useDispatch } from 'react-redux';
import { ColumnProps } from '@arco-design/web-react/es/Table/interface';
import { IconTag } from '@arco-design/web-react/icon';
import { ReducerState } from '../../redux';
import styles from './style/index.module.less';
import {
  UPDATE_LOADING,
  UPDATE_LIST,
  UPDATE_SEARCH_KEYWORD,
  TOGGLE_VISIBLE,
  UPDATE_EDIT_MODAL_CONTENT,
} from './redux/actionTypes';
import getUrlParams from '../../utils/getUrlParams';
import { deleteTag, getTagList } from '../../api/tag';
import history from '../../history';
import EditTag from './edit';

function Tag(props) {
  const dispatch = useDispatch();
  const tagState = useSelector((state: ReducerState) => state.tag);
  const { data, pagination, loading, searchKeyWord } = tagState;

  const onEdit = (row) => {
    dispatch({
      type: UPDATE_EDIT_MODAL_CONTENT,
      payload: { content: { id: row.id, name: row.name } },
    });
    dispatch({ type: TOGGLE_VISIBLE, payload: { visible: true } });
  };

  const onDelete = async (row) => {
    try {
      const res: any = await deleteTag(row.id);
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

  const columns: ColumnProps[] = [
    {
      title: '标签名称',
      width: 300,
      render: (_, record) => {
        return (
          <ArcoTag size="medium" icon={<IconTag />}>
            {record.name}
          </ArcoTag>
        );
      },
    },
    {
      title: '文章数量',
      dataIndex: 'count',
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

  async function fetchData() {
    dispatch({ type: UPDATE_LOADING, payload: { loading: true } });
    const urlParams = getUrlParams();
    const page = urlParams.page ? Number(urlParams.page) : 0;
    const size = urlParams.size ? Number(urlParams.size) : 10;
    const keyword = urlParams.keyword ? urlParams.keyword : '';
    try {
      const res: any = await getTagList(keyword, page, size);
      if (res.code === 0) {
        dispatch({
          type: UPDATE_LIST,
          payload: {
            data: res.data != null ? res.data.tag_list : [],
            pagination:
              res.data != null
                ? {
                    pageSize: res.data.pagination.limit,
                    current: res.data.pagination.page,
                    total: res.data.pagination.total,
                  }
                : {},
          },
        });
        dispatch({ type: UPDATE_SEARCH_KEYWORD, payload: { searchKeyWord: keyword } });
      } else {
        Message.error(res.msg);
      }
    } finally {
      dispatch({ type: UPDATE_LOADING, payload: { loading: false } });
    }
  }

  useEffect(() => {
    fetchData();
  }, [props.location]);

  const onChangeTable = (pagination) => {
    history.push(
      `/tag?page=${pagination.current}&size=${pagination.pageSize}${
        searchKeyWord === '' ? `&${searchKeyWord}` : ''
      }`
    );
  };

  const onSearch = (value) => {
    history.push(`/tag?page=1&size=10${value !== '' ? `&keyword=${value}` : ''}`);
  };

  const onAdd = () => {
    dispatch({ type: TOGGLE_VISIBLE, payload: { visible: true } });
  };

  return (
    <div className={styles.container}>
      <Breadcrumb style={{ marginBottom: 20 }}>
        <Breadcrumb.Item>
          <a
            onClick={() => {
              history.push(`/tag`);
              fetchData();
            }}
            href="#"
          >
            标签管理
          </a>
        </Breadcrumb.Item>
      </Breadcrumb>
      <Card bordered={false}>
        <div className={styles.toolbar}>
          <div>
            <Button type="primary" onClick={onAdd}>
              添加标签
            </Button>
          </div>
          <div>
            <Input.Search
              style={{ width: 300 }}
              searchButton
              placeholder="搜索"
              onSearch={onSearch}
            />
          </div>
        </div>
        <Table
          rowKey="id"
          loading={loading}
          onChange={onChangeTable}
          pagination={pagination}
          columns={columns}
          data={data}
          border={false}
        />
        <EditTag title="添加标签" />
      </Card>
    </div>
  );
}

export default Tag;
