// pages/Home.jsx
import { useState } from "react";
import { useAuth } from "../../contexts/authContext";
import { apiFetch } from "../../lib/apiFetch";

const Home = () => {
  const { currentUser } = useAuth();
  const [resp, setResp] = useState(null);
  const [err, setErr] = useState("");
  const apiGatewayURL = "localhost:3000";
  const sendRequest = async () => {
    setErr("");
    setResp(null);
    try {
      // hit your gateway or backend route
      const res = await apiFetch(`http://${apiGatewayURL}/health`, {
        method: "GET",
      });
      if (!res.ok) {
        const text = await res.text();
        throw new Error(`${res.status}: ${text}`);
      }
      const data = await res.json();
      setResp(data);
    } catch (e) {
      setErr(e.message);
    }
  };

  return (
    <div className="pt-14">
      <div className="text-2xl font-bold">
        Hello{" "}
        {currentUser.email},
        you are now logged in.
      </div>

      <button
        onClick={sendRequest}
        className="mt-6 px-4 py-2 rounded bg-indigo-600 text-white hover:bg-indigo-700"
      >
        Send req to api gateway
      </button>

      {resp && (
        <pre className="mt-4 p-3 border rounded bg-gray-50 text-sm">
          {JSON.stringify(resp, null, 2)}
        </pre>
      )}
      {err && <div className="mt-4 text-red-600 font-semibold">{err}</div>}
    </div>
  );
};

export default Home;
