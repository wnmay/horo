Example how to use ws hook (by chat)
1. Handle sending msg
```ts
  const { connected, messages, send } = useWebSocket();
  const [input, setInput] = useState("");

  // Handle sending a message
  const handleSend = () => {
    if (!input.trim()) return;
    send({
      type: "text",
      roomId: "room-123",
      senderId: "user-999",
      content: input,
    });
    setInput("");
  };
```

2. Handle receiving new msg/notificaiton
```ts
  // 💬 Handle receiving messages or notifications
  useEffect(() => {
    if (messages.length === 0) return;
    const latest = messages[messages.length - 1];

    switch (latest.type) {
      // ----- Text messages -----
      case "text":
        console.log("💬 New chat message:", latest.content);
        break;

      // ----- Notifications -----
      case "notification":
        switch (msg.trigger) {
            case Trigger.OrderPaid:
                break;
            case Trigger.OrderCompleted:
                break;
            }

          default:
            console.warn("⚠️ Unknown notification:", latest);
        }
        break;
    }
  }, [messages]);

```