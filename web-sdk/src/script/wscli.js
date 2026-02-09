// wscli.js

/**
 * 消息类型协议（前后端约定）
 */
export const MSG_TYPES = {
  // 客户端 → 服务端
  REQ_MESSAGE: "message.req",
  TYPING_START: "typing.start",
  SESSION_CLOSE: "session.close",

  // 服务端 → 客户端
  RSP_MESSAGE: "message.rsp",
  SESSION_UPDATE: "session.update",
  TYPING_INDICATOR: "typing.start",
};

/**
 * 消息内容类型（media type）
 */
export const CONTENT_TYPES = {
  TEXT: "text", // 包含 Emoji、URL、换行符的纯文本
  IMAGE: "image", // payload.url 必填
  AUDIO: "audio", // payload.url + duration
  FILE: "file", // payload.url + name + size
};

/**
 * 会话状态
 */
export const SESSION_STATUS = {
  WAITING: "waiting",
  CHATTING: "chatting",
  CLOSED: "closed",
};

/**
 * 客服 WebSocket 客户端
 */
export class WSClient {
  constructor(appid, visitorId, options = {}) {
    this.appid = appid;
    this.visitorId = visitorId;
    this.wsUrl = options.wsUrl || "ws://127.0.0.1:5400/ws/chat";
    this.onMessage = options.onMessage || (() => {});
    this.onStatusChange = options.onStatusChange || (() => {});
    this.onError = options.onError || console.error;
    this.onConnected = options.onConnected || (() => {});

    this.ws = null;
    this.isConnected = false;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
  }

  connect() {
    if (this.ws?.readyState === WebSocket.OPEN) return;

    const url = `${this.wsUrl}?appid=${encodeURIComponent(this.appid)}&visitor_id=${encodeURIComponent(this.visitorId)}`;
    this.ws = new WebSocket(url);

    this.ws.onopen = () => {
      this.isConnected = true;
      this.reconnectAttempts = 0;
      this.onConnected();
      this.onStatusChange("connected");
    };

    this.ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data);
        this._handleIncoming(msg);
      } catch (e) {
        this.onError("invalid WS message:", event.data);
      }
    };

    this.ws.onclose = () => {
      this.isConnected = false;
      this.onStatusChange("disconnected");
      this._reconnect();
    };

    this.ws.onerror = (err) => {
      this.onError("WS error:", err);
      this.ws.close(); // 触发 onclose
    };
  }

  _handleIncoming(msg) {
    switch (msg.type) {
      case MSG_TYPES.RSP_MESSAGE:
        this.onMessage({
          type: "message",
          id: msg.payload.id,
          from: msg.payload.from,
          contentType: msg.payload.msg_type,
          content: msg.payload.content,
          url: msg.payload.url,
          name: msg.payload.name,
          size: msg.payload.size,
          duration: msg.payload.duration,
          timestamp: msg.payload.timestamp,
        });
        break;

      case MSG_TYPES.SESSION_UPDATE:
        this.onStatusChange("session", msg.payload);
        break;

      case MSG_TYPES.TYPING_INDICATOR:
        this.onMessage({ type: "typing", from: msg.payload.from });
        break;

      default:
      // 忽略未知类型（未来兼容）
    }
  }

  // --- 发送 API ---

  sendText(content) {
    if (!this.isConnected) return;
    this._send(MSG_TYPES.REQ_MESSAGE, {
      msg_type: CONTENT_TYPES.TEXT,
      content,
    });
  }

  sendImage(url, name = "image.jpg") {
    if (!this.isConnected) return;
    this._send(MSG_TYPES.REQ_MESSAGE, {
      msg_type: CONTENT_TYPES.IMAGE,
      url,
      name,
    });
  }

  sendAudio(url, duration) {
    if (!this.isConnected) return;
    this._send(MSG_TYPES.REQ_MESSAGE, {
      msg_type: CONTENT_TYPES.AUDIO,
      url,
      duration,
    });
  }

  sendFile(url, name, size) {
    if (!this.isConnected) return;
    this._send(MSG_TYPES.REQ_MESSAGE, {
      msg_type: CONTENT_TYPES.FILE,
      url,
      name,
      size,
    });
  }

  sendTyping() {
    if (!this.isConnected) return;
    this._send(MSG_TYPES.TYPING_START);
  }

  closeSession() {
    if (!this.isConnected) return;
    this._send(MSG_TYPES.SESSION_CLOSE);
  }

  // --- 内部方法 ---
  _send(type, payload = {}) {
    this.ws.send(JSON.stringify({ type, payload }));
  }

  _reconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      this.onStatusChange("reconnect-failed");
      return;
    }

    const delay = 1000 * Math.pow(2, this.reconnectAttempts);
    setTimeout(() => {
      this.reconnectAttempts++;
      this.connect();
    }, delay);
  }

  disconnect() {
    this.reconnectAttempts = this.maxReconnectAttempts;
    this.ws?.close();
  }
}
