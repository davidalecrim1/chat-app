import { useEffect, useRef } from 'react';

type WebSocketHandleIncomingMessage = (message: string) => void;
type WebSocketSendMessage = (message: string) => void;

const useWebSocket = (url: string, onMessage: WebSocketHandleIncomingMessage): WebSocketSendMessage => {
  const socket = useRef<WebSocket | null>(null);

  useEffect(() => {
    socket.current = new WebSocket(url);

    socket.current.onopen = () => {
      console.log('WebSocket connected');
    };

    socket.current.onmessage = (event: MessageEvent) => {
      onMessage?.(event.data);
    };

    socket.current.onclose = () => {
      console.log('WebSocket disconnected');
    };

    socket.current.onerror = (error: Event) => {
      console.error('WebSocket error', error);
    };

    return () => {
      socket.current?.close();
    };
  }, [url, onMessage])

  const sendMessage = (message: string) => {
    if (socket.current?.readyState === WebSocket.OPEN) {
      socket.current.send(message);
    } else {
      console.warn("WebSocket not ready");
    }
  };

  return sendMessage;
}

export default useWebSocket;
export type { WebSocketHandleIncomingMessage, WebSocketSendMessage };
