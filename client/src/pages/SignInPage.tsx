import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { v4 as uuidv4 } from "uuid";

type User = {
  id: string;
  name: string;
};

function SignInPage() {
  const [user, setUser] = useState<User>({
    id: "",
    name: "",
  });

  const navigate = useNavigate();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    navigate("/chat", { state: { user } });
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-r from-indigo-500 to-purple-600 p-4">
      <form
        onSubmit={handleSubmit}
        className="bg-white rounded-lg shadow-lg p-8 w-full max-w-md"
      >
        <h1 className="text-2xl font-semibold mb-6 text-center text-gray-800">
          Sign In
        </h1>
        <input
          type="text"
          placeholder="Your name"
          value={user.name}
          onChange={(e) =>
            setUser({ ...user, name: e.target.value, id: uuidv4() })
          }
          className="w-full px-4 py-3 mb-6 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-400"
          required
        />
        <button
          type="submit"
          className="w-full bg-indigo-600 text-white py-3 rounded-md font-semibold hover:bg-indigo-700 transition"
        >
          Join Chat
        </button>
      </form>
    </div>
  );
}

export default SignInPage;
