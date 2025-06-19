
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import './App.css'
import SignInPage from './pages/SignInPage';
import ChatPage from './pages/ChatPage';

function App() {
  // TODOS
  // Create the Login in Page
  // Create the Chat Page
  // Create the Chat Component
  // Create the messages within the Chat
  // See best practices for calling the websocket connection.


  return (
    <>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<SignInPage />} />
          <Route path="/chat" element={<ChatPage />} />
        </Routes>
      </BrowserRouter>
    </>
  )
}

export default App
