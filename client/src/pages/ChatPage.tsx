import { useState, useCallback } from "react";
import useWebSocket from "../hooks/useWebSocket";
import { useLocation } from "react-router-dom";
import type { User } from './../types/User';
import type { ChatMessage } from './../types/chat';

function ChatPage() {
  const [message, setMessage] = useState('');
  const [chatMessages, setChatMessages] = useState<ChatMessage[]>([]);

  const location = useLocation();
  const user: User = location.state?.user;

  const handleIncomingMessage = useCallback((data: string) => {
    console.log('Received:', data);
    const parsedData = JSON.parse(data);

    if (parsedData.type === 'chat') {
      setChatMessages(prev => [...prev, parsedData]);
    } else {
      console.log('Unknown message type:', parsedData);
    }
  }, []);

  const sendMessage = useWebSocket(
    `ws://localhost:8201/ws/connect?id=${user.id}&name=${user.name}`,
    handleIncomingMessage
  );

  const sendChatMessage = () => {
    const trimmed = message.trim();
    if (!trimmed) return;

    sendMessage(JSON.stringify({
      type: "chat",
      payload: {
        user: { id: user.id, name: user.name },
        message: trimmed
      }
    }));

    setMessage('');
  };

  const handleSend = (e: React.FormEvent | React.MouseEvent) => {
    e.preventDefault?.();
    sendChatMessage();
  };

  const handleKeyPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      sendChatMessage();
    }
  };

  if (!user) {
    return <div>No user data found. Please sign in first.</div>;
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100 p-4">
      <div className="flex flex-col bg-white rounded-lg shadow-lg w-full max-w-md max-h-[80vh]">
        <header className="bg-blue-600 text-white p-4 rounded-t-lg">
          <h1 className="text-xl font-semibold">Chat</h1>
          <p className="text-sm opacity-80">Welcome, {user.name}!</p>
        </header>

        <main className="flex-1 overflow-y-auto p-4 space-y-4 min-h-[300px]">
          {chatMessages.map((msg, idx) => (
            <div
              key={idx}
              className={`flex ${msg.payload.user.id === user.id ? 'justify-end' : 'justify-start'}`}
            >
              <div
                className={`max-w-xs p-3 rounded-lg text-white ${msg.payload.user.id === user.id ? 'bg-blue-500' : 'bg-gray-500'
                  }`}
              >
                <p className="font-semibold">{msg.payload.user.name}</p>
                <p>{msg.payload.message}</p>
              </div>
            </div>
          ))}
        </main>

        <footer className="p-4 border-t flex">
          <form onSubmit={handleSend} className="flex w-full space-x-2">
            <input
              type="text"
              value={message}
              onChange={e => setMessage(e.target.value)}
              onKeyDown={handleKeyPress}
              placeholder="Type a message..."
              className="flex-1 p-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <button
              type="submit"
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              Send
            </button>
          </form>
        </footer>
      </div>
    </div>
  );
}

export default ChatPage;
