import React, { useState } from 'react';
import { Button, Comment, Grid, Input, Message } from '@arco-design/web-react';
import { commentArticle, getCommentList, replyComment } from '../../api';
import { useComments } from '../../CommentsContext';

const Row = Grid.Row;
const Col = Grid.Col;
const TextArea = Input.TextArea;

export default function CommentReply (props: {
  comment_id?: string;
  onClose?: any;
}) {
  const [nickname, setNickname] = useState("")
  const [nicknameStatus, setNicknameStatus] = useState<'error' | 'warning' | undefined>(undefined)
  const [email, setEmail] = useState("")
  const [emailStatus, setEmailStatus] = useState<'error' | 'warning' | undefined>(undefined)
  const [site, setSite] = useState("")
  const [content, setContent] = useState("")
  const [contentStatus, setContentStatus] = useState<'error' | 'warning' | undefined>(undefined)
  const { state, dispatch } = useComments();

  function clear () {
    setNickname('')
    setEmail('')
    setSite('')
    setContent('')
    if (props.comment_id && props.onClose) {
      props.onClose()
    }
  }

  function isEmail (email: string): boolean {
    const regex: RegExp = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return regex.test(email);
  }

  async function reply () {
    let paramOK = true;
    if (nickname === '') {
      setNicknameStatus('error')
      paramOK = false
    }
    if ((email === '') || !isEmail(email)) {
      setEmailStatus('error')
      paramOK = false
    }
    if (content === '') {
      setContentStatus('error')
      paramOK = false
    }
    if (!paramOK) {
      return
    }
    if (props.comment_id) {
      try {
        const resp = await replyComment(props.comment_id, state.article_id, nickname, email, site, content)
        if (resp.data.code === 0) {
          if (resp.data.data.comment_status === 2) {
            Message.success('success')
            clear()
            dispatch({ type: 'SET_COMMENTS', payload: resp.data.data.comments });
            return
          }
        }
        Message.error('fail')
      } catch (error) {
        Message.error('fail')
        console.error('Failed to reply:', error);
      }
      return
    }
    try {
      const resp = await commentArticle(state.article_id, nickname, email, site, content)
      if (resp.data.code === 0) {
        if (resp.data.data.comment_status === 2) {
          Message.success('success')
          clear()
          dispatch({ type: 'SET_COMMENTS', payload: resp.data.data.comments });
          return
        }
      }
      Message.error('fail')
    } catch (error) {
      console.error('Failed to reply:', error);
      Message.error('fail')
    }
  }

  return (<Comment
    align='right'
    content={
      <div>
        <Row>
          <Col sm={8} xs={24}>
            <Input addBefore={<><span style={{ color: 'red' }}>*</span>昵称</>}
                   status={nicknameStatus}
                   className="info-input"
                   onChange={(val) => {
                     if (nicknameStatus !== undefined) {
                       setNicknameStatus(undefined)
                     }
                     setNickname(val)
                   }}
                   value={nickname}
            />
          </Col>
          <Col sm={8} xs={24}>
            <Input addBefore={<><span style={{ color: 'red' }}>*</span>邮箱</>}
                   status={emailStatus}
                   className="info-input"
                   onChange={(val) => {
                     if (emailStatus !== undefined) {
                       setEmailStatus(undefined)
                     }
                     setEmail(val)
                   }}
                   value={email}
            />
          </Col>
          <Col sm={8} xs={24}>
            <Input addBefore='网址'
                   className="info-input"
                   onChange={setSite}
                   value={site}
            />
          </Col>
        </Row>
        <div className="info-input">
          <TextArea placeholder='Here is you content.'
                    style={{ minHeight: 128 }}
                    onChange={(val) => {
                      if (emailStatus !== undefined) {
                        setContentStatus(undefined)
                      }
                      setContent(val)
                    }}
                    value={content}
                    status={contentStatus}
          />
        </div>
        <div className="info-input">
          <div style={{ display: 'flex', justifyContent: 'right' }}>
            < div style={{ paddingRight: 5 }}>
              <Button key='0' type='secondary' onClick={clear}>
                Cancel
              </Button>
            </div>
            < div style={{ paddingLeft: 5 }}>
              <Button key='1' type='primary' onClick={reply}>
                Reply
              </Button>
            </div>
          </div>
        </div>
      </div>
    }
  />)
}
