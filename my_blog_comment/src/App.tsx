import React from 'react';
import './App.css'
import CommentReply from './components/CommentReply';
import CommentList from './components/CommentList';
import { CommentsProvider } from './CommentsContext';
import './Reset.css'


export default function App () {

  return (
    <CommentsProvider>
      <div className="App">
        <CommentReply/>
        <CommentList/>
      </div>
    </CommentsProvider>
  )
}
