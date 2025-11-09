'use client';

import {useState} from 'react'
import ChatRoomList from '@/components/ChatRoomList'

function page() {
  const [roomId, setRoomId] = useState<string | null>(null)
  
  return (
    <div className='w-full flex'>
      <ChatRoomList onRoomSelect={setRoomId} />
      <p>Selected Room ID: {roomId}</p>
    </div>
  )
}

export default page
