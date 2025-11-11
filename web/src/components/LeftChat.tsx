'use client';

import React from 'react';

export interface ChatRoom {
  ID: string;
  ProphetID: string;
  CustomerID: string;
  CourseID: string;
  CreatedAt: string;
  LastMessage: string;
  IsDone: boolean;
  ProphetName?: string;
  CustomerName?: string;
  courseName?: string; 
}

interface LeftChatProps {
  room: ChatRoom;
  onClick?: () => void;
}

const LeftChat: React.FC<LeftChatProps> = ({ room, onClick }) => {
  const formatDate = (dateString: string): string => {
    const date = new Date(dateString);
    const months = ['JAN', 'FEB', 'MAR', 'APR', 'MAY', 'JUN', 'JUL', 'AUG', 'SEP', 'OCT', 'NOV', 'DEC'];
    const day = date.getDate();
    const month = months[date.getMonth()];
    return `${day} ${month}`;
  };

  return (
    <div 
      onClick={onClick} 
      className="p-4 border-b border-gray-200 cursor-pointer hover:bg-gray-50 transition-colors"
    >
      {/* Top section: Course name/ID and created date */}
      <div className="flex justify-between items-start mb-2">
        <div className="flex-1">
          <h3 className="font-semibold text-base truncate">
            {room.courseName || 'Loading...'}
          </h3>
          <p className="text-sm text-gray-500 truncate">
            {room.CourseID}
          </p>
        </div>
        <span className="text-xs text-gray-400 ml-2 whitespace-nowrap">
          {formatDate(room.CreatedAt)}
        </span>
      </div>

      {/* Middle section: Last message */}
      <div className="mb-2">
        <p className="text-sm text-gray-600 truncate">
          {room.LastMessage || 'No messages yet'}
        </p>
      </div>

      {/* Bottom section: Room status */}
      <div className="flex justify-end">
        <span
          className={`text-xs px-2 py-1 rounded-full ${
            room.IsDone
              ? 'bg-green-100 text-green-700'
              : 'bg-blue-100 text-blue-700'
          }`}
        >
          {room.IsDone ? 'Done' : 'Active'}
        </span>
      </div>
    </div>
  );
};

export default LeftChat;