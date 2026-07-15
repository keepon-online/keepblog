import { ref, onMounted, onUnmounted } from "vue";
import { getToken } from "@/utils/auth";

type MessageHandler = (data: any) => void;

interface WebSocketMessage {
  type: string;
  payload: any;
  time: number;
}

class WebSocketClient {
  private socket: WebSocket | null = null;
  private url: string;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 3000;
  private handlers: Map<string, MessageHandler[]> = new Map();
  private isManualClose = false;

  constructor() {
    const baseUrl = import.meta.env.VITE_BASE_URL || "";
    const wsProtocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const wsHost = baseUrl.replace(/^https?:/, wsProtocol);
    this.url = `${wsHost}/ws`;
  }

  connect(): void {
    if (this.socket?.readyState === WebSocket.OPEN) {
      return;
    }

    const token = getToken()?.accessToken;
    if (!token) {
      console.warn("WebSocket: No token available");
      return;
    }

    this.isManualClose = false;
    this.socket = new WebSocket(`${this.url}?token=${token}`);

    this.socket.onopen = () => {
      console.log("WebSocket connected");
      this.reconnectAttempts = 0;
    };

    this.socket.onmessage = (event: MessageEvent) => {
      try {
        const message: WebSocketMessage = JSON.parse(event.data);
        this.dispatchMessage(message);
      } catch (error) {
        console.error("WebSocket: Failed to parse message", error);
      }
    };

    this.socket.onclose = () => {
      console.log("WebSocket disconnected");
      if (!this.isManualClose) {
        this.attemptReconnect();
      }
    };

    this.socket.onerror = (error: Event) => {
      console.error("WebSocket error:", error);
    };
  }

  disconnect(): void {
    this.isManualClose = true;
    if (this.socket) {
      this.socket.close();
      this.socket = null;
    }
  }

  private attemptReconnect(): void {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      console.log(
        `WebSocket: Reconnecting... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`
      );
      setTimeout(() => this.connect(), this.reconnectDelay);
    }
  }

  on(type: string, handler: MessageHandler): void {
    if (!this.handlers.has(type)) {
      this.handlers.set(type, []);
    }
    this.handlers.get(type)!.push(handler);
  }

  off(type: string, handler?: MessageHandler): void {
    if (!handler) {
      this.handlers.delete(type);
    } else {
      const handlers = this.handlers.get(type);
      if (handlers) {
        const index = handlers.indexOf(handler);
        if (index > -1) {
          handlers.splice(index, 1);
        }
      }
    }
  }

  private dispatchMessage(message: WebSocketMessage): void {
    const handlers = this.handlers.get(message.type);
    if (handlers) {
      handlers.forEach(handler => handler(message.payload));
    }

    // 也分发给 "*" 处理器（监听所有消息）
    const allHandlers = this.handlers.get("*");
    if (allHandlers) {
      allHandlers.forEach(handler => handler(message));
    }
  }

  send(type: string, payload: any): void {
    if (this.socket?.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify({ type, payload }));
    }
  }

  get isConnected(): boolean {
    return this.socket?.readyState === WebSocket.OPEN;
  }
}

// 单例
let wsClient: WebSocketClient | null = null;

export function getWebSocketClient(): WebSocketClient {
  if (!wsClient) {
    wsClient = new WebSocketClient();
  }
  return wsClient;
}

// Composable Hook
export function useWebSocket() {
  const isConnected = ref(false);

  const connect = () => {
    const client = getWebSocketClient();
    client.connect();
    isConnected.value = client.isConnected;
  };

  const disconnect = () => {
    const client = getWebSocketClient();
    client.disconnect();
    isConnected.value = false;
  };

  const on = (type: string, handler: MessageHandler) => {
    const client = getWebSocketClient();
    client.on(type, handler);
  };

  const off = (type: string, handler?: MessageHandler) => {
    const client = getWebSocketClient();
    client.off(type, handler);
  };

  const send = (type: string, payload: any) => {
    const client = getWebSocketClient();
    client.send(type, payload);
  };

  return {
    isConnected,
    connect,
    disconnect,
    on,
    off,
    send
  };
}
