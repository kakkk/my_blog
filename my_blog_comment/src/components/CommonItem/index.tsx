import React from 'react';
import { Comment } from '@arco-design/web-react';
import { IconMessage } from '@arco-design/web-react/icon';
import './index.css'
import CommentReply from '../CommentReply';

export default function CommentItem (props: {
  id: string;
  nickname: string;
  avatar: string;
  website?: string;
  content: string;
  comment_at: string;
  reply_user?: string;
  children?: React.ReactNode;
}) {
  const [showReply, setShowReply] = React.useState(false);
  const nameLink = props.website ? (
    <a href={props.website} target="_blank" rel="noopener noreferrer">{props.nickname}</a>
  ) : (
    <a href="#" onClick={e => e.preventDefault()}>{props.nickname}</a>
  );
  // 当reply_user存在时，追加前缀
  const atUser = props.reply_user ? (
    <span style={{ color: 'blue' }}>@{props.reply_user}</span>
  ) : null;
  const handleHideReply = () => {
    setShowReply(false);
  };
  return (
    <Comment
      actions={
        <button className='custom-comment-action' onClick={() => {
          setShowReply(!showReply)
        }}>
          <IconMessage/> Reply
        </button>
      }
      author={nameLink}
      avatar={props.avatar}
      content={
        <div>
          {atUser} {props.content}
        </div>
      }
      datetime={props.comment_at}
    >
      {showReply ? <CommentReply comment_id={props.id} onClose={handleHideReply}/> : <></>}
      {props.children}
    </Comment>
  )
}
