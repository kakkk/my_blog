import React from 'react';
import { Comment } from '@arco-design/web-react';
import {IconMessage } from '@arco-design/web-react/icon';
import './App.css'
import CommentReply from './components/CommentReply';
import CommentItem from './components/CommonItem';



export default function App () {
  const actions = (
    <span className='custom-comment-action'>
      <IconMessage/> Reply
    </span>
  );
  return (
    <div>
      <CommentReply/>
      <Comment
        actions={actions}
        author={'kakkk'}
        avatar='//dn-qiniu-avatar.qbox.me/avatar/381206956e1d704103be7530dadd2e90'
        content={<div>Comment body content.</div>}
        datetime='1 hour'
      >
        <CommentItem/>
      </Comment>
      <CommentItem/>
    </div>
  )
}
